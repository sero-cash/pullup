package app

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/sero-cash/go-czero-import/c_type"
	"github.com/sero-cash/go-sero/common"
	"github.com/sero-cash/go-sero/common/hexutil"
	"github.com/sero-cash/go-sero/zero/localdb"
	"github.com/sero-cash/go-sero/zero/txs/assets"
	"github.com/sero-cash/go-sero/zero/txtool"
	"github.com/sero-cash/go-sero/zero/txtool/prepare"
)

func (self *SEROLight) GenTx(param prepare.PreTxParam) (txParam *txtool.GTxParam, e error) {
	txParam, e = prepare.GenTxParam(&param, self, self)
	return
}

//===== TxParamGenerator interface impl
func (self *SEROLight) findUtxos(pk *c_type.Uint512, currency string, amount *big.Int) (utxos []Utxo, remain *big.Int) {
	remain = new(big.Int).Set(amount)

	currency = strings.ToUpper(currency)
	prefix := append(pkPrefix, append(pk[:], common.LeftPadBytes([]byte(currency), 32)...)...)
	iterator := self.db.NewIteratorWithPrefix(prefix)

	for iterator.Next() {
		key := iterator.Key()
		var root c_type.Uint256
		copy(root[:], key[98:130])

		if utxo, err := self.getUtxo(root); err == nil {
			if utxo.Asset.Tkn != nil {
				if _, ok := self.usedFlag.Load(utxo.Root); !ok {
					utxos = append(utxos, utxo)
					remain.Sub(remain, utxo.Asset.Tkn.Value.ToIntRef())
					if remain.Sign() <= 0 {
						break
					}
				}
			}
		}
	}
	return
}

func (self *SEROLight) findUtxosByTicket(pk *c_type.Uint512, tickets []assets.Ticket) (utxos []Utxo, remain map[c_type.Uint256]c_type.Uint256) {
	remain = map[c_type.Uint256]c_type.Uint256{}
	for _, t := range tickets {
		remain[t.Value] = t.Category
		prefix := append(pkPrefix, append(pk[:], t.Value[:]...)...)
		iterator := self.db.NewIteratorWithPrefix(prefix)
		if iterator.Next() {
			key := iterator.Key()
			var root c_type.Uint256
			copy(root[:], key[98:130])

			if utxo, err := self.getUtxo(root); err == nil {
				if utxo.Asset.Tkt != nil && utxo.Asset.Tkt.Category == t.Category {
					if _, ok := self.usedFlag.Load(utxo.Root); !ok {
						utxos = append(utxos, utxo)
						delete(remain, t.Value)
					}
				}
			}
		}
	}
	return
}

func (self *SEROLight) FindRoots(pk *c_type.Uint512, currency string, amount *big.Int) (roots prepare.Utxos, remain big.Int) {
	utxos, r := self.findUtxos(pk, currency, amount)
	for _, utxo := range utxos {
		roots = append(roots, prepare.Utxo{utxo.Root, utxo.Asset})
	}
	remain = *r
	return
}

// tickets map[keys.Uint256]keys.Uint256) (utxos []Utxo, remain map[keys.Uint256]keys.Uint256)
func (self *SEROLight) FindRootsByTicket(pk *c_type.Uint512, tickets []assets.Ticket) (roots prepare.Utxos, remain map[c_type.Uint256]c_type.Uint256) {
	utxos, remain := self.findUtxosByTicket(pk, tickets)
	for _, utxo := range utxos {
		roots = append(roots, prepare.Utxo{utxo.Root, utxo.Asset})
	}
	return
}

func (self *SEROLight) DefaultRefundTo(from *c_type.Uint512) (ret *c_type.PKr) {
	if value, ok := self.accounts.Load(from); ok {
		account := value.(*Account)
		return &account.mainPkr
	} else {
		return nil
	}
}

func (self *SEROLight) GetRoot(root *c_type.Uint256) (utxos *prepare.Utxo) {
	if u, e := self.getUtxo(*root); e != nil {
		return nil
	} else {
		return &prepare.Utxo{u.Root, u.Asset}
	}
}

//===== TxParamState interface impl

func (self *SEROLight) GetAnchor(roots []c_type.Uint256) ([]txtool.Witness, error) {
	params := []string{}
	for _, each := range roots {
		params = append(params, hexutil.Encode(each[:]))
	}
	sync := Sync{RpcHost: GetRpcHost(), Method: "sero_getAnchor", Params: []interface{}{params}}
	rpcResp, err := sync.Do()
	if err != nil {
		return nil, err
	}
	if rpcResp.Result != nil {
		var witnesses []txtool.Witness
		err = json.Unmarshal(*rpcResp.Result, &witnesses)
		return witnesses, err
	}
	return nil, nil
}

func (self *SEROLight) GetOut(root *c_type.Uint256) (out *localdb.RootState) {
	if u, e := self.getUtxo(*root); e != nil {
		return nil
	} else {
		return &u.Out.State
	}
}
func (self *SEROLight) GetPkgById(id *c_type.Uint256) (ret *localdb.ZPkg) {
	return nil
}

func (self *SEROLight) GetSeroGasLimit(to *common.Address, tfee *assets.Token, gasPrice *big.Int) (gaslimit uint64, e error) {
	return big.NewInt(0).Div(tfee.Value.ToInt(), gasPrice).Uint64(), nil
}
