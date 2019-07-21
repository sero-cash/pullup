package app

import "github.com/sero-cash/go-czero-import/keys"

var (
	numPrefix     = []byte("NUM")
	pkPrefix      = []byte("PK")
	pkrPrefix     = []byte("PKR")
	utxoPrefix    = []byte("UTXO")
	rootPrefix    = []byte("ROOT")
	allRootPrefix = []byte("ALLROOT")
	nilPrefix     = []byte("NIL")
	nilIdPrefix   = []byte("NILID")
	syncNilKEY    = []byte("SYNCNILNUM")

	nilRootPrefix   = []byte("NOILTOROOT")
	peddingTxPrefix = []byte("PEDDINGTX")
	decimalPrefix   = []byte("DECIMAL")
	pkrIndexPrefix =  []byte("PKRINDEX")

	hostKey =  []byte("HOST")
)

// PKR + PK + r
func pkrKey(pk keys.Uint512, r keys.Uint256) []byte {
	key := append(pkrPrefix, pk[:]...)
	key = append(key, r[:]...)
	return key
}

func numKey(pk keys.Uint512) []byte {
	return append(numPrefix, pk[:]...)
}

func nilKey(nil keys.Uint256) []byte {
	return append(nilPrefix, nil[:]...)
}

func rootKey(root keys.Uint256) []byte {
	return append(rootPrefix, root[:]...)
}

type PkKey struct {
	PK  keys.Uint512
	Num uint64
}

type pkrAndIndex struct {
	pkr   keys.PKr
	index uint64
}


func nilToRootKey(nil keys.Uint256) []byte {
	return append(nilRootPrefix, nil[:]...)
}

func penddingTxKey(pk keys.Uint512, hash keys.Uint256) []byte {
	key := append(peddingTxPrefix, pk[:]...)
	return append(key, hash[:]...)
}