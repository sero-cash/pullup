package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/sero-cash/go-sero/zero/account"

	"github.com/pkg/errors"

	"github.com/sero-cash/go-czero-import/c_superzk"

	"github.com/sero-cash/go-czero-import/c_type"

	"github.com/btcsuite/btcutil/base58"
	"github.com/sero-cash/go-sero/accounts"
	"github.com/sero-cash/go-sero/common"
	"github.com/sero-cash/go-sero/common/address"
	"github.com/sero-cash/go-sero/common/hexutil"
	"github.com/sero-cash/go-sero/core/types"
	"github.com/sero-cash/go-sero/crypto"
	"github.com/sero-cash/go-sero/event"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/rlp"
	"github.com/sero-cash/go-sero/serodb"
	"github.com/sero-cash/go-sero/zero/txs/assets"
	"github.com/sero-cash/go-sero/zero/txs/stx"
	"github.com/sero-cash/go-sero/zero/txtool"
	"github.com/sero-cash/go-sero/zero/txtool/flight"
	"github.com/sero-cash/go-sero/zero/txtool/prepare"
	"github.com/sero-cash/go-sero/zero/utils"
)

type SEROLight struct {
	db             *serodb.LDBDatabase
	dbConfig       *serodb.LDBDatabase
	accountManager *accounts.Manager
	accounts       sync.Map
	usedFlag       sync.Map
	pkrIndexMap    sync.Map
	sli            flight.SRI
	// SERO wallet subscription
	feed       event.Feed
	updater    event.Subscription        // Wallet update subscriptions for all backends
	update     chan accounts.WalletEvent // Subscription sink for backend wallet changes
	quit       chan chan error
	lock       sync.RWMutex
	useHashPkr sync.Map
}

var currentLight *SEROLight

func NewSeroLight() {

	logex.Info("App start ,version: ", version)
	// new AccountManage
	accountManager, err := makeAccountManager()
	if err != nil {
		logex.Fatalf("makeAccountManager, err=[%v]", err)
	}

	// new config db
	configdb, err := serodb.NewLDBDatabase(GetConfigPath(), 1024, 1024)
	if err != nil {
		logex.Fatalf("NewLDBDatabase, err=[%v]", err)
	}

	// check this version clean data on start
	if cleanData {
		versionByte, err := configdb.Get(VersonKey[:])
		if err != nil {
			configdb.Put(VersonKey[:], []byte(GetVersion()))
			// clean data
			CleanData()
		} else {
			if string(versionByte[:]) == GetVersion() {
				logex.Info("latest version:", string(versionByte[:]))
			} else {
				configdb.Put(VersonKey[:], []byte(GetVersion()))
				// clean data
				CleanData()
			}
		}
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
	light.dbConfig = configdb
	light.pkrIndexMap = sync.Map{}
	light.accounts = sync.Map{}
	light.usedFlag = sync.Map{}
	light.useHashPkr = sync.Map{}

	currentLight = light

	for _, w := range accountManager.Wallets() {
		light.initWallet(w)
	}

	AddJob("0/20 * * * * ?", light.SyncOut)
	go light.keystoreListener()
}

// sync out request base params
type outReq struct {
	PkrIndex uint64
	Pkr      c_type.PKr
	Num      uint64
}

type fetchReturn struct {
	utxoMap         map[PkKey][]Utxo
	again           bool
	remoteNum       uint64
	nextNum         uint64
	useHashPkr      bool
	lastBlockNumber uint64
}

func (self *SEROLight) SyncOut() {
	if rpcHost == "" {
		return
	}
	self.pkrIndexMap.Range(func(key, value interface{}) bool {
		pk := key.(c_type.Uint512)
		otreq := value.(outReq)

		start := otreq.Num
		end := uint64(0)
		for {
			// var start, end = otreq.Num, otreq.Num + fetchCount
			account := self.getAccountByKey(pk)
			rtn, err := self.fetchAndDecOuts(account, otreq.PkrIndex, start, end)
			if err != nil {
				logex.Errorf("fetchAndDecOuts,err=[%s]", err.Error())
				return false
			}
			if len(rtn.utxoMap) > 0 {
				account.isChanged = true
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
			}

			if rtn.useHashPkr {
				self.useHashPkr.Store(account.pk, 1)
				self.db.Put(append(onlyUseHashPkrKey, pk[:]...), encodeNumber(1))
			}

			if rtn.remoteNum > 0 {
				self.db.Put(remoteNumKey, encodeNumber(rtn.remoteNum+12))
			}

			if rtn.again {
				// otreq.Num = rtn.nextNum
				otreq.PkrIndex = otreq.PkrIndex + 1
				nextPkr, _ := self.createPkrHash(account.tk, otreq.PkrIndex)
				otreq.Pkr = *nextPkr
				data, _ := rlp.EncodeToBytes(otreq)
				self.pkrIndexMap.Store(pk, otreq)
				self.db.Put(append(pkrIndexPrefix, pk[:]...), data)
				if (end == 0) {
					end = rtn.lastBlockNumber
				}
				continue
			} else {
				otreq.Num = rtn.lastBlockNumber + 1
				// otreq.Num = rtn.nextNum
				data, _ := rlp.EncodeToBytes(otreq)
				self.pkrIndexMap.Store(pk, otreq)
				self.db.Put(append(pkrIndexPrefix, pk[:]...), data)
				// if end >= rtn.remoteNum {
				// 	break
				// }
				break
			}
		}
		return true
	})
	self.CheckNil()
}

func (self *SEROLight) fetchAndDecOuts(account *Account, pkrIndex uint64, start, end uint64) (rtn fetchReturn, err error) {

	pkrTypeMap, currentPkrsMap, pkrs := self.genPkrs(pkrIndex, account)

	param := []interface{}{pkrs, start}
	if (end != 0) {
		param = append(param, end)
	}

	sync := Sync{RpcHost: GetRpcHost(), Method: "light_getOutsByPKr", Params: param}
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
	rtn.lastBlockNumber = bor.CurrentNum;
	rtn.remoteNum = bor.CurrentNum
	if rtn.remoteNum > end {
		rtn.nextNum = end + 1
	} else {
		rtn.nextNum = bor.CurrentNum + 1
	}

	var hasResWithHashPkr = false
	var hasResWithOldPkr = false
	for _, blockOut := range blockOuts {
		datas := blockOut.Data
		for _, data := range datas {
			out := data.Out
			var pkr c_type.PKr

			if out.State.OS.Out_C != nil {
				pkr = out.State.OS.Out_C.PKr
			} else if out.State.OS.Out_O != nil {
				pkr = out.State.OS.Out_O.Addr
			} else if out.State.OS.Out_P != nil {
				pkr = out.State.OS.Out_P.PKr
			} else if out.State.OS.Out_Z != nil {
				pkr = out.State.OS.Out_Z.PKr
			}

			if currentPkrsMap[pkr] == 1 {
				rtn.again = true
				// gen min block Num
				if rtn.nextNum > blockOut.Num {
					rtn.nextNum = blockOut.Num
				}
			}

			if _, ok := self.useHashPkr.Load(account.pk); !ok {
				if pkrTypeMap[pkr] == PRK_TYPE_HASH {
					hasResWithHashPkr = true
				} else if pkrTypeMap[pkr] == PKR_TYPE_NUM {
					hasResWithOldPkr = true
				}
			}

			// dout := DecOuts([]txtool.Out{out}, &account.skr)[0]
			dout := flight.DecOut(account.tk, []txtool.Out{out})[0]

			dy, _ := json.Marshal(dout)
			fmt.Println("dout", string(dy[:]))

			key := PkKey{Pk: *account.pk, Num: blockOut.Num}
			utxo := Utxo{Pkr: pkr, Root: out.Root, Nils: dout.Nils, TxHash: out.State.TxHash, Num: out.State.Num, Asset: dout.Asset, IsZ: out.State.OS.Out_Z != nil, Out: out}

			// log.Info("DecOuts", "PK", base58.Encode(account.pk[:]), "root", common.Bytes2Hex(out.Root[:]), "currency", common.BytesToString(utxo.Asset.Tkn.Currency[:]), "value", utxo.Asset.Tkn.Value)
			if list, ok := utxosMap[key]; ok {
				utxosMap[key] = append(list, utxo)
			} else {
				utxosMap[key] = []Utxo{utxo}
			}

			// index base tx info
			txInfo := data.TxInfo
			txData, _ := rlp.EncodeToBytes(txInfo)
			self.db.Put(txHashKey(txInfo.TxHash[:]), txData)
		}

		// getBlock RPC
		// self.storeBlockInfo(blockOut.Num)
	}
	// if hash pkr return >0 and old pkr return = 0 ,set use hash pkr flag
	if _, ok := self.useHashPkr.Load(account.pk); !ok && (hasResWithHashPkr && !hasResWithOldPkr) {
		rtn.useHashPkr = true
	}

	rtn.utxoMap = utxosMap
	return rtn, nil
}

//
// func (self *SEROLight) storeBlockInfo(number uint64) {
//	sync := Sync{RpcHost: GetRpcHost(), Method: "sero_getBlockByNumber", Params: []interface{}{hexutil.EncodeUint64(number), false}}
//	resp, err := sync.Do()
//	if err != nil {
//		logex.Error("sero_getBlockByNumber request.do err: ", err)
//	} else {
//		if resp.Result != nil {
//			var b map[string]interface{}
//			err := json.Unmarshal(*resp.Result, &b)
//			if err != nil {
//				logex.Error("sero_getBlockByNumber json.Unmarshal: ", err)
//			} else {
//				blockEx := BlockEx{}
//				for key, value := range b {
//					if key == "number" {
//						numberHex := value.(string)
//						num, _ := hexutil.DecodeUint64(numberHex)
//						blockEx.BlockNumber = num
//					}
//					if key == "hash" {
//						blockEx.BlockHash = value.(string)
//					}
//					if key == "timestamp" {
//						timeHex := value.(string)
//						time, _ := hexutil.DecodeUint64(timeHex)
//						blockEx.Timestamp = time
//					}
//				}
//				if blockEx.BlockHash != "" {
//					bData, _ := rlp.EncodeToBytes(blockEx)
//					self.db.Put(blockIndex(number), bData)
//				}
//			}
//		}
//	}
// }

func (self *SEROLight) genPkrs(pkrIndex uint64, account *Account) (map[c_type.PKr]int8, map[c_type.PKr]int8, []string) {
	pkrTypeMap := map[c_type.PKr]int8{}
	// check loop again
	currentPkrsMap := map[c_type.PKr]int8{}
	var pkrs = []string{}
	pkrNum := int(1)
	// need append two main pkr
	pkrs = append(pkrs, base58.Encode(account.mainPkr[:]))
	if !c_superzk.IsSzkPKr(&account.mainPkr) {
		pkrs = append(pkrs, base58.Encode(account.mainOldPkr[:]))
	}

	if pkrIndex == 1 {
		currentPkrsMap[account.mainPkr] = 1
		currentPkrsMap[account.mainOldPkr] = 1
		pkrTypeMap[account.mainPkr] = PRK_TYPE_HASH

		if !c_superzk.IsSzkPKr(&account.mainPkr) {
			pkrTypeMap[account.mainOldPkr] = PKR_TYPE_NUM
		}
	}
	if pkrIndex > 5 {
		pkrNum = int(pkrIndex) - 5
	}
	for i := int(pkrIndex); i > pkrNum; i-- {
		pkrHash, _ := self.createPkrHash(account.tk, uint64(i))
		pkrs = append(pkrs, base58.Encode(pkrHash[:]))
		if _, ok := self.useHashPkr.Load(account.pk); !ok {
			pkrTypeMap[*pkrHash] = PRK_TYPE_HASH

			if !c_superzk.IsSzkPKr(&account.mainPkr) {
				pkrOld, _ := self.createPkr(account.tk, uint64(i))
				pkrs = append(pkrs, base58.Encode(pkrOld[:]))
				pkrTypeMap[*pkrOld] = PKR_TYPE_NUM
				if i == int(pkrIndex) {
					currentPkrsMap[*pkrOld] = 1
				}
			}
		}
		if i == int(pkrIndex) {
			currentPkrsMap[*pkrHash] = 1
		}
	}
	return pkrTypeMap, currentPkrsMap, pkrs
}

// if the currentpkr in the outs, again = true, then loop continue next Pkr
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

func (self *SEROLight) indexUtxo(utxosMap map[PkKey][]Utxo, batch serodb.Batch) (opsReturn map[string]string, err error) {
	ops := map[string]string{}
	for key, list := range utxosMap {
		roots := []c_type.Uint256{}
		for _, utxo := range list {
			data, err := rlp.EncodeToBytes(&utxo)
			if err != nil {
				return nil, err
			}
			// "ROOT" + root
			batch.Put(rootKey(utxo.Root), data)

			// "TXHASH" + PK + hash + root + outType
			batch.Put(indexTxKey(key.Pk, utxo.TxHash, utxo.Root, uint64(1)), data)

			// nil => root
			for _, Nil := range utxo.Nils {
				batch.Put(nilToRootKey(Nil), utxo.Root[:])
			}

			var pkKey []byte
			if utxo.Asset.Tkn != nil {
				// "PK" + PK + currency + root
				pkKey = utxoPkKey(key.Pk, utxo.Asset.Tkn.Currency[:], &utxo.Root)

			} else if utxo.Asset.Tkt != nil {
				// "PK" + PK + tkt + root
				pkKey = utxoPkKey(key.Pk, utxo.Asset.Tkt.Value[:], &utxo.Root)
			}
			// "PK" + PK + currency + root => 0
			ops[common.Bytes2Hex(pkKey)] = common.Bytes2Hex([]byte{0})

			// "NIL" + PK + tkt + root => "PK" + PK + currency + root
			for _, Nil := range utxo.Nils {
				// nilIdkey := nilIdKey(utxo.Nils)
				nilkey := nilKey(Nil)
				// "NIL" +nil/root => pkKey
				ops[common.Bytes2Hex(nilkey)] = common.Bytes2Hex(pkKey)
			}
			rootkey := nilKey(utxo.Root)

			// "NIL" +nil/root => pkKey
			ops[common.Bytes2Hex(rootkey)] = common.Bytes2Hex(pkKey)
			// ops[common.Bytes2Hex(nilIdkey)] = common.Bytes2Hex(encodeNumber(key.Num))
			roots = append(roots, utxo.Root)
			// log.Info("Index add", "PK", base58.Encode(key.PK[:]), "Nils", common.Bytes2Hex(utxo.Nils[:]), "root", common.Bytes2Hex(utxo.Root[:]), "Value", utxo.Asset.Tkn.Value)

			// self.genTxReceipt(utxo.TxHash, batch)
		}
		data, err := rlp.EncodeToBytes(roots)
		if err != nil {
			return nil, err
		}

		// utxo PK + at  => [roots]
		batch.Put(utxoKey(key.Num, key.Pk), data)
	}
	return ops, nil
}

// type NilValue struct {
//	Nil    c_type.Uint256
//	Num    uint64
//	TxHash c_type.Uint256
//	TxFee  big.Int
// }

type NilValue struct {
	Nil    c_type.Uint256
	Num    uint64
	TxHash c_type.Uint256
	TxInfo TxInfo
}

func (self *SEROLight) CheckNil() {

	iterator := self.db.NewIteratorWithPrefix(nilPrefix)
	// Nils := []keys.Uint256{}
	Nils := []string{}
	for iterator.Next() {
		key := iterator.Key()
		var Nil c_type.Uint256
		copy(Nil[:], key[3:])
		nilkey := nilKey(Nil)
		value, _ := self.db.Get(nilkey)
		if value != nil {
			Nils = append(Nils, hexutil.Encode(Nil[:]))
		}
	}

	sync := Sync{RpcHost: GetRpcHost(), Method: "light_checkNil", Params: []interface{}{Nils}}
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
		fmt.Println("nilvs:", nilvs)
		logex.Infof("light_checkNil result=[%d]", len(nilvs))
		if len(nilvs) > 0 {
			batch := self.db.NewBatch()
			for _, nilv := range nilvs {
				fmt.Println(nilv)
				var pk c_type.Uint512
				Nil := nilv.Nil

				value, _ := self.db.Get(nilKey(Nil))
				if value != nil {
					copy(pk[:], value[2:66])
					if account := self.getAccountByKey(pk); account != nil {
						account.isChanged = true
					}
					var root c_type.Uint256
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

					//TODO remove pending tx
					//fmt.Println("batch indexTxKey:",string(pk[:]), string(nilv.TxHash[:]), string(nilv.TxHash[:]),2)
					batch.Delete(indexTxKey(pk, nilv.TxHash, nilv.TxHash, uint64(2)))
					utxoI := Utxo{Root: root, TxHash: nilv.TxHash, Num: nilv.Num, Nils: []c_type.Uint256{nilv.Nil}, Asset: utxo.Asset, Pkr: utxo.Pkr}
					data, err := rlp.EncodeToBytes(&utxoI)
					if err !=nil{
						fmt.Println("EncodeToBytes err:",err)
						continue
					}
					batch.Put(indexTxKey(pk, nilv.TxHash, root, uint64(2)), data)


					txInfo := nilv.TxInfo
					txData, _ := rlp.EncodeToBytes(txInfo)
					batch.Put(txHashKey(nilv.TxHash[:]), txData)


					self.usedFlag.Delete(root)
				}
			}
			batch.Write()
		}
	}

}

func (self *SEROLight) genTxReceipt(txHash c_type.Uint256, batch serodb.Batch) {

	if *powReward.HashToUint256() == txHash || *posReward.HashToUint256() == txHash || *posMiner.HashToUint256() == txHash {
		logex.Info("txHash=", txHash, " is rewards Hash")
		return
	}

	var r *types.Receipt
	sync := Sync{RpcHost: GetRpcHost(), Method: "sero_getTransactionReceipt", Params: []interface{}{txHash}}
	resp, err := sync.Do()
	if err != nil {
		logex.Error("sero_getTransactionReceipt request.do err: ", err)
	} else {
		if resp.Result != nil {
			err := json.Unmarshal(*resp.Result, &r)
			if err != nil {
				logex.Error("sero_getTransactionReceipt json Unmarshal  err: ", err)
			} else {
				txReceipt := TxReceipt{
					Status:            r.Status,
					CumulativeGasUsed: r.CumulativeGasUsed,
					TxHash:            *r.TxHash.HashToUint256(),
					ContractAddress:   r.ContractAddress.Base58(),
					GasUsed:           r.GasUsed,
				}
				if r.PoolId != nil {
					txReceipt.PoolId = r.PoolId.String()
					txReceipt.ShareId = r.ShareId.String()
				}
				bData, err := rlp.EncodeToBytes(txReceipt)
				if err != nil {
					logex.Error("sero_getTransactionReceipt rlp.EncodeToBytes err: ", err)
				} else {
					err = batch.Put(txReceiptIndex(*r.TxHash.HashToUint256()), bData)
					logex.Error("batch.Put(txReceiptIndex,err :", r.TxHash.Hex(), err)
				}
			}
		}
	}
}

func (self *SEROLight) getAccountByKey(pk c_type.Uint512) *Account {
	if value, ok := self.accounts.Load(pk); ok {
		return value.(*Account)
	}
	return nil
}

func (self *SEROLight) getAccountByPk(pk c_type.Uint512) *Account {
	acc, err := self.accountManager.FindAccountByPk(pk)
	if err != nil {
		return nil
	}
	return self.getAccountByKey(acc.GetPk())

}

func (self *SEROLight) getAccountByPkr(pkr c_type.PKr) *Account {
	acc, err := self.accountManager.FindAccountByPkr(pkr)
	if err != nil {
		return nil
	}
	return self.getAccountByKey(acc.GetPk())
}

// "UTXO" + pk + number
func utxoKey(number uint64, pk c_type.Uint512) []byte {
	return append(utxoPrefix, append(pk[:], encodeNumber(number)...)...)
}

// utxoKey = PK + currency +root
func utxoPkKey(pk c_type.Uint512, currency []byte, root *c_type.Uint256) []byte {
	key := append(pkPrefix, pk[:]...)
	if len(currency) > 0 {
		key = append(key, currency...)
	}
	if root != nil {
		key = append(key, root[:]...)
	}
	return key
}

func (self *SEROLight) GetUtxoNum(pk c_type.Uint512) map[string]uint64 {
	if account := self.getAccountByKey(pk); account != nil {
		return account.utxoNums
	}
	return map[string]uint64{}
}

func (self *SEROLight) GetBalances(pk c_type.Uint512) (balances map[string]*big.Int) {
	if value, ok := self.accounts.Load(pk); ok {
		account := value.(*Account)
		// if account.isChanged {
		// } else {
		//	return account.balances
		// }

		prefix := append(pkPrefix, pk[:]...)
		iterator := self.db.NewIteratorWithPrefix(prefix)
		balances = map[string]*big.Int{}
		utxoNums := map[string]uint64{}
		for iterator.Next() {
			key := iterator.Key()
			var root c_type.Uint256
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
	}
	return
}

func (self *SEROLight) getUtxo(root c_type.Uint256) (utxo Utxo, e error) {
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

func (self *SEROLight) setZ() bool {
	data, err := self.db.Get(remoteNumKey)
	if err != nil {
		return false
	}
	num := decodeNumber(data[8:])
	if num >= useZNum {
		return true
	} else {
		return false
	}
}

func (self *SEROLight) commitTx(from, to, currency, passwd string, amount, gasprice *big.Int) (hash c_type.Uint256, err error) {

	fee := new(big.Int).Mul(big.NewInt(25000), gasprice)
	addr := address.StringToPk(from)

	var RefundTo *c_type.PKr
	ac := self.getAccountByPk(addr.ToUint512())
	pk := *ac.pk
	if ac == nil {
		logex.Errorf("account not found")
		return hash, fmt.Errorf("account not found")
	}

	if value, ok := self.pkrIndexMap.Load(pk); !ok {
		logex.Errorf("pkrIndexMap not store from accountKey")
		return hash, fmt.Errorf("account not found")
	} else {
		outReq := value.(outReq)
		RefundTo = &outReq.Pkr
	}

	account := accounts.Account{Address: addr}
	wallet, err := self.accountManager.Find(account)
	if err != nil {
		return hash, err
	}
	seed, err := wallet.GetSeedWithPassphrase(passwd)
	if err != nil {
		return hash, err
	}
	var toPkr c_type.PKr
	copy(toPkr[:], base58.Decode(to)[:])
	reception := self.genReceiption(currency, amount, toPkr)

	preTxParam := prepare.PreTxParam{}
	preTxParam.From = addr.ToUint512()
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = gasprice
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}
	preTxParam.Receptions = []prepare.Reception{reception}

	param, err := self.GenTx(preTxParam)
	self.needSzk(param)

	if err != nil {
		return hash, err
	}
	sk := c_superzk.Seed2Sk(seed.SeedToUint256())
	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return hash, err
	}
	hash = gtx.Hash
	sync := Sync{RpcHost: GetRpcHost(), Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return hash, err
	}

	utxoIn := Utxo{Pkr: toPkr, Root: hash, TxHash: hash, Fee: *fee}
	self.storePeddingUtxo(param, currency, amount, utxoIn, ac.pk)
	ac.isChanged = true

	return hash, nil
}

func (self *SEROLight) needSzk(param *txtool.GTxParam) {
	var need_szk = true
	data, err := self.db.Get(remoteNumKey)
	if err == nil {
		num := decodeNumber(data[:])
		fmt.Println("needSzk num:", num)
		fmt.Println("needSzk useZNum:", useZNum)
		if num >= useZNum {
			param.Z = &need_szk
		}
	}
}

func (self *SEROLight) storePeddingUtxo(param *txtool.GTxParam, currency string, amount *big.Int, utxoIn Utxo, pk *c_type.Uint512) {
	roots := []c_type.Uint256{}
	for _, in := range param.Ins {
		roots = append(roots, in.Out.Root)
		self.usedFlag.Store(in.Out.Root, 1)
	}
	tknc := assets.Token{Currency: utils.CurrencyToUint256(currency), Value: utils.U256(*amount)}
	assetc := assets.Asset{}
	assetc.Tkn = &tknc
	utxoIn.Asset = assetc
	dataIn, _ := rlp.EncodeToBytes(utxoIn)
	self.db.Put(indexTxKey(*pk, utxoIn.TxHash, utxoIn.TxHash, uint64(2)), dataIn)
}

func (self *SEROLight) genReceiption(currency string, amount *big.Int, toPkr c_type.PKr) prepare.Reception {
	tkn := assets.Token{Currency: utils.CurrencyToUint256(currency), Value: utils.U256(*amount)}
	asset := assets.NewAsset(&tkn, nil)
	reception := prepare.Reception{
		Addr:  toPkr,
		Asset: asset,
	}
	return reception
}

func (self *SEROLight) registerStakePool(from, vote, passwd string, feeRate uint32, amount, gasprice *big.Int) (hash c_type.Uint256, err error) {

	fee := new(big.Int).Mul(big.NewInt(25000), gasprice)
	fromAddress := address.StringToPk(from)
	fromPk := fromAddress.ToUint512()

	if len(base58.Decode(vote)) != 96 {
		return hash, fmt.Errorf("Invalid Vote Address ")
	}

	var RefundTo *c_type.PKr
	ac := self.getAccountByPk(fromPk)
	if ac != nil {
		RefundTo = &ac.mainPkr
	} else {
		return hash, errors.New("unknown account")
	}
	// check pk register pool
	poolId := crypto.Keccak256(ac.mainPkr[:])
	sync := Sync{RpcHost: GetRpcHost(), Method: "stake_poolState", Params: []interface{}{hexutil.Encode(poolId)}}
	_, err = sync.Do()
	if err != nil {
		if err.Error() != "stake pool not exists" {
			logex.Errorf("jsonRep err=[%s]", err.Error())
			return
		}
	} else {
		// amount > 0 is register . amount = 0 is modify
		if amount.Sign() > 0 {
			err = fmt.Errorf("stake pool exists")
			logex.Errorf("jsonRep err=[%s]", err.Error())
			return
		}
	}
	seed, err := ac.wallet.GetSeedWithPassphrase(passwd)
	if err != nil {
		return hash, err
	}
	var votePkr c_type.PKr
	if vote == "" {
		votePkr = ac.mainPkr
	} else {
		copy(votePkr[:], base58.Decode(vote)[:])
	}

	registerPool := stx.RegistPoolCmd{Value: utils.U256(*amount), Vote: votePkr, FeeRate: feeRate}
	preTxParam := prepare.PreTxParam{}
	preTxParam.From = fromPk
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = gasprice
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}
	preTxParam.Cmds = prepare.Cmds{RegistPool: &registerPool}

	param, err := self.GenTx(preTxParam)
	self.needSzk(param)

	if err != nil {
		return hash, err
	}

	sk := c_superzk.Seed2Sk(seed.SeedToUint256())

	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return hash, err
	}

	hash = gtx.Hash
	logex.Info("commit txhash: ", hash)
	sync = Sync{RpcHost: GetRpcHost(), Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return hash, err
	}

	utxoIn := Utxo{Pkr: votePkr, Root: hash, TxHash: hash, Fee: *fee}
	self.storePeddingUtxo(param, "SERO", amount, utxoIn, ac.pk)
	ac.isChanged = true

	return hash, nil
}

func (self *SEROLight) modifyStakePool(from, vote, passwd, idPkrStr string, feeRate uint32, amount, gasprice *big.Int) (hash c_type.Uint256, err error) {

	fee := new(big.Int).Mul(big.NewInt(25000), gasprice)
	fromAddress := address.StringToPk(from)
	fromPk := fromAddress.ToUint512()

	if len(base58.Decode(vote)) != 96 {
		return hash, fmt.Errorf("Invalid Vote Address ")
	}
	if idPkrStr != "" && len(base58.Decode(idPkrStr)) != 96 {
		return hash, fmt.Errorf("Invalid IdPkr ")
	}
	var idPkr c_type.PKr
	copy(idPkr[:], base58.Decode(idPkrStr)[:])

	var RefundTo *c_type.PKr
	ac := self.getAccountByPk(fromPk)
	if ac != nil {
		RefundTo = &ac.mainPkr
	} else {
		return hash, errors.New("unknown account")
	}

	// check pk register pool
	poolId := crypto.Keccak256(ac.mainPkr[:])
	sync := Sync{RpcHost: GetRpcHost(), Method: "stake_poolState", Params: []interface{}{hexutil.Encode(poolId)}}
	_, err = sync.Do()
	if err != nil {
		if err.Error() != "stake pool not exists" {
			logex.Errorf("jsonRep err=[%s]", err.Error())
			return
		}
	}
	seed, err := ac.wallet.GetSeedWithPassphrase(passwd)
	if err != nil {
		return hash, err
	}
	var votePkr c_type.PKr
	if vote == "" {
		votePkr = ac.mainPkr
	} else {
		copy(votePkr[:], base58.Decode(vote)[:])
	}

	registerPool := stx.RegistPoolCmd{Value: utils.U256(*amount), Vote: votePkr, FeeRate: feeRate}
	preTxParam := prepare.PreTxParam{}
	preTxParam.From = fromPk
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = gasprice
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}
	preTxParam.Cmds = prepare.Cmds{RegistPool: &registerPool}

	param, err := self.GenTx(preTxParam)
	self.needSzk(param)

	if err != nil {
		return hash, err
	}

	sk := c_superzk.Seed2Sk(seed.SeedToUint256())

	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return hash, err
	}

	hash = gtx.Hash
	logex.Info("commit txhash: ", hash)
	sync = Sync{RpcHost: GetRpcHost(), Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return hash, err
	}

	utxoIn := Utxo{Pkr: votePkr, Root: hash, TxHash: hash, Fee: *fee}
	self.storePeddingUtxo(param, "SERO", amount, utxoIn, ac.pk)
	ac.isChanged = true

	return hash, nil
}

func (self *SEROLight) closeStakePool(from, idPkrStr, passwd string) (hash c_type.Uint256, err error) {

	fee := new(big.Int).Mul(big.NewInt(25000), big.NewInt(1000000000))
	fromAddress := address.StringToPk(from)
	fromPk := fromAddress.ToUint512()

	if idPkrStr != "" && len(base58.Decode(idPkrStr)) != 96 {
		return hash, fmt.Errorf("Invalid IdPkr ")
	}
	var RefundTo *c_type.PKr
	var idPkr c_type.PKr
	copy(idPkr[:], base58.Decode(idPkrStr)[:])
	RefundTo = &idPkr

	ac := self.getAccountByPk(fromPk)
	if ac == nil {
		return hash, errors.New("unknown account")
	}
	// check pk register pool
	poolId := crypto.Keccak256(ac.mainPkr[:])
	sync := Sync{RpcHost: GetRpcHost(), Method: "stake_poolState", Params: []interface{}{hexutil.Encode(poolId)}}
	_, err = sync.Do()
	if err != nil {
		if err.Error() != "stake pool not exists" {
			logex.Errorf("jsonRep err=[%s]", err.Error())
			return
		}
	}
	seed, err := ac.wallet.GetSeedWithPassphrase(passwd)
	if err != nil {
		return hash, err
	}
	closePool := stx.ClosePoolCmd{}
	preTxParam := prepare.PreTxParam{}
	preTxParam.From = fromPk
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = big.NewInt(1000000000)
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}
	preTxParam.Cmds = prepare.Cmds{ClosePool: &closePool}

	param, err := self.GenTx(preTxParam)

	if err != nil {
		return hash, err
	}
	self.needSzk(param)

	sk := c_superzk.Seed2Sk(seed.SeedToUint256())
	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return hash, err
	}

	hash = gtx.Hash
	logex.Info("commit txhash: ", hash)
	sync = Sync{RpcHost: GetRpcHost(), Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return hash, err
	}

	utxoIn := Utxo{Pkr: ac.mainPkr, Root: hash, TxHash: hash, Fee: *fee}
	self.storePeddingUtxo(param, "SERO", big.NewInt(0), utxoIn, ac.pk)
	ac.isChanged = true

	return hash, nil
}

func (self *SEROLight) buyShare(from, vote, passwd, pool string, amount, gasprice *big.Int) (hash c_type.Uint256, err error) {

	fee := new(big.Int).Mul(big.NewInt(25000), gasprice)
	fromAddress := address.StringToPk(from)
	fromPk := fromAddress.ToUint512()

	var RefundTo *c_type.PKr
	ac := self.getAccountByPk(fromPk)
	if ac != nil {
		RefundTo = &ac.mainPkr
	} else {
		return hash, errors.New("unknown account")
	}

	seed, err := ac.wallet.GetSeedWithPassphrase(passwd)
	if err != nil {
		return hash, err
	}
	var votePkr c_type.PKr
	if len(vote) == 0 {
		votePkr = ac.mainPkr
	} else {
		copy(votePkr[:], base58.Decode(vote)[:])
	}
	poolId := common.HexToHash(pool)
	buyShareCmd := stx.BuyShareCmd{Value: utils.U256(*amount), Vote: votePkr, Pool: poolId.HashToUint256()}
	preTxParam := prepare.PreTxParam{}
	preTxParam.From = fromPk
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = gasprice
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}
	preTxParam.Cmds = prepare.Cmds{BuyShare: &buyShareCmd}
	param, err := self.GenTx(preTxParam)
	self.needSzk(param)

	if err != nil {
		return hash, err
	}
	sk := c_superzk.Seed2Sk(seed.SeedToUint256())
	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return hash, err
	}
	hash = gtx.Hash
	logex.Info("commit txhash: ", hash)
	sync := Sync{RpcHost: GetRpcHost(), Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return hash, err
	}

	utxoIn := Utxo{Pkr: votePkr, Root: hash, TxHash: hash, Fee: *fee}
	self.storePeddingUtxo(param, "SERO", amount, utxoIn, ac.pk)
	ac.isChanged = true

	return hash, nil
}

func (self *SEROLight) getDecimal(currency string) uint64 {
	if decimalData, err := self.db.Get(append(decimalPrefix, []byte(currency)[:]...)); err != nil {
		if decimalData == nil {
			sync := Sync{RpcHost: GetRpcHost(), Method: "sero_getDecimal", Params: []interface{}{currency}}
			if jsonResp, err := sync.Do(); err != nil {
				return 0
			} else {
				var decimalStr string
				if err = json.Unmarshal(*jsonResp.Result, &decimalStr); err != nil {
					logex.Error("json.Unmarshal err=[%s]", err.Error())
					return 0
				}
				decimalStr = decimalStr[2:]
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
	self.pkrIndexMap.Range(func(key, value interface{}) bool {
		data := value.(outReq)
		if number < data.Num {
			number = data.Num
		}
		return true
	})
	return number
}

func (self *SEROLight) getLatestPKrs(pk c_type.Uint512) (pais []pkrAndIndex) {
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
		var pkr c_type.PKr
		copy(pkr[:], value[:])
		pai.pkr = pkr
		pais = append(pais, pai)
		count++
	}
	return pais
}

func (self *SEROLight) DeployContractTx(ctq ContractTxReq, password string) (txHash string, err error) {

	gasPrice, err := NewBigIntFromString(ctq.GasPrice, 10)
	if err != nil {
		return "", err
	} else {
		if gasPrice.Sign() < 0 {
			return "", fmt.Errorf("gasPrice < 0")
		}
	}
	gas, err := NewBigIntFromString(ctq.Gas, 10)
	if err != nil {
		return "", err
	} else {
		if gas.Sign() < 0 {
			return "", fmt.Errorf("gas < 0")
		}
	}
	amount, err := NewBigIntFromString(ctq.Value, 10)
	if err != nil {
		return "", err
	} else {
		if amount.Sign() < 0 {
			return "", fmt.Errorf("amount < 0")
		}
	}
	fromAddress := address.StringToPk(ctq.From)
	fromPk := fromAddress.ToUint512()
	var RefundTo *c_type.PKr
	ac := self.getAccountByPk(fromPk)
	if ac == nil {
		logex.Errorf("account not found")
		return txHash, fmt.Errorf("account not found")
	}
	// random := keys.RandUint128()
	// copy(random[:], ctq.Data[:16])
	// fromPkr := self.genPkrContract(fromPk, random)
	// RefundTo = &fromPkr
	RefundTo = &ac.mainPkr

	seed, err := ac.wallet.GetSeedWithPassphrase(password)
	if err != nil {
		return txHash, err
	}

	fee := big.NewInt(0).Mul(gas, gasPrice)
	preTxParam := prepare.PreTxParam{}
	preTxParam.From = fromPk
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = gasPrice
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}
	preTxParam.Cmds = prepare.Cmds{
		Contract: &stx.ContractCmd{
			Data: ctq.Data,
			Asset: assets.Asset{
				Tkn: &assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*amount)},
			},
		},
	}

	param, err := self.GenTx(preTxParam)
	self.needSzk(param)

	if err != nil {
		return txHash, err
	}
	sk := c_superzk.Seed2Sk(seed.SeedToUint256())
	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return txHash, err
	}

	txHash = hexutil.Encode(gtx.Hash[:])

	sync := Sync{RpcHost: GetRpcHost(), Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return txHash, err
	}

	utxoIn := Utxo{Pkr: *RefundTo, Root: gtx.Hash, TxHash: gtx.Hash, Fee: *fee}
	self.storePeddingUtxo(param, "SERO", amount, utxoIn, ac.pk)
	ac.isChanged = true

	ctq.Token.TxHash = txHash
	if data, err := rlp.EncodeToBytes(ctq.Token); err == nil {
		self.db.Put(append(tokenPrefix[:], []byte(txHash)[:]...), data[:])
	}

	return txHash, nil
}

func (self *SEROLight) ExecuteContractTx(ctq ContractTxReq, password string) (txHash string, err error) {

	gasPrice, err := NewBigIntFromString(ctq.GasPrice[2:], 16)
	if err != nil {
		return "", err
	} else {
		if gasPrice.Sign() < 0 {
			return "", fmt.Errorf("gasPrice < 0")
		}
	}
	gas, err := NewBigIntFromString(ctq.Gas[2:], 16)
	if err != nil {
		return "", err
	} else {
		if gas.Sign() < 0 {
			return "", fmt.Errorf("gas < 0")
		}
	}
	amount := big.NewInt(0)
	if ctq.Value != "" {
		amount, err = NewBigIntFromString(ctq.Value[2:], 16)
		if err != nil {
			return "", err
		} else {
			if amount.Sign() < 0 {
				return "", fmt.Errorf("amount < 0")
			}
		}
	}
	var ac *Account
	add, err := account.NewAddressByString(ctq.From)
	if err != nil {
		return txHash, err
	}
	fromByte := add.Bytes
	if len(fromByte) == 96 {
		pkr := c_type.PKr{}
		copy(pkr[:], fromByte[:])
		ac = self.getAccountByPkr(pkr)
	} else if len(fromByte) == 64 {
		pk := c_type.Uint512{}
		copy(pk[:], fromByte[:])
		ac = self.getAccountByPk(pk)

	}
	var RefundTo *c_type.PKr
	if ac == nil {
		logex.Errorf("account not found")
		return txHash, fmt.Errorf("account not found")
	}

	// random := keys.RandUint128()
	// copy(random[:], ctq.Data[:16])
	// fromPkr := self.genPkrContract(fromPk, random)
	// RefundTo = &fromPkr

	RefundTo = &ac.mainPkr

	seed, err := ac.wallet.GetSeedWithPassphrase(password)
	if err != nil {
		return txHash, err
	}
	var toPkr c_type.PKr
	copy(toPkr[:], base58.Decode(ctq.To)[:])

	cy := "SERO"
	if ctq.Currency != "" {
		cy = ctq.Currency
	}
	fee := big.NewInt(0).Mul(gas, gasPrice)
	preTxParam := prepare.PreTxParam{}
	preTxParam.From = *ac.pk
	preTxParam.RefundTo = RefundTo
	preTxParam.GasPrice = gasPrice
	preTxParam.Fee = assets.Token{Currency: utils.CurrencyToUint256("SERO"), Value: utils.U256(*fee)}
	preTxParam.Cmds = prepare.Cmds{
		Contract: &stx.ContractCmd{
			Data: ctq.Data,
			To:   &toPkr,
			Asset: assets.Asset{
				Tkn: &assets.Token{Currency: utils.CurrencyToUint256(cy), Value: utils.U256(*amount)},
			},
		},
	}

	param, err := self.GenTx(preTxParam)
	self.needSzk(param)
	if err != nil {
		return txHash, err
	}
	sk := c_superzk.Seed2Sk(seed.SeedToUint256())
	gtx, err := flight.SignTx(&sk, param)
	if err != nil {
		return txHash, err
	}

	txHash = hexutil.Encode(gtx.Hash[:])

	sync := Sync{RpcHost: GetRpcHost(), Method: "sero_commitTx", Params: []interface{}{gtx}}
	if _, err := sync.Do(); err != nil {
		return txHash, err
	}

	utxoIn := Utxo{Pkr: *RefundTo, Root: gtx.Hash, TxHash: gtx.Hash, Fee: *fee}
	self.storePeddingUtxo(param, "SERO", amount, utxoIn, ac.pk)
	ac.isChanged = true

	return txHash, nil
}

func (self *SEROLight) getTokens() ([]TokenReq, error) {
	prefix := append(tokenPrefix)
	iterator := self.db.NewIteratorWithPrefix(prefix)
	tokens := []TokenReq{}
	for iterator.Next() {
		// key := iterator.Key()
		value := iterator.Value()
		token := TokenReq{}
		err := rlp.DecodeBytes(value, &token)
		if err != nil {
			return nil, err
		}
		// //get Transaction Receipt
		if token.TxHash != "" && token.Symbol != "" {
			sync := Sync{RpcHost: GetRpcHost(), Method: "sero_currencyToContractAddress", Params: []interface{}{token.Symbol}}
			jsonResp, err := sync.Do()
			if err == nil {
				var ctrtAddr string
				json.Unmarshal(*jsonResp.Result, &ctrtAddr)
				token.ContractAddress = string(ctrtAddr[:])
				token.TxHash = ""
				data, err := rlp.EncodeToBytes(token)
				if err == nil {
					self.db.Put(append(tokenPrefix, []byte(token.ContractAddress)[:]...), data[:])
					self.db.Delete(append(tokenPrefix[:], []byte(token.TxHash)[:]...))
				}
			}
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

type ContractTxReq struct {
	From     string        `json:"from"`
	To       string        `json:"to"`
	Value    string        `json:"value"`
	GasPrice string        `json:"gas_price"`
	Gas      string        `json:"gas"`
	Currency string        `json:"cy"`
	Data     hexutil.Bytes `json:"data"`
	Token    TokenReq      `json:"token"`
}

type TokenReq struct {
	TxHash          string
	ContractAddress string
	Name            string
	Symbol          string
	Decimal         uint8
	Total           string
}
