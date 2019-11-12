package app

import (
	"github.com/sero-cash/go-czero-import/c_type"
	"github.com/sero-cash/go-sero/common"
	"github.com/sero-cash/go-sero/common/hexutil"
	"github.com/sero-cash/go-sero/pullup/common/logex"
	"github.com/sero-cash/go-sero/rlp"
	"math/big"
	"time"
)

type Transaction struct {
	Type      uint64
	Hash      c_type.Uint256
	Block     uint64
	PK        c_type.Uint512
	To        c_type.PKr
	Currency  c_type.Uint256
	Amount    *big.Int
	Fee       *big.Int
	Receipt   TxReceipt
	Timestamp uint64
}

var txHashPrefix = []byte("TXHASH")

//"TXHASH"+PK+hash+root+outType = utxo
func indexTxKey(pk c_type.Uint512, hash c_type.Uint256, root c_type.Uint256, outType uint64) []byte {
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

func (self *SEROLight) findTx(pk c_type.Uint512, pageCount uint64) (map[string]Transaction, error) {
	prefix := append(txHashPrefix, pk[:]...)
	iterator := self.db.NewIteratorWithPrefix(prefix)
	txMap := map[string]Transaction{}
	//i := uint64(0)
	for iterator.Next() {
		//i++
		key := iterator.Key()
		value := iterator.Value()
		doutroot := c_type.Uint256{}
		douthash := c_type.Uint256{}
		copy(douthash[:], key[70:102])
		copy(doutroot[:], key[102:134])
		outType := decodeNumber(key[134:142])

		//fmt.Println(hexutil.Encode(doutroot[:]),"-----",hexutil.Encode(douthash[:]))

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

			if outType == 2 {
				amount = big.NewInt(0).Mul(amount, big.NewInt(-1))
			}
			if tx, ok := txMap[ukey]; ok {
				tx.Amount = big.NewInt(0).Add(tx.Amount, amount)
				//if outType == 2 {
					//tx.Fee = fee
					//tx.To = utxo.Pkr
				//}
				txMap[ukey] = tx
			} else {
				tx = Transaction{Type: outType, Hash: douthash, Block: utxo.Num, PK: pk, To: utxo.Pkr, Amount: amount, Currency: utxo.Asset.Tkn.Currency}

				rData, err := self.db.Get(txHashKey(douthash[:]))
				if err != nil {
					if *powReward.HashToUint256() == douthash || *posReward.HashToUint256() == douthash || *posMiner.HashToUint256() == douthash {
					}else{
						batch := self.db.NewBatch()
						self.genTxReceipt(douthash, batch)
						batch.Write()
					}
					logex.Error("txHash not indexed, hash: ", douthash, err)
					tx.Timestamp = uint64(time.Now().Unix())
				}else{
					var txInfo TxInfo
					rlp.DecodeBytes(rData,&txInfo)
					tx.Receipt.TxHash = txInfo.TxHash
					tx.Receipt.BlockHash = txInfo.BlockHash.String()
					tx.Receipt.BlockNumber = txInfo.Num
					tx.Receipt.GasUsed=txInfo.GasUsed
					tx.Timestamp = txInfo.Time.Uint64()
					tx.Fee = big.NewInt(0).Mul(big.NewInt(int64(txInfo.GasUsed)),&txInfo.GasPrice)
				}
				txMap[ukey] = tx
			}
		}
	}
	return txMap, nil
}
