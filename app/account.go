package app

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/sero-cash/go-czero-import/superzk"

	"github.com/btcsuite/btcutil/base58"

	"github.com/sero-cash/go-sero/common/address"

	"github.com/sero-cash/go-czero-import/c_superzk"
	"github.com/sero-cash/go-czero-import/seroparam"

	"github.com/pborman/uuid"
	"github.com/sero-cash/go-czero-import/c_type"
	"github.com/sero-cash/go-sero/accounts"
	"github.com/sero-cash/go-sero/accounts/keystore"
	"github.com/sero-cash/go-sero/common"
	"github.com/sero-cash/go-sero/crypto"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/rlp"
)

type Account struct {
	wallet     accounts.Wallet
	pk         *c_type.Uint512
	tk         *c_type.Tk
	skr        c_type.PKr
	mainPkr    c_type.PKr
	mainOldPkr c_type.PKr
	balances   map[string]*big.Int
	utxoNums   map[string]uint64
	//use for map sort
	at            uint64
	isChanged     bool
	keyPath       string
	initTimestamp int64
	name          string
	version       int
}

func makeAccountManager() (*accounts.Manager, error) {
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	keydir := GetKeystorePath()
	var err error
	if err != nil {
		return nil, err
	}
	// Assemble the account manager and supported backends
	backends := []accounts.Backend{
		keystore.NewKeyStore(keydir, scryptN, scryptP),
	}
	return accounts.NewManager(backends...), nil
}

func (account *Account) Create(passphrase string, at uint64) error {

	var privateKey *ecdsa.PrivateKey
	// If not loaded, generate random.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}
	version := 1
	if at >= seroparam.SIP5() {
		version = 2
	}
	// Create the keyfile object with a random UUID.
	id := uuid.NewRandom()
	tk :=crypto.PrivkeyToTk(privateKey,version)
	address := tk.ToPk()
	key := &keystore.Key{
		Id:         id,
		Address:    address,
		Tk:         tk,
		PrivateKey: privateKey,
		Version:    version,
		At:         at,
	}

	// Encrypt key with passphrase.
	keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return err
	}
	// Store the file to disk.
	if err := os.MkdirAll(filepath.Dir(GetKeystorePath()+"/"+address.String()), 0700); err != nil {
		logex.Fatalf("Could not create directory %s", filepath.Dir(GetKeystorePath()))
		return err
	}
	if err := ioutil.WriteFile(GetKeystorePath()+"/"+address.String(), keyjson, 0600); err != nil {
		logex.Fatalf("Failed to write keyfile to %s: %v", GetKeystorePath(), err)
		return err
	}
	// Output some information.
	//logex.Infof("Create account successful. address =[%s]", key.Address)
	return nil

}

func (account *Account) Import(passphrase, keyPath string) error {
	keyJson, err := ioutil.ReadFile(keyPath)
	if err != nil {
		logex.Errorf("Failed to read the keyfile at '%s': %v", keyPath, err)
		return err
	}
	// Decrypt key with passphrase.
	key, err := keystore.DecryptKey(keyJson, passphrase)
	if err != nil {
		logex.Errorf("Error decrypting key: %v", err)
		return err
	}
	// Then write the new keyfile in place of the old one.
	if err := ioutil.WriteFile(GetKeystorePath()+"/"+key.Address.String(), keyJson, 0600); err != nil {
		logex.Errorf("Error writing new keyFile to disk: %v", err)
		return err
	}
	//logex.Infof("Import account successful. address=[%s]", key.Address)
	return nil
}

func (account *Account) UpdatePass(oldPas, newPass string) error {
	keyJson, err := ioutil.ReadFile(account.keyPath)
	if err != nil {
		//logex.Errorf("Failed to read the keyfile at '%s': %v", account.keyPath, err)
		return err
	}
	// Decrypt key with passphrase.
	key, err := keystore.DecryptKey(keyJson, oldPas)
	if err != nil {
		//logex.Errorf("Error decrypting key: %v", err)
		return err
	}

	// Encrypt the key with the new passphrase.
	newJson, err := keystore.EncryptKey(key, newPass, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		//logex.Errorf("Error encrypting with new passphrase: %v", err)
	}

	// Then write the new keyfile in place of the old one.
	if err := ioutil.WriteFile(GetKeystorePath()+"/"+key.Address.String(), newJson, 0600); err != nil {
		//logex.Errorf("Error writing new keyFile to disk: %v", err)
		return err
	}
	//logex.Infof("Update account pass successful. address=[%s]", key.Address)
	return nil
}

func (account *Account) Export() {

}
func (account *Account) PkString() string {
	var addr address.PKAddress
	copy(addr[:], account.pk[:])
	return addr.String()
}
func (a Account) PkrString(pkr c_type.PKr) string {
	return base58.Encode(pkr[:])
}

func (self *SEROLight) keystoreListener() {
	// Close all subscriptions when the manager terminates
	defer func() {
		self.lock.Lock()
		self.updater.Unsubscribe()
		self.updater = nil
		self.lock.Unlock()
	}()
	// Loop until termination
	for {
		select {
		case event := <-self.update:
			// Wallet event arrived, update local cache
			self.lock.Lock()
			switch event.Kind {
			case accounts.WalletArrived:
				//wallet := event.Wallet
				self.initWallet(event.Wallet)
			case accounts.WalletDropped:
				//pk := *event.Wallet.Accounts()[0].Address.ToUint512()
				//self.pkrIndexMap.Delete(pk)
			}
			self.lock.Unlock()

		case errc := <-self.quit:
			// Manager terminating, return
			errc <- nil
			return
		}
	}
}

func (self *SEROLight) initWallet(w accounts.Wallet) {
	if _, ok := self.accounts.Load(w.Accounts()[0].GetPk()); !ok {
		account := Account{}
		account.wallet = w
		account.pk = w.Accounts()[0].GetPk().NewRef()
		tk := w.Accounts()[0].Tk.ToTk()
		account.tk = &tk
		account.at = w.Accounts()[0].At
		account.version = w.Accounts()[0].Version
		copy(account.skr[:], account.tk[:])
		mainPkr, err := self.createPkrHash(account.tk, 1)
		if err != nil {
			panic("init account failed pk = " + account.PkString())
		}
		account.mainPkr = *mainPkr
		if account.version == 1 {
			oldPkr, err := self.createPkr(account.tk, 1)
			if err != nil {
				panic("init account failed pk = " + account.PkString())
			}
			account.mainOldPkr = *oldPkr
		}

		self.accounts.Store(*account.pk, &account)
		account.isChanged = true
		account.initTimestamp = time.Now().UnixNano()
		self.recoverPkrIndex(account, w.Accounts()[0].At)

		var keystoreName = w.URL().Path[len(GetKeystorePath()):]
		var split = "ac_"
		if keystoreName[:len(split)] == split {
			fmt.Println("customer account name : ", keystoreName[len(split):])
			account.name = keystoreName[len(split):]
		}

		fmt.Println("init wallet :", account.PkString())
	}
}

func (self *SEROLight) recoverPkrIndex(account Account, at uint64) {
	pk := *account.pk
	value, _ := self.db.Get(append(pkrIndexPrefix, pk[:]...))
	if value == nil {
		self.pkrIndexMap.Store(pk, outReq{Num: at, Pkr: account.mainPkr, PkrIndex: 1})
	} else {
		var otq outReq
		err := rlp.DecodeBytes(value, &otq)
		if err != nil {
			return
		}
		self.pkrIndexMap.Store(pk, otq)
	}

	if data, err := self.db.Get(append(onlyUseHashPkrKey, pk[:]...)); err == nil {
		value := decodeNumber(data)
		if value == 1 {
			self.useHashPkr.Store(account.pk, 1)
		}
	}
}

func (self *SEROLight) createPkr(tk *c_type.Tk, index uint64) (*c_type.PKr, error) {
	r := c_type.Uint256{}
	copy(r[:], common.LeftPadBytes(encodeNumber(index), 32))
	pk, err := superzk.Tk2Pk(tk)
	if err != nil {
		return nil, err
	}
	pkr, err := c_superzk.Czero_PK2PKr(&pk, &r)
	if err != nil {
		return nil, err
	}
	return &pkr, nil
}

func (self *SEROLight) createPkrHash(tk *c_type.Tk, index uint64) (*c_type.PKr, error) {
	random := append(tk[:], encodeNumber(index)[:]...)
	r := crypto.Keccak256Hash(random).HashToUint256()
	pk, err := superzk.Tk2Pk(tk)
	if err != nil {
		return nil, err
	}
	pkr := superzk.Pk2PKr(&pk, r)
	if err != nil {
		return nil, err
	}
	//fmt.Println("hashPkr: ", base58.Encode(pkr[:]))

	return &pkr, nil
}
