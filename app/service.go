package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/accounts"
	"github.com/sero-cash/go-sero/accounts/keystore"
	"github.com/sero-cash/go-sero/common"
	"github.com/sero-cash/go-sero/common/address"
	"github.com/sero-cash/go-sero/common/hexutil"
	"github.com/sero-cash/go-sero/crypto"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/pullup/common/transport"
	"github.com/tyler-smith/go-bip39"
	"io/ioutil"
	"math/big"
	"net/http"
	"sort"
)

//keystore file upload
const maxUploadSize = 1 * 1024 * 2014 // 2 MB

type Service interface {
	NewAccountWithMnemonic(passphrase string) (map[string]string, error)
	UploadKeystoreHandler() http.HandlerFunc
	ImportAccountFromMnemonic(mnemonic, password string) (map[string]string, error)
	ImportAccountFromRawKey(privkey, password string) (map[string]string, error)
	ExportMnemonic(addressStr, password string) (string, error)
	AccountList() accountResps
	AccountDetail(pkStr string) accountResp
	AccountBalance(pkStr string) map[string]*big.Int
	TXNum(pkStr string) map[string]uint64
	TXList(pkStr string, request transport.PageRequest) (utxosResp, error)

	Transfer(from, to, currency, amount, gasPrice, pwd string) (hash keys.Uint256, err error)
	GetDecimal(currency string) uint64

	registerStakePool(from, vote, passwd string, feeRate uint32) (txHash string, err error)
	buyStake(from, vote, passwd, pool, amountStr, gaspriceStr string) (txHash string, err error)

	getSetNetwork(host string) string

	InitHost(rpcHostCustomer, webHostCustomer string)
}

func NewServiceAPI() Service {
	return &ServiceApi{
		SL: currentLight,
	}
}

type ServiceApi struct {
	SL *SEROLight
}

func (s *ServiceApi) ExportMnemonic(addressStr, password string) (string, error) {
	return fetchKeystore(s.SL.accountManager).ExportMnemonic(accounts.Account{Address: address.Base58ToAccount(addressStr)}, password)
}

// fetchKeystore retrives the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) *keystore.KeyStore {
	return am.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

func (s *ServiceApi) NewAccountWithMnemonic(passphrase string) (map[string]string, error) {
	blockNum := s.SL.getAccountBlock()

	mnemonic, acc, err := fetchKeystore(s.SL.accountManager).NewAccountWithMnemonic(passphrase, blockNum)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	result["mnemonic"] = mnemonic
	result["address"] = acc.Address.Base58()
	return result, nil
}

func (s *ServiceApi) ImportAccountFromMnemonic(mnemonic, password string) (map[string]string, error) {
	_, err := bip39.MnemonicToByteArray(mnemonic)
	if err != nil {
		return nil, err
	}
	seed, err := bip39.EntropyFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	if len(seed) != 32 {
		return nil, errors.New("EntropyFromMnemonic error seed not 256bits")
	}
	key, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return nil, err
	}
	acc, err := fetchKeystore(s.SL.accountManager).ImportECDSA(key, password)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	result["address"] = acc.Address.Base58()
	return result, nil
}

func (s *ServiceApi) ImportAccountFromRawKey(privkey, password string) (map[string]string, error) {
	key, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return nil, err
	}
	acc, err := fetchKeystore(s.SL.accountManager).ImportECDSA(key, password)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	result["address"] = acc.Address.Base58()
	return result, nil
}

type accountResp struct {
	PK         string
	MainPKr    string
	MainOldPKr string
	Balance    map[string]*big.Int
	UtxoNums   map[string]uint64
	PkrBase58  string
	at         uint64
	Name       string

	initTimestamp int64
}

type accountResps []accountResp

func (acrs accountResps) Len() int {
	return len(acrs)
}
func (acrs accountResps) Less(i, j int) bool {
	return acrs[i].initTimestamp < acrs[j].initTimestamp
}
func (acrs accountResps) Swap(i, j int) {
	acrs[i], acrs[j] = acrs[j], acrs[i]
}

func (s *ServiceApi) AccountList() (accountListResps accountResps) {
	s.SL.accounts.Range(func(key, value interface{}) bool {
		pk := key.(keys.Uint512)
		account := value.(*Account)
		latestPKr := keys.PKr{}
		if v, ok := s.SL.pkrIndexMap.Load(pk); ok {
			o := v.(outReq)
			latestPKr = o.Pkr
		}
		balance := s.SL.GetBalances(pk)
		accountListResp := accountResp{PK: base58.Encode(pk[:]), MainPKr: base58.Encode(account.mainPkr[:]), MainOldPKr: base58.Encode(account.mainOldPkr[:]), Balance: balance, UtxoNums: account.utxoNums, PkrBase58: base58.Encode(latestPKr[:]), at: account.at, initTimestamp: account.initTimestamp, Name: account.name}
		accountListResps = append(accountListResps, accountListResp)
		return true
	})

	sort.Sort(accountListResps)

	return accountListResps
}

func (s *ServiceApi) AccountDetail(pkStr string) (account accountResp) {
	pk := *address.Base58ToAccount(pkStr).ToUint512()
	if ac := s.SL.getAccountByPk(pk); ac != nil {
		latestPKr := keys.PKr{}
		if v, ok := s.SL.pkrIndexMap.Load(pk); ok {
			o := v.(outReq)
			latestPKr = o.Pkr
		}
		balance := s.SL.GetBalances(pk)
		account := accountResp{PK: base58.Encode(pk[:]), MainPKr: base58.Encode(ac.mainPkr[:]), MainOldPKr: base58.Encode(ac.mainOldPkr[:]), Balance: balance, UtxoNums: ac.utxoNums, PkrBase58: base58.Encode(latestPKr[:]), Name: ac.name}

		return account
	}
	return account
}

func (s *ServiceApi) AccountBalance(pkStr string) map[string]*big.Int {
	pk := address.Base58ToAccount(pkStr)
	return s.SL.GetBalances(*pk.ToUint512())
}

type utxoResp struct {
	Id        uint64
	Type      uint64
	To        string
	Hash      keys.Uint256
	Block     uint64
	Currency  string
	Amount    *big.Int
	Fee       *big.Int
	Receipt   TxReceipt
	Timestamp uint64
}

type assetResp struct {
	Tkn tknResp
	Tkt tktResp
}

type tknResp struct {
	Currency string
	Value    big.Int
}

type tktResp struct {
	Category string
	Value    string
}

type utxosResp []utxoResp

func (u utxosResp) Len() int {
	return len(u)
}
func (u utxosResp) Less(i, j int) bool {
	return u[i].Block > u[j].Block
}
func (u utxosResp) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (s *ServiceApi) TXList(pkStr string, request transport.PageRequest) (utxos utxosResp, err error) {
	pk := address.Base58ToAccount(pkStr)

	if txs, err := s.SL.findTx(*pk.ToUint512(), uint64(request.PageSize)); err == nil {
		for _, tx := range txs {
			if tx.Block == 0 {
				tx.Block = uint64(999999999999999)
			}
			utxo := utxoResp{
				Type:      tx.Type,
				To:        base58.Encode(tx.To[:]),
				Currency:  common.BytesToString(tx.Currency[:]),
				Amount:    tx.Amount,
				Fee:       tx.Fee,
				Hash:      tx.Hash,
				Block:     tx.Block,
				Receipt:   tx.Receipt,
				Timestamp: tx.Timestamp,
			}
			if big.NewInt(0).Add(tx.Amount, tx.Fee).Sign() == 0 {
			} else {
				utxos = append(utxos, utxo)
			}
		}
		sort.Sort(utxos)
	}

	return
}

func (s *ServiceApi) Transfer(from, to, currency, amountStr, gasPriceStr, password string) (hash keys.Uint256, err error) {

	amount, err := NewBigIntFromString(amountStr, 10)
	if err != nil {
		return hash, err
	} else {
		if amount.Sign() < 0 {
			return hash, fmt.Errorf("amount < 0")
		}
	}

	gasPrice, err := NewBigIntFromString(gasPriceStr, 10)
	if err != nil {
		return hash, err
	} else {
		if gasPrice.Sign() < 0 {
			return hash, fmt.Errorf("gasPrice < 0")
		}
	}
	if toBytes := base58.Decode(to); len(toBytes) != 96 {
		return hash, fmt.Errorf("Invalid colleaction address ")
	}
	h, err := s.SL.commitTx(from, to, currency, password, amount, gasPrice)
	if err != nil {
		return hash, err
	}
	return h, nil
}

func (s *ServiceApi) TXNum(pkStr string) map[string]uint64 {
	pk := address.Base58ToAccount(pkStr)
	return s.SL.GetUtxoNum(*pk.ToUint512())
}

func (s *ServiceApi) GetDecimal(currency string) uint64 {
	return s.SL.getDecimal(currency)
}

func renderError(w http.ResponseWriter, errcode string, code int) {
	//w.WriteHeader(code)
	w.Write([]byte(errcode))
}

func (s *ServiceApi) UploadKeystoreHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "FILE_TOO_BIG", http.StatusOK)
			return
		}
		file, _, err := r.FormFile("uploadFile")
		passphrase := r.FormValue("passphrase")
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusOK)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusOK)
			return
		}
		key, err := keystore.DecryptKey(fileBytes, passphrase)
		if err != nil {
			//renderError(w, "INVALID_FILE_TYPE", http.StatusOK)
			w.Write([]byte("INVALID_FILE_TYPE"))
			return
		}

		if err := ioutil.WriteFile(GetKeystorePath()+"/"+key.Address.String(), fileBytes, 0600); err != nil {
			renderError(w, "INVALID_FILE", http.StatusOK)
			return
		}

		logex.Infof("Import account successful. address=[%s]", key.Address)
		w.Write([]byte("SUCCESS"))
		return
	})
}

func (s *ServiceApi) registerStakePool(from, vote, passwd string, feeRate uint32) (txHash string, err error) {

	decimal := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil)
	amount := big.NewInt(0).Mul(big.NewInt(200000), decimal)
	if IsDev {
		amount = big.NewInt(0).Mul(big.NewInt(1), decimal)
	}
	gasprice := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(9), nil)
	hash, err := s.SL.registerStakePool(from, vote, passwd, feeRate, amount, gasprice)
	if err != nil {
		return txHash, err
	}
	return hexutil.Encode(hash[:]), nil
}

func (s *ServiceApi) buyStake(from, vote, passwd, pool, amountStr, gaspriceStr string) (txHash string, err error) {

	amount, err := NewBigIntFromString(amountStr, 10)
	if err != nil {
		return txHash, err
	} else {
		if amount.Sign() < 0 {
			return txHash, fmt.Errorf("amount < 0")
		}
	}

	gasPrice, err := NewBigIntFromString(gaspriceStr, 10)
	if err != nil {
		return txHash, err
	} else {
		if gasPrice.Sign() < 0 {
			return txHash, fmt.Errorf("gasPrice < 0")
		}
	}
	if len(vote) > 0 {
		if toBytes := base58.Decode(vote); len(toBytes) != 96 {
			return txHash, fmt.Errorf("Invalid vote address ")
		}
	}

	hash, err := s.SL.buyShare(from, vote, passwd, pool, amount, gasPrice)
	if err != nil {
		return txHash, err
	}
	return hexutil.Encode(hash[:]), nil
}

func (self *ServiceApi) getSetNetwork(hostReq string) string {
	if hostReq == "" {
		hostByte, err := self.SL.dbConfig.Get(hostKey)
		if err != nil {
			return GetRpcHost()
		}
		return string(hostByte[:])
	} else {
		self.SL.dbConfig.Put(hostKey, []byte(hostReq))
		setRpcHost(hostReq)
		return hostReq
	}
}

func (self *ServiceApi) InitHost(rpcHostCustomer, webHostCustomer string) {

	defaultRpcHost := "http://148.70.169.73:8545"
	defaultWebHost := "http://129.211.98.114:3006/web/v0_1_6/"

	//get remote rpc host
	resp, err := http.Get(remoteRpcHost)
	if err != nil {
		logex.Error("get remoteRpcHost Get err: ",err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logex.Error("get remoteRpcHost ReadAll err: ",err.Error())
		return
	}
	fmt.Println("get remote config success : ",string(body[:]))
	config := RpcConfig{}
	err =json.Unmarshal(body,&config)
	if err != nil {
		logex.Error("get remoteRpcHost Unmarshal err: ",err.Error())
		return
	}

	if config.Default.Rpc != "" {
		defaultRpcHost = config.Default.Rpc
	}
	if config.Default.Web != "" {
		defaultWebHost = config.Default.Web
	}
	fmt.Println("defaultRpcHoGst : ", defaultRpcHost)
	fmt.Println("defaultWebHost : ", defaultWebHost)
	if rpcHostCustomer != "" {
		setRpcHost(rpcHostCustomer)
		self.SL.dbConfig.Put(hostKey, []byte(rpcHostCustomer))
	} else {
		hostByte, err := self.SL.dbConfig.Get(hostKey)
		if err != nil {
			setRpcHost(defaultRpcHost)
		} else {
			setRpcHost(string(hostByte[:]))
		}
	}

	if webHostCustomer != "" {
		setWebHost(webHostCustomer)
	} else {
		setWebHost(defaultWebHost)
	}

}
