package app

import (
	"math/big"

	"github.com/sero-cash/go-czero-import/c_type"
	"github.com/sero-cash/go-sero/zero/txs/assets"
	"github.com/sero-cash/go-sero/zero/txtool"
)

type BlockOutResp struct {
	CurrentNum uint64
	BlockOuts  []BlockOut
}

type BlockOut struct {
	Num  uint64
	Outs []txtool.Out
}

type BlockInfo struct {
	Num  uint64
	Hash c_type.Uint256
	Ins  []c_type.Uint256
	Outs []Utxo
}

type Utxo struct {
	Pkr       c_type.PKr
	Root      c_type.Uint256
	TxHash    c_type.Uint256
	Nils      []c_type.Uint256
	Num       uint64
	Asset     assets.Asset
	IsZ       bool
	flag      int
	Out       txtool.Out
	Fee       big.Int
	Timestamp uint64
}

// stake pool

type StakePool struct {
	ChoicedNum  string `json:"choicedNum"`
	Closed      bool   `json:"closed"`
	CreateAt    string `json:"createAt"`
	ExpireNum   string `json:"expireNum"`
	Fee         string `json:"fee"`
	Id          string `json:"id"`
	IdPkr       string `json:"idPkr"`
	LastPayTime string `json:"lastPayTime"`
	MissedNum   string `json:"missedNum"`
	Own         string `json:"own"`
	Profit      string `json:"profit"`
	ShareNum    string `json:"shareNum"`
	Tx          string `json:"tx"`
	VoteAddress string `json:"voteAddress"`
	WishVoteNum string `json:"wishVoteNum"`
}

type TxReceipt struct {
	// Consensus fields
	Status            uint64
	CumulativeGasUsed uint64

	// Implementation fields (don't reorder!)
	TxHash          c_type.Uint256
	ContractAddress string
	GasUsed         uint64

	//Staking
	PoolId  string
	ShareId string

	BlockNumber uint64
	BlockHash   string
}

type BlockEx struct {
	BlockNumber uint64
	Timestamp   uint64
	BlockHash   string
}
