package types

import (
	"fs.video/blockchain/util"
)

const (
	RelationKeyDeflationVoteInforResult  = "DeflationVoteInforResult"
	RelationKeyDeflationVotePowerResult  = "DeflationVotePowerResult"
	RelationKeyTradeInforResult          = "TradeInforResult"
	RelationKeyBuySourceInforResult      = "BuySourceInforResult"
	RelationKeySourceSizeInforResult     = "SourceSizeforResult"
	RelationKeyMortgageRedeemResult      = "MortgageRedeemResult"
	RelationKeyCopyrightMoveResult       = "CopyrightMoveResult"
	RelationKeySpaceMinerResult          = "SpaceMinerResult"
	RelationKeyDeflationVoteResult       = "DeflationVoteResult"
	RelationKeyCopyrightApproveResult    = "CopyrightApproveResult"
	RelationKeyCopyrightVoteRedeemResult = "CopyrightVoteRedeemResult"
	RelationKeyCopyrightVoteResult       = "CopyrightVoteResult"
)

type BlockRelationData interface {
	Type() string
	Marshal() []byte
}

type BlockRDS struct {
	List []BlockRD `json:"list"`
}

type BlockRD struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

func (this *BlockRDS) Add(newList ...BlockRD) {
	this.List = append(this.List, newList...)
}

func (this *BlockRD) Unmarshal(xx interface{}) error {
	return util.Json.Unmarshal(this.Data, xx)
}

//
func NewBlockRD(data BlockRelationData) BlockRD {
	return BlockRD{Type: data.Type(), Data: data.Marshal()}
}

func NewBlockRDS() BlockRDS {
	rds := BlockRDS{}
	rds.List = []BlockRD{}
	return rds
}

type DeflationAccountVote struct {
	Account   string `json:"account"`
	Power     string `json:"power"`
	VoteIndex string `json:"vote_index"`
}

func (m DeflationAccountVote) Type() string { return RelationKeyDeflationVotePowerResult }

func (m DeflationAccountVote) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type DeflationRateVoteInfor struct {
	VoteOption   string `json:"vote_option"`
	VoteTitle    string `json:"vote_title"`
	BeginNum     int64  `json:"begin_num"`
	EndNum       int64  `json:"end_num"`
	VoteIndex    string `json:"vote_index" `
	Status       int64  `json:"status" `
	DeflationPre string `json:"deflation_pre" `
	Deflation    string `json:"deflation" `
}

func (m DeflationRateVoteInfor) Type() string { return RelationKeyDeflationVoteInforResult }

func (m DeflationRateVoteInfor) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type TradeInfor struct {
	Txhash         string      `json:"tx_hash" `
	From           string      `json:"from"`
	To             string      `json:"to"`
	Amount         RealCoin    `json:"amount"`
	Fee            RealCoin    `json:"fee"`
	TradeType      TranserType `json:"trade_type" `
	BlockNum       int64       `json:"block_num" `
	TimeStemp      int64       `json:"time_stemp" `
	FromFsvBalance string      `json:"from_fsv_balance"`
	ToFsvBalance   string      `json:"to_fsv_balance"`
	FromTipBalance string      `json:"from_tip_balance"`
	ToTipBalance   string      `json:"to_tip_balance"`
}


func (m TradeInfor) Type() string { return RelationKeyTradeInforResult }


func (m TradeInfor) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type BuySourceInfor struct {
	Txhash    string   `json:"tx_hash" `
	DataHash  string   `json:"data_hash"`
	Creator   string   `json:"creator"`
	Purchaser string   `json:"purchaser"`
	Price     RealCoin `json:"price"`
	Fee       RealCoin `json:"fee"`
	Remark    string   `json:"remark" `
	TimeStemp int64    `json:"time_stemp" `
}

func (m BuySourceInfor) Type() string { return RelationKeyBuySourceInforResult }

func (m BuySourceInfor) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type SourceSizeInfor struct {
	Txhash    string `json:"tx_hash" `
	Account   string `json:"account"`
	DataHash  string `json:"data_hash"`
	Size      int64  `json:"size"`
	BlockNum  int64  `json:"block_num" `
	TimeStemp int64  `json:"time_stemp" `
}

func (m SourceSizeInfor) Type() string { return RelationKeySourceSizeInforResult }

func (m SourceSizeInfor) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type MortgageRedeem struct {
	MortageAccount string   `json:"mortage_account"`
	MortgageAmount RealCoin `json:"mortgage_amount"`
	RedeemAccount  string   `json:"redeem_account"`
	RedeemAmount   RealCoin `json:"redeem_amount"`
	BlockNum       int64    `json:"block_num" `
	TimeStemp      int64    `json:"time_stemp" `
	TxHash         string   `json:"tx_hash" `
}

func (m MortgageRedeem) Type() string { return RelationKeyMortgageRedeemResult }

func (m MortgageRedeem) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type CopyrightMove struct {
	FromAccount    string `json:"from_account"`
	ToAccount      string `json:"to_account"`
	DataHash       string `json:"data_hash"`
	BlockNum       int64  `json:"block_num" `
	TimeStemp      int64  `json:"time_stemp" `
	ComplainId     string `json:"complain_id"`
	ComplainStatus string `json:"complain_status"`
	Status         int    `json:"status"`
}

func (m CopyrightMove) Type() string { return RelationKeyCopyrightMoveResult }

func (m CopyrightMove) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type SpaceMiner struct {
	Txhash string `json:"tx_hash"`
	Space  string `json:"space"`
}

func (m SpaceMiner) Type() string { return RelationKeySpaceMinerResult }

func (m SpaceMiner) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type DeflationVote struct {
	Txhash    string `json:"tx_hash"`
	VoteIndex string `json:"vote_index"`
}

func (m DeflationVote) Type() string { return RelationKeyDeflationVoteResult }

func (m DeflationVote) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type CopyrightApprove struct {
	PrimaryHash  string   `json:"primary_hash"`
	DataHash     string   `json:"data_hash"`
	BlockNum     int64    `json:"block_num" `
	TimeStemp    int64    `json:"time_stemp" `
	Status       int      `json:"status"`
	RemoveName   []string `json:"remove_name"`
	NodeByte     string   `json:"node_byte"`
	NewSize      uint64   `json:"new_size"`
	PicNodeByte  string   `json:"pic_node_byte"`
	Creator      string   `json:"creator"`
	EditorStatus int      `json:"editor_status"`
}

func (m CopyrightApprove) Type() string { return RelationKeyCopyrightApproveResult }

func (m CopyrightApprove) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type CopyrightVoteRedeem struct {
	DataHash  string `json:"data_hash"`
	BlockNum  int64  `json:"block_num" `
	TimeStemp int64  `json:"time_stemp" `
	TxHash    string `json:"tx_hash"`
}

func (m CopyrightVoteRedeem) Type() string { return RelationKeyCopyrightVoteRedeemResult }

func (m CopyrightVoteRedeem) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}

type CopyrightVote struct {
	DataHash   string `json:"data_hash"`
	VideoHash  string `json:"video_hash"`
	SourceName string `json:"source_name"`
	VideoName  string `json:"video_name"`
	Account    string `json:"account"`
	Award      string `json:"award"`
	Power      string `json:"power"`
	VoteStatus int    `json:"vote_status"`
	BlockNum   int64  `json:"block_num" `
	TimeStemp  int64  `json:"time_stemp" `
	TxHash     string `json:"tx_hash"`
}

func (m CopyrightVote) Type() string { return RelationKeyCopyrightVoteResult }

func (m CopyrightVote) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}
