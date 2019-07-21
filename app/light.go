package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/accounts"
	"github.com/sero-cash/go-sero/common"
	"github.com/sero-cash/go-sero/common/address"
	"github.com/sero-cash/go-sero/common/hexutil"
	"github.com/sero-cash/go-sero/crypto"
	"github.com/sero-cash/go-sero/event"
	"github.com/sero-cash/go-sero/light-wallet/common/logex"
	"github.com/sero-cash/go-sero/rlp"
	"github.com/sero-cash/go-sero/serodb"
	"github.com/sero-cash/go-sero/zero/txs/assets"
	"github.com/sero-cash/go-sero/zero/txs/stx"
	"github.com/sero-cash/go-sero/zero/txtool"
	"github.com/sero-cash/go-sero/zero/txtool/flight"
	"github.com/sero-cash/go-sero/zero/txtool/prepare"
	"github.com/sero-cash/go-sero/zero/utils"
	"math/big"
	"strconv"
	"sync"
)

var (
	maxUint64  = ^uint64(0)
	fetchCount = uint64(5000)
)

type SEROLight struct {
	db *serodb.LDBDatabase

	accountManager *accounts.Manager
	accounts       sync.Map
	usedFlag       sync.Map
	accountMap     sync.Map
	pkrIndexMap    sync.Map
	sli            flight.SLI

	// SERO wallet subscription
	feed    event.Feed
	updater event.Subscription        // Wallet update subscriptions for all backends
	update  chan accounts.WalletEvent // Subscription sink for backend wallet changes
	quit    chan chan error
	lock    sync.RWMutex

}

var currentLight *SEROLight

func NewSeroLight() {

	accountManager, err := makeAccountManager()
	if err != nil {
		logex.Fatalf("makeAccountManager, err=[%v]", err)
	}
	db, err := serodb.NewLDBDatabase(GetDataPath(), 1024, 1024)
	if err != nil {
		logex.Fatalf("NewLDBDatabase, err=[%v]", err)
	}
	update := make(chan accounts.WalletEvent, 1)
	updater := accountManager.Subscribe(update)

	light := &SEROLight{}
	light.accountManager = accountManager
	light.update = update
	light.updater = updater
	light.db = db
	light.accountMap = sync.Map{}
	light.accounts = sync.Map{}
	light.usedFlag = sync.Map{}
	light.pkrIndexMap = sync.Map{}

	currentLight = light

	for _, w := range accountManager.Wallets() {
		light.initWallet(w)
	}

	AddJob("0/10 * * * * ?", light.SyncOut)
	go light.keystoreListener()
}

// sync out request base params
type outReq struct {
	PkrIndex uint64
	Pkr      keys.PKr
	Num      uint64
}

type fetchReturn struct {
	utxoMap   map[PkKey][]Utxo
	again     bool
	remoteNum uint64
	nextNum   uint64
}

func (self *SEROLight) SyncOut() {
	if host == ""{
		return
	}
	self.accountMap.Range(func(key, value interface{}) bool {
		pk := key.(keys.Uint512)
		otreq := value.(outReq)
		for {
			fmt.Printf("otreq.Pkr:%s,[%d],Num:[%d] \n ", base58.Encode(pk[:]), otreq.PkrIndex, otreq.Num)
			pkrs := self.getBeforePKrs(pk, otreq.PkrIndex)
			if len(pkrs) == 0 {
				return false
			}
			var start, end = otreq.Num, otreq.Num+fetchCount
			account := self.getAccountByPk(pk)
			rtn, err := self.fetchAndDecOuts(account, pkrs, otreq.Pkr, start, end)
			if err != nil {
				logex.Errorf(err.Error())
				return false
			}
			otreq.Num = rtn.nextNum
			if len(rtn.utxoMap) > 0 {
				batch := self.db.NewBatch()
				err = self.indexOuts(rtn.utxoMap, batch)
				if err != nil {
					logex.Errorf(err.Error())
					return false
				}
				err = batch.Write()
				if err != nil {
					return false
				}
				account.isChanged = true
			}

			if rtn.again {
				otreq.PkrIndex = otreq.PkrIndex + 1
				otreq.Pkr = self.createPkr(&pk, otreq.PkrIndex)
				data, _ := rlp.EncodeToBytes(otreq)
				var test outReq
				rlp.DecodeBytes(data, &test)
				self.accountMap.Store(pk, otreq)
				self.db.Put(append(pkrIndexPrefix, pk[:]...), data)
				continue
			}
			data, _ := rlp.EncodeToBytes(otreq)
			self.accountMap.Store(pk, otreq)
			self.db.Put(append(pkrIndexPrefix, pk[:]...), data)
			if end >= rtn.remoteNum {
				break
			}
		}
		return true
	})
	self.CheckNil()
}

func (self *SEROLight) fetchAndDecOuts(account *Account, pkrs []string, currentPkr keys.PKr, start, end uint64) (rtn fetchReturn, err error) {

	sync := Sync{RpcHost: host, Method: "light_getOutsByPKr", Params: []interface{}{pkrs, start, end}}
	jsonResp, err := sync.Do()
	if err != nil {
		logex.Errorf("jsonRep err=[%s]", err.Error())
		return
	}
	bor := BlockOutResp{}
	if err = json.Unmarshal(*jsonResp.Result, &bor); err != nil {
		logex.Errorf("json.Unmarshal err=[%s]", err.Error())
		return
	}
	utxosMap := map[PkKey][]Utxo{}
	// if not find outs , the end block return query current block
	blockOuts := bor.BlockOuts
	rtn.remoteNum = bor.CurrentNum
	if rtn.remoteNum > end {
		rtn.nextNum = end + 1
	} else {
		rtn.nextNum = bor.CurrentNum + 1
	}

	for _, blockOut := range blockOuts {
		outs := blockOut.Outs
		for _, out := range outs {
			var pkr keys.PKr
			if out.State.OS.Out_Z != nil {
				pkr = out.State.OS.Out_Z.PKr
			}
			if out.State.OS.Out_O != nil {
				pkr = out.State.OS.Out_O.Addr
			}
			if pkr == currentPkr {
				rtn.again = true
				//gen min block Num
				if rtn.nextNum > blockOut.Num {
					rtn.nextNum = blockOut.Num
				}
			}
			dout := DecOuts([]txtool.Out{out}, &account.skr)[0]
			key := PkKey{PK: *account.pk, Num: blockOut.Num}
			utxo := Utxo{Pkr: pkr, Root: out.Root, Nil: dout.Nil, TxHash: out.State.TxHash, Num: out.State.Num, Asset: dout.Asset, IsZ: out.State.OS.Out_Z != nil, Out: out}
			//log.Info("DecOuts", "PK", base58.Encode(account.pk[:]), "root", common.Bytes2Hex(out.Root[:]), "currency", common.BytesToString(utxo.Asset.Tkn.Currency[:]), "value", utxo.Asset.Tkn.Value)
			if list, ok := utxosMap[key]; ok {
				utxosMap[key] = append(list, utxo)
			} else {
				utxosMap[key] = []Utxo{utxo}
			}
		}
	}
	rtn.utxoMap = utxosMap
	return rtn, nil
}

//if the currentpkr in the outs, again = true, then loop continue next Pkr
func (self *SEROLight) indexOuts(utxosMap map[PkKey][]Utxo, batch serodb.Batch) (err error) {
	if len(utxosMap) > 0 {
		ops, err := self.indexUtxo(utxosMap, batch)
		if err != nil {
			return err
		}
		for key, value := range ops {
			batch.Put(common.Hex2Bytes(key), common.Hex2Bytes(value))
		}
	}
	return err
}

func (self *SEROLight) getBeforePKrs(pk keys.Uint512, currentPkrIndex uint64) (pkrs []string) {
	pkrNum := int(0)
	if currentPkrIndex > 5 {
		pkrNum = int(currentPkrIndex) - 5
	}
	for i := int(currentPkrIndex); i > pkrNum; i-- {
		pkr, err := self.getPKrIndex(pk, uint64(i))
		if err != nil {
			pkr = self.createPkr(&pk, uint64(i))
			pkrs = append(pkrs, base58.Encode(pkr[:]))
		} else {
			pkrs = append(pkrs, base58.Encode(pkr[:]))
		}
	}
	return pkrs
}

type pkrIndexKey struct {
	pk    keys.Uint512
	index uint64
}

func (self *SEROLight) getPKrIndex(pk keys.Uint512, index uint64) (pkr keys.PKr, err error) {
	if value, ok := self.pkrIndexMap.Load(pkrIndexKey{pk, index}); ok {
		return value.(keys.PKr), nil
	} else {
		return pkr, fmt.Errorf("not fund")
	}
}

func (self *SEROLight) setPKrIndex(pk keys.Uint512, index uint64, pkr keys.PKr) {
	self.pkrIndexMap.Store(pkrIndexKey{pk, index}, pkr)
}

func (self *SEROLight) indexUtxo(utxosMap map[PkKey][]Utxo, batch serodb.Batch) (opsReturn map[string]string, err error) {
	ops := map[string]string{}
	for key, list := range utxosMap {
		roots := []keys.Uint256{}
		for _, utxo := range list {
			data, err := rlp.EncodeToBytes(utxo)
			if err != nil {
				return nil, err
			}
			// "ROOT" + root
			batch.Put(rootKey(utxo.Root), data)

			//"TXHASH" + PK + hash + root + outType
			batch.Put(indexTxKey(key.PK, utxo.TxHash, utxo.Root, uint64(1)), data)

			//nil => root
			batch.Put(nilToRootKey(utxo.Nil), utxo.Root[:])

			var pkKey []byte
			if utxo.Asset.Tkn != nil {
				// "PK" + PK + currency + root
				pkKey = utxoPkKey(key.PK, utxo.Asset.Tkn.Currency[:], &utxo.Root)

			} else if utxo.Asset.Tkt != nil {
				// "PK" + PK + tkt + root
				pkKey = utxoPkKey(key.PK, utxo.Asset.Tkt.Value[:], &utxo.Root)
			}
			// "PK" + PK + currency + root => 0
			ops[common.Bytes2Hex(pkKey)] = common.Bytes2Hex([]byte{0})

			// "NIL" + PK + tkt + root => "PK" + PK + currency + root
			nilkey := nilKey(utxo.Nil)
			rootkey := nilKey(utxo.Root)
			//nilIdkey := nilIdKey(utxo.Nil)

			// "NIL" +nil/root => pkKey
			ops[common.Bytes2Hex(nilkey)] = common.Bytes2Hex(pkKey)
			ops[common.Bytes2Hex(rootkey)] = common.Bytes2Hex(pkKey)
			//ops[common.Bytes2Hex(nilIdkey)] = common.Bytes2Hex(encodeNumber(key.Num))
			roots = append(roots, utxo.Root)
			//log.Info("Index add", "PK", base58.Encode(key.PK[:]), "Nil", common.Bytes2Hex(utxo.Nil[:]), "root", common.Bytes2Hex(utxo.Root[:]), "Value", utxo.Asset.Tkn.Value)
		}
		data, err := rlp.EncodeToBytes(roots)
		if err != nil {
			return nil, err
		}

		// utxo PK + at  => [roots]
		batch.Put(utxoKey(key.Num, key.PK), data)
	}
	return ops, nil
}

type NilValue struct {
	Nil    keys.Uint256
	Num    uint64
	TxHash keys.Uint256
	TxFee  big.Int
}

func (self *SEROLight) CheckNil() {

	iterator := self.db.NewIteratorWithPrefix(nilPrefix)
	Nils := []keys.Uint256{}

	for iterator.Next() {
		key := iterator.Key()
		var Nil keys.Uint256
		copy(Nil[:], key[3:])
		nilkey := nilKey(Nil)
		value, _ := self.db.Get(nilkey)
		if value != nil {
			Nils = append(Nils, Nil)
		}

	}

	sync := Sync{RpcHost: host, Method: "light_checkNil", Params: []interface{}{Nils}}
	jsonResp, err := sync.Do()
	if err != nil {
		logex.Errorf("jsonRep err=[%s]", err.Error())
		return
	}
	if jsonResp.Result != nil {
		nilvs := []NilValue{}
		if err = json.Unmarshal(*jsonResp.Result, &nilvs); err != nil {
			logex.Errorf("json.Unmarshal err=[%s]", err.Error())
			return
		}
		if len(nilvs) > 0 {
			batch := self.db.NewBatch()
			for _, nilv := range nilvs {
				var pk keys.Uint512
				Nil := nilv.Nil
				value, _ := self.db.Get(nilKey(Nil))
				if value != nil {
					copy(pk[:], value[2:66])
					var root keys.Uint256
					copy(root[:], value[98:130])
					utxo, err := self.getUtxo(root)
					if err == nil {
						batch.Delete(penddingTxKey(pk, utxo.TxHash))
					}
					if len(value) == 130 {
						batch.Delete(value)
					} else {
						batch.Delete(value[0:130])
						batch.Delete(value[130:260])
					}
					batch.Delete(nilKey(Nil))
					batch.Delete(nilKey(root))

					//remove pending tx
					batch.Delete(indexTxKey(pk, nilv.TxHash, nilv.TxHash, uint64(2)))
					utxoI := Utxo{Root: root, TxHash: nilv.TxHash, Fee: nilv.TxFee, Num: nilv.Num, Nil: nilv.Nil, Asset: utxo.Asset, Pkr: utxo.Pkr}
					data, _ := rlp.EncodeToBytes(utxoI)
					batch.Put(indexTxKey(pk, nilv.TxHash, root, uint64(2)), data)

					self.usedFlag.Delete(root)
					logex.Info("delete :", hexutil.Encode(root[:]))
				}
				if account := self.getAccountByPk(pk); account != nil {
					account.isChanged = true
				}
			}
			batch.Write()
		}
	}

}

func (self *SEROLight) getAccountByPk(pk keys.Uint512) *Account {
	if value, ok := self.accounts.Load(pk); ok {
		return value.(*Account)
	}
	return nil
}

// "UTXO" + pk + number
func utxoKey(number uint64, pk keys.Uint512) []byte {
	return append(utxoPrefix, append(pk[:], encodeNumber(number)...)...)
}

// utxoKey = PK + currency +root
func utxoPkKey(pk keys.Uint512, currency []byte, root *keys.Uint256) []byte {
	key := append(pkPrefix, pk[:]...)
	if len(currency) > 0 {
		key = append(key, currency...)
	}
	if root != nil {
		key = append(key, root[:]...)
	}
	return key
}

func (self *SEROLight) GetUtxoNum(pk keys.Uint512) map[string]uint64 {
	if account := self.getAccountByPk(pk); account != nil {
		return account.utxoNums
	}
	return map[string]uint64{}
}

func (self *SEROLight) GetBalances(pk keys.Uint512) (balances map[string]*big.Int) {
	if value, ok := self.accounts.Load(pk); ok {
		account := value.(*Account)
		if account.isChanged {
			prefix := append(pkPrefix, pk[:]...)
			iterator := self.db.NewIteratorWithPrefix(prefix)
			balances = map[string]*big.Int{}
			utxoNums := map[string]uint64{}
			for iterator.Next() {
				key := iterator.Key()
				var root keys.Uint256
				copy(root[:], key[98:130])
				if utxo, err := self.getUtxo(root); err == nil {
					if utxo.Asset.Tkn != nil {
						currency := common.BytesToString(utxo.Asset.Tkn.Currency[:])
						if amount, ok := balances[currency]; ok {
							amount.Add(amount, utxo.Asset.Tkn.Value.ToIntRef())
							utxoNums[currency] += 1
						} else {
							balances[currency] = new(big.Int).Set(utxo.Asset.Tkn.Value.ToIntRef())
							utxoNums[currency] = 1
						}
					}
				}
			}
			account.balances = balances
			account.utxoNums = utxoNums
			account.isChanged = false
		} else {
			return account.balances
		}
	}
	return
}

func (self *SEROLight) getUtxo(root keys.Uint256) (utxo Utxo, e error) {
	data, err := self.db.Get(rootKey(root))
	if err != nil {
		return
	}
	if err := rlp.Decode(bytes.NewReader(data), &utxo); err != nil {
		logex.Error("Light Invalid utxo RLP", "root", common.Bytes2Hex(root[:]), "err", err)
		e = err
		return
	}
	if value, ok := self.usedFlag.Load(utxo.Root); ok {
		utxo.flag = value.(int)
	}
	return
}

func (self *SEROLight) commitTx(from, to, currency, passwd string, amount, gasprice *big.Int) (hash keys.Uint256, err error) {

	fee := new(big.Int).Mul(big.NewInt(25000), gasprice)
	fromPk := address.Base58ToAccount(from).ToUint512()

	var RefundTo *keys.PKr
	ac := self.getAccountByPk(*fromPk)
	if ac != nil {
		RefundTo = &ac.mainPkr
	}

	account := accounts.Account{Address: ac.wallet.Accounts()[0].Address}
	wallet, err := self.accountManager.Find(account)
	if err != nil {
		return hash, err
	}
	seed, err := wallet.GetSeedWithPassphrase(passwd)
	if err != nil {
		return hash, err
	}
	var toPkr keys.PKr
	copy(toPkr[:], base58.Decode(to)[:])
	reception := self.genReceiption(currency, amount, toPkr)

	preTxParam := prepare.PreTxParam{}
	preTxParam.From = *fromPk
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = gasprice
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}
	preTxParam.Receptions = []prepare.Reception{reception}

	param, err := self.GenTx(preTxParam)
	if err != nil {
		return hash, err
	}
	sk := keys.Seed2Sk(seed.SeedToUint256())
	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return hash, err
	}
	hash = gtx.Hash
	sync := Sync{RpcHost: host, Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return hash, err
	}

	utxoIn := Utxo{Pkr: toPkr, Root: hash, TxHash: hash, Fee: *fee}
	self.storePeddingUtxo(param, currency, amount, utxoIn, fromPk)

	return hash, nil
}

func (self *SEROLight) storePeddingUtxo(param *txtool.GTxParam, currency string, amount *big.Int, utxoIn Utxo, fromPk *keys.Uint512) {
	roots := []keys.Uint256{}
	for _, in := range param.Ins {
		roots = append(roots, in.Out.Root)
		self.usedFlag.Store(in.Out.Root, 1)
	}
	tknc := assets.Token{Currency: utils.CurrencyToUint256(currency), Value: utils.U256(*amount)}
	assetc := assets.Asset{}
	assetc.Tkn = &tknc
	utxoIn.Asset = assetc
	dataIn, _ := rlp.EncodeToBytes(utxoIn)
	self.db.Put(indexTxKey(*fromPk, utxoIn.TxHash, utxoIn.TxHash, uint64(2)), dataIn)
}

func (self *SEROLight) genReceiption(currency string, amount *big.Int, toPkr keys.PKr) prepare.Reception {
	tkn := assets.Token{Currency: utils.CurrencyToUint256(currency), Value: utils.U256(*amount)}
	asset := assets.NewAsset(&tkn, nil)
	reception := prepare.Reception{
		Addr:  toPkr,
		Asset: asset,
	}
	return reception
}

func (self *SEROLight) registerStakePool(from, vote, passwd string,feeRate uint32, amount, gasprice *big.Int) (hash keys.Uint256, err error) {

	fee := new(big.Int).Mul(big.NewInt(25000), gasprice)
	fromPk := address.Base58ToAccount(from).ToUint512()

	var RefundTo *keys.PKr
	ac := self.getAccountByPk(*fromPk)
	if ac != nil {
		RefundTo = &ac.mainPkr
	}
	//check pk register pool
	poolId :=crypto.Keccak256(ac.mainPkr[:])
	sync := Sync{RpcHost: host, Method: "stake_poolState", Params: []interface{}{hexutil.Encode(poolId)}}
	_, err = sync.Do()
	if err != nil {
		if err.Error() != "stake pool not exists"{
			logex.Errorf("jsonRep err=[%s]", err.Error())
			return
		}
	}else{
		err = fmt.Errorf("stake pool exists")
		logex.Errorf("jsonRep err=[%s]", err.Error())
		return
	}
	account := accounts.Account{Address: ac.wallet.Accounts()[0].Address}
	wallet, err := self.accountManager.Find(account)
	if err != nil {
		return hash, err
	}
	seed, err := wallet.GetSeedWithPassphrase(passwd)
	if err != nil {
		return hash, err
	}
	var votePkr keys.PKr
	if vote == "" {
		votePkr = ac.mainPkr
	}else{
		copy(votePkr[:], base58.Decode(vote)[:])
	}
	registerPool := stx.RegistPoolCmd{Value: utils.U256(*amount), Vote: votePkr, FeeRate:feeRate}
	preTxParam := prepare.PreTxParam{}
	preTxParam.From = *fromPk
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = gasprice
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}

	preTxParam.Cmds = prepare.Cmds{RegistPool:&registerPool}

	param, err := self.GenTx(preTxParam)

	if err != nil {
		return hash, err
	}

	sk := keys.Seed2Sk(seed.SeedToUint256())

	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return hash, err
	}

	hash = gtx.Hash
	logex.Info("commit txhash: ", hash)
	sync = Sync{RpcHost: host, Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return hash, err
	}

	utxoIn := Utxo{Pkr: votePkr, Root: hash, TxHash: hash, Fee: *fee}
	self.storePeddingUtxo(param, "SERO", amount, utxoIn, fromPk)

	return hash, nil
}

func (self *SEROLight) buyShare(from, vote, passwd, pool string, amount, gasprice *big.Int) (hash keys.Uint256, err error) {

	fee := new(big.Int).Mul(big.NewInt(25000), gasprice)
	fromPk := address.Base58ToAccount(from).ToUint512()

	var RefundTo *keys.PKr
	ac := self.getAccountByPk(*fromPk)
	if ac != nil {
		RefundTo = &ac.mainPkr
	}

	account := accounts.Account{Address: ac.wallet.Accounts()[0].Address}
	wallet, err := self.accountManager.Find(account)
	if err != nil {
		return hash, err
	}
	seed, err := wallet.GetSeedWithPassphrase(passwd)
	if err != nil {
		return hash, err
	}
	var votePkr keys.PKr
	if len(vote) == 0 {
		votePkr = ac.mainPkr
	}else{
		copy(votePkr[:], base58.Decode(vote)[:])
	}
	poolId := common.HexToHash(pool)
	buyShareCmd := stx.BuyShareCmd{Value: utils.U256(*amount), Vote: votePkr, Pool: poolId.HashToUint256()}
	preTxParam := prepare.PreTxParam{}
	preTxParam.From = *fromPk
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = gasprice
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}
	preTxParam.Cmds = prepare.Cmds{BuyShare:&buyShareCmd}
	param, err := self.GenTx(preTxParam)

	if err != nil {
		return hash, err
	}
	sk := keys.Seed2Sk(seed.SeedToUint256())
	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return hash, err
	}
	hash = gtx.Hash
	logex.Info("commit txhash: ", hash)
	sync := Sync{RpcHost: host, Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return hash, err
	}

	utxoIn := Utxo{Pkr: votePkr, Root: hash, TxHash: hash, Fee: *fee}
	self.storePeddingUtxo(param, "SERO", amount, utxoIn, fromPk)


	return hash, nil
}

func (self *SEROLight) getDecimal(currency string) uint64 {
	if decimalData, err := self.db.Get(append(decimalPrefix, []byte(currency)[:]...)); err != nil {
		if decimalData == nil {
			sync := Sync{RpcHost: host, Method: "sero_getDecimal", Params: []interface{}{currency}}
			if jsonResp, err := sync.Do(); err != nil {
				return 0
			} else {
				var decimalStr string
				if err = json.Unmarshal(*jsonResp.Result, &decimalStr); err != nil {
					logex.Error("json.Unmarshal err=[%s]", err.Error())
					return 0
				}
				decimalStr=decimalStr[2:]
				decimal, _ := strconv.ParseUint(decimalStr, 16, 64)
				self.db.Put(append(decimalPrefix, []byte(currency)[:]...), encodeNumber(decimal))
				return decimal
			}
		} else {
			return 0
		}
	} else {
		return decodeNumber(decimalData)
	}
}

func (self *SEROLight) getAccountBlock() uint64 {
	number := uint64(0)
	self.accountMap.Range(func(key, value interface{}) bool {
		data := value.(outReq)
		if number < data.Num {
			number = data.Num
		}
		return true
	})
	return number
}

func (self *SEROLight) getLatestPKrs(pk keys.Uint512) (pais []pkrAndIndex) {
	prefix := append(pkrPrefix, pk[:]...)
	iterator := self.db.NewIteratorWithPrefix(prefix)
	count := 0
	for iterator.Next() {
		pai := pkrAndIndex{}
		key := iterator.Key()
		keyLen := len(key)
		pai.index = decodeNumber(key[keyLen-8:])
		// remove at=0 , save latest five pkrs
		if count > 5 {
			pais = append(pais[:1], pais[2:]...)
		}
		value := iterator.Value()
		var pkr keys.PKr
		copy(pkr[:], value[:])
		pai.pkr = pkr
		pais = append(pais, pai)
		count++
	}
	return pais
}

func (self *SEROLight) GetCurrencyNumber(pk keys.Uint512) uint64 {
	value, ok := self.accountMap.Load(pk)
	if !ok {
		return 0
	}
	return value.(outReq).Num
}

func (self *SEROLight) createPkr(pk *keys.Uint512, index uint64) keys.PKr {
	r := keys.Uint256{}
	copy(r[:], common.LeftPadBytes(encodeNumber(index), 32))
	pkr := keys.Addr2PKr(pk, &r)
	self.setPKrIndex(*pk, index, pkr)
	logex.Infof("createPkr,at=[%d],Pkr=[%s]", index, base58.Encode(pkr[:]))
	return pkr
}