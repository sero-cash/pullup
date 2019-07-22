package app

import (
	"crypto/ecdsa"
	"github.com/pborman/uuid"
	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/accounts"
	"github.com/sero-cash/go-sero/accounts/keystore"
	"github.com/sero-cash/go-sero/crypto"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/rlp"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
)

type Account struct {
	wallet   accounts.Wallet
	pk       *keys.Uint512
	tk       *keys.Uint512
	skr      keys.PKr
	mainPkr  keys.PKr
	balances map[string]*big.Int
	utxoNums map[string]uint64

	//use for map sort
	at uint64

	isChanged bool
	keyPath   string
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

func (account *Account) Create(passphrase string) error {

	var privateKey *ecdsa.PrivateKey
	// If not loaded, generate random.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}
	// Create the keyfile object with a random UUID.
	id := uuid.NewRandom()
	address := crypto.PrivkeyToAddress(privateKey)
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PrivkeyToAddress(privateKey),
		Tk:         crypto.PrivkeyToTk(privateKey),
		PrivateKey: privateKey,
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
				//self.accountMap.Delete(pk)
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
	if _, ok := self.accounts.Load(*w.Accounts()[0].Address.ToUint512()); !ok {
		account := Account{}
		account.wallet = w
		account.pk = w.Accounts()[0].Address.ToUint512()
		account.tk = w.Accounts()[0].Tk.ToUint512()
		account.at= w.Accounts()[0].At

		copy(account.skr[:], account.tk[:])
		account.mainPkr = self.createPkr(account.pk, 1)
		self.accounts.Store(*account.pk, &account)
		account.isChanged = true

		self.recoverPkrIndex(account, w.Accounts()[0].At)
	}
}

func (self *SEROLight) recoverPkrIndex(account Account, at uint64) {
	pk := *account.pk
	value, _ := self.db.Get(append(pkrIndexPrefix, pk[:]...))
	if value == nil {
		self.accountMap.Store(pk, outReq{Num: at, Pkr: account.mainPkr, PkrIndex: 1})
		self.createPkr(account.pk, uint64(1))
	} else {
		var otq outReq
		err :=rlp.DecodeBytes(value, &otq)
		if err!=nil{
			return
		}
		self.accountMap.Store(pk, otq)
		pkrNum := int(0)
		currentPkrIndex := otq.PkrIndex
		for i := int(currentPkrIndex); i > pkrNum; i-- {
			self.createPkr(account.pk, uint64(i))
		}
	}
}