package app

import (
	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/common"
	"github.com/sero-cash/go-sero/common/hexutil"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/rlp"
	"math/big"
	"time"
)

type Transaction struct {
	Type      uint64
	Hash      keys.Uint256
	Block     uint64
	PK        keys.Uint512
	To        keys.PKr
	Currency  keys.Uint256
	Amount    *big.Int
	Fee       *big.Int
	Receipt   TxReceipt
	Timestamp uint64
}

var txHashPrefix = []byte("TXHASH")

//"TXHASH"+PK+hash+root+outType = utxo
func indexTxKey(pk keys.Uint512, hash keys.Uint256, root keys.Uint256, outType uint64) []byte {
	key := append(txHashPrefix, pk[:]...)
	key = append(key, hash[:]...)
	key = append(key, root[:]...)
	return append(key, encodeNumber(outType)...)
}

var (
	powReward = common.BytesToHash([]byte{1})
	posReward = common.BytesToHash([]byte{2})
	posMiner  = common.BytesToHash([]byte{3})
)

func (self *SEROLight) findTx(pk keys.Uint512, pageCount uint64) (map[string]Transaction, error) {
	prefix := append(txHashPrefix, pk[:]...)
	iterator := self.db.NewIteratorWithPrefix(prefix)
	txMap := map[string]Transaction{}
	i:=uint64(0)
	for iterator.Next() {
		i++
		key := iterator.Key()
		value := iterator.Value()
		doutroot := keys.Uint256{}
		douthash := keys.Uint256{}
		copy(douthash[:], key[70:102])
		copy(doutroot[:], key[102:134])
		outType := decodeNumber(key[134:142])
		utxo := Utxo{}
		rlp.DecodeBytes(value, &utxo)

		ukeyb := douthash[:]
		if *powReward.HashToUint256() == douthash || *posReward.HashToUint256() == douthash || *posMiner.HashToUint256() == douthash {
			ukeyb = append(encodeNumber(utxo.Num), utxo.Pkr[:]...)
		}

		if utxo.Asset.Tkn != nil {
			ukeyb = append(ukeyb[:], utxo.Asset.Tkn.Currency[:]...)
		}

		ukey := hexutil.Encode(ukeyb[:])
		if outType == 2 {
			if tx, ok := txMap[ukey]; ok {
				tx.Type = 2
			}
		}
		if utxo.Asset.Tkn != nil {
			amount := utxo.Asset.Tkn.Value.ToIntRef()
			fee := &utxo.Fee
			if outType == 2 {
				amount = big.NewInt(0).Mul(amount, big.NewInt(-1))
			}
			if tx, ok := txMap[ukey]; ok {
				tx.Amount = big.NewInt(0).Add(tx.Amount, amount)
				if outType == 2 {
					tx.Fee = fee
					tx.To = utxo.Pkr
				}
				txMap[ukey] = tx
			} else {

				tx = Transaction{Type: outType, Hash: douthash, Block: utxo.Num, PK: pk, To: utxo.Pkr, Amount: amount, Currency: utxo.Asset.Tkn.Currency, Fee: fee}
				rData, err := self.db.Get(txReceiptIndex(douthash))
				if err != nil {
					logex.Error("txHash not indexed, hash: ", douthash, err)
				} else {
					var r TxReceipt
					err := rlp.DecodeBytes(rData, &r)
					if err != nil {
						logex.Error("txReceipt rlp.decode err: ", err)
					} else {
						tx.Receipt = r
					}
				}
				if utxo.Num == 0{
					tx.Timestamp = uint64(time.Now().Unix())
				}else{
					bData, err := self.db.Get(blockIndex(utxo.Num))
					if err != nil {
						logex.Error("block not indexed, hash: ", utxo.Num, err)
					} else {
						var b BlockEx
						err := rlp.DecodeBytes(bData, &b)
						if err != nil {
							logex.Error("rlp.decode err: ", err)
						} else {
							tx.Receipt.BlockHash = b.BlockHash
							tx.Receipt.BlockNumber = b.BlockNumber
							tx.Timestamp = b.Timestamp+i
						}
					}
				}

				txMap[ukey] = tx
			}
		}
	}
	return txMap, nil
}
