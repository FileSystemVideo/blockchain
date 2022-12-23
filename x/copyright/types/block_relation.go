package types

import (
	"fs.video/blockchain/core"
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

// BLockRelation Type() 
func (m DeflationAccountVote) Type() string { return RelationKeyDeflationVotePowerResult }

// BLockRelation Marshal() 
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

// BLockRelation Type() 
func (m DeflationRateVoteInfor) Type() string { return RelationKeyDeflationVoteInforResult }

// BLockRelation Marshal() 
func (m DeflationRateVoteInfor) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}


type TradeInfor struct {
	Txhash         string           `json:"tx_hash" `    //txhash
	From           string           `json:"from"`        
	To             string           `json:"to"`          
	Amount         RealCoin         `json:"amount"`      
	Fee            RealCoin         `json:"fee"`         
	TradeType      core.TranserType `json:"trade_type" ` 
	BlockNum       int64            `json:"block_num" `  
	TimeStemp      int64            `json:"time_stemp" ` 
	FromFsvBalance string           `json:"from_fsv_balance"`
	ToFsvBalance   string           `json:"to_fsv_balance"`
	FromTipBalance string           `json:"from_tip_balance"`
	ToTipBalance   string           `json:"to_tip_balance"`
}

// BLockRelation Type() 
func (m TradeInfor) Type() string { return RelationKeyTradeInforResult }

// BLockRelation Marshal() 
func (m TradeInfor) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}


type BuySourceInfor struct {
	Txhash    string   `json:"tx_hash" `    //txhash
	DataHash  string   `json:"data_hash"`   //hash
	Creator   string   `json:"creator"`     
	Purchaser string   `json:"purchaser"`   
	Price     RealCoin `json:"price"`       
	Fee       RealCoin `json:"fee"`         
	Remark    string   `json:"remark" `     
	TimeStemp int64    `json:"time_stemp" ` 
}

// BLockRelation Type() 
func (m BuySourceInfor) Type() string { return RelationKeyBuySourceInforResult }

// BLockRelation Marshal() 
func (m BuySourceInfor) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}


type SourceSizeInfor struct {
	Txhash    string `json:"tx_hash" `    //txhash
	Account   string `json:"account"`     
	DataHash  string `json:"data_hash"`   //hash
	Size      int64  `json:"size"`        
	BlockNum  int64  `json:"block_num" `  
	TimeStemp int64  `json:"time_stemp" ` 
}

// BLockRelation Type() 
func (m SourceSizeInfor) Type() string { return RelationKeySourceSizeInforResult }

// BLockRelation Marshal() 
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
	TxHash         string   `json:"tx_hash" `        //txhash,
}

// BLockRelation Type() 
func (m MortgageRedeem) Type() string { return RelationKeyMortgageRedeemResult }

// BLockRelation Marshal() 
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
	DataHash       string `json:"data_hash"`       //hash
	BlockNum       int64  `json:"block_num" `      
	TimeStemp      int64  `json:"time_stemp" `     
	ComplainId     string `json:"complain_id"`     //id
	ComplainStatus string `json:"complain_status"` 
	Status         int    `json:"status"`          
}

// BLockRelation Type() 
func (m CopyrightMove) Type() string { return RelationKeyCopyrightMoveResult }

// BLockRelation Marshal() 
func (m CopyrightMove) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}


type SpaceMiner struct {
	Txhash string `json:"tx_hash"` //txhash
	Space  string `json:"space"`   
}

// BLockRelation Type() 
func (m SpaceMiner) Type() string { return RelationKeySpaceMinerResult }

// BLockRelation Marshal() 
func (m SpaceMiner) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}


type DeflationVote struct {
	Txhash    string `json:"tx_hash"`    //txhash
	VoteIndex string `json:"vote_index"` 
}

// BLockRelation Type() 
func (m DeflationVote) Type() string { return RelationKeyDeflationVoteResult }

// BLockRelation Marshal() 
func (m DeflationVote) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}


type CopyrightApprove struct {
	PrimaryHash  string   `json:"primary_hash"`  //hash
	DataHash     string   `json:"data_hash"`     //hash
	BlockNum     int64    `json:"block_num" `    
	TimeStemp    int64    `json:"time_stemp" `   
	Status       int      `json:"status"`        
	RemoveName   []string `json:"remove_name"`   //,,files
	NodeByte     string   `json:"node_byte"`     
	NewSize      uint64   `json:"new_size"`      
	PicNodeByte  string   `json:"pic_node_byte"` 
	Creator      string   `json:"creator"`       
	EditorStatus int      `json:"editor_status"` //datahash
}

// BLockRelation Type() 
func (m CopyrightApprove) Type() string { return RelationKeyCopyrightApproveResult }

// BLockRelation Marshal() 
func (m CopyrightApprove) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}


type CopyrightVoteRedeem struct {
	DataHash  string `json:"data_hash"`   //hash
	BlockNum  int64  `json:"block_num" `  
	TimeStemp int64  `json:"time_stemp" ` 
	TxHash    string `json:"tx_hash"`     //hash
}

// BLockRelation Type() 
func (m CopyrightVoteRedeem) Type() string { return RelationKeyCopyrightVoteRedeemResult }

// BLockRelation Marshal() 
func (m CopyrightVoteRedeem) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}


type CopyrightVote struct {
	DataHash   string `json:"data_hash"`   //hash
	VideoHash  string `json:"video_hash"`  //hash
	SourceName string `json:"source_name"` 
	VideoName  string `json:"video_name"`  
	Account    string `json:"account"`     
	Award      string `json:"award"`       
	Power      string `json:"power"`       
	VoteStatus int    `json:"vote_status"` 
	BlockNum   int64  `json:"block_num" `  
	TimeStemp  int64  `json:"time_stemp" ` 
	TxHash     string `json:"tx_hash"`     //hash
}

// BLockRelation Type() 
func (m CopyrightVote) Type() string { return RelationKeyCopyrightVoteResult }

// BLockRelation Marshal() 
func (m CopyrightVote) Marshal() []byte {
	data, err := util.Json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return data
}
