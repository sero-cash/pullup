package app

import (
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
)

const (
	PRK_TYPE_HASH int8 = 2
	PKR_TYPE_NUM  int8 = 1
)

func txHashKey(hash []byte) []byte  {
	return append(hashPrefix,hash[:]...)
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

func penddingTxKey(pk c_type.Uint512, hash c_type.Uint256) []byte {
	key := append(peddingTxPrefix, pk[:]...)
	return append(key, hash[:]...)
}
