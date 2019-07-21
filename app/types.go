package app

import (
	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/zero/txs/assets"
	"github.com/sero-cash/go-sero/zero/txtool"
	"math/big"
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
	Hash keys.Uint256
	Ins  []keys.Uint256
	Outs []Utxo
}

type Utxo struct {
	Pkr    keys.PKr
	Root   keys.Uint256
	TxHash keys.Uint256
	Nil    keys.Uint256
	Num    uint64
	Asset  assets.Asset
	IsZ    bool
	flag   int
	Out    txtool.Out
	Fee    big.Int
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
