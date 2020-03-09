package app

import (
	"encoding/binary"
	"github.com/sero-cash/go-czero-import/c_type"
)

var (
	numPrefix         = []byte("NUM")
	pkPrefix          = []byte("PK")
	pkrPrefix         = []byte("PKR")
	utxoPrefix        = []byte("UTXO")
	rootPrefix        = []byte("ROOT")
	allRootPrefix     = []byte("ALLROOT")
	nilPrefix         = []byte("NIL")
	nilIdPrefix       = []byte("NILID")
	syncNilKEY        = []byte("SYNCNILNUM")
	nilRootPrefix     = []byte("NOILTOROOT")
	peddingTxPrefix   = []byte("PEDDINGTX")
	decimalPrefix     = []byte("DECIMAL")
	pkrIndexPrefix    = []byte("PKRINDEX")
	hostKey           = []byte("HOST")
	VersonKey         = []byte("VERSION")
	onlyUseHashPkrKey = []byte("USEHASHPKR")

	tokenPrefix     = []byte("Token")
	txReceiptPrefix = []byte("TXRECEIPT")
	blockPrefix     = []byte("BLOCK")

	dappPrefix = []byte("DAPPS")
	hashPrefix = []byte("HASH")
	remoteNumKey= []byte("REMOTENUM")
	txHashPrefix = []byte("TXHASH")

	txPendingHashPrefix = []byte("TXPENDINGHASH")
)

const (
	PRK_TYPE_HASH int8 = 2
	PKR_TYPE_NUM  int8 = 1
)

//"TXHASH"+PK+hash+root+outType = utxo
func indexTxKey(pk c_type.Uint512,num uint64, hash c_type.Uint256, root c_type.Uint256, outType uint64) []byte {
	maxNum := uint64(9999999999) - num
	key := append(txHashPrefix, pk[:]...)
	key = append(key, encodeNumber(maxNum)...)
	key = append(key, hash[:]...)
	key = append(key, root[:]...)
	return append(key, encodeNumber(outType)...)
}

func txPendingHashKey(pk c_type.Uint512,hash c_type.Uint256, now uint64) []byte {
	key := append(txPendingHashPrefix, pk[:]...)
	key = append(key, hash[:]...)
	enc := make([]byte, 16)
	binary.BigEndian.PutUint64(enc, now)
	return append(key, enc...)
}

func txPendingHashKey2(pk c_type.Uint512,hash c_type.Uint256) []byte {
	key := append(txPendingHashPrefix, pk[:]...)
	key = append(key, hash[:]...)
	return key
}

func txHashKey(hash []byte,num uint64) []byte  {
	var key = append(hashPrefix,hash[:]...);
	return append(key,encodeNumber(num)...)
}
func dappKey(dappId string) []byte {
	key := append(dappPrefix, dappId[:]...)
	return key
}

// PKR + PK + r
func pkrKey(pk c_type.Uint512, r c_type.Uint256) []byte {
	key := append(pkrPrefix, pk[:]...)
	key = append(key, r[:]...)
	return key
}

func txReceiptIndex(txHash c_type.Uint256) []byte {
	return append(txReceiptPrefix, txHash[:]...)
}

func blockIndex(num uint64) []byte {
	return append(blockPrefix, encodeNumber(num)...)
}

func nilKey(nil c_type.Uint256) []byte {
	return append(nilPrefix, nil[:]...)
}

func rootKey(root c_type.Uint256) []byte {
	return append(rootPrefix, root[:]...)
}

type PkKey struct {
	Pk  c_type.Uint512
	Num uint64
}

type pkrAndIndex struct {
	pkr   c_type.PKr
	index uint64
}

func nilToRootKey(nil c_type.Uint256) []byte {
	return append(nilRootPrefix, nil[:]...)
}
