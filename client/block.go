package client

import (
	"context"
	"encoding/hex"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	copyTypes "fs.video/blockchain/x/copyright/types"
	"fs.video/trerr"
	"github.com/sirupsen/logrus"
	abci "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"strings"
)

type Block struct {
	ChainId            string   //id
	Height             int64    
	Time               string   
	LastCommitHash     string   //hash
	Datahash           string   //hash
	ValidatorsHash     string   //hash
	NextValidatorsHash string   //hash
	ConsensusHash      string   //Hash
	Apphash            string   //hash
	LastResultsHash    string   //Hash
	EvidenceHash       string   //hash
	ProposerAddress    string   
	Txs                []string //hash
	Signatures         []Signature
	LastBlockId        string //ID
	BlockId            string //ID
}

type Signature struct {
	ValidatorAddress string 
	TimeStamp        string 
	Sign             string 
}

type BlockClient struct {
	ServerUrl string
	logPrefix string
}

//func (bc *BlockClient) Find(height string) (block *Block, err error) {
//	var reponse string
//	block = new(Block)
//	if height != "" {
//		reponse, err = GetRequest(bc.ServerUrl, "/blocks/"+height)
//	} else {
//		reponse, err = GetRequest(bc.ServerUrl, "/blocks/latest")
//	}
//	if err != nil {
//		return
//	}
//	jsonData, err := simplejson.NewJson([]byte(reponse))
//	if err != nil {
//		fmt.Println("error:", err)
//		return
//	}
//	
//	block = &Block{}
//	header := jsonData.Get("block").Get("header")
//	block.Height = header.Get("height").MustInt64()
//	block.Datahash = header.Get("data_hash").MustString()
//	block.Time = header.Get("time").MustString()
//	block.Apphash = header.Get("app_hash").MustString()
//	block.ChainId = header.Get("chain_id").MustString()
//	block.ConsensusHash = header.Get("consensus_hash").MustString()
//	block.EvidenceHash = header.Get("evidence_hash").MustString()
//	block.LastCommitHash = header.Get("last_commit_hash").MustString()
//	block.LastResultsHash = header.Get("last_results_hash").MustString()
//	block.ValidatorsHash = header.Get("validators_hash").MustString()
//	block.NextValidatorsHash = header.Get("next_validators_hash").MustString()
//	block.ProposerAddress = header.Get("proposer_address").MustString()

//	
//	signaturesJson := jsonData.Get("block").Get("last_commit").Get("signatures")
//	block.Signatures = []Signature{}
//	signatures := signaturesJson.MustArray()
//	for i := 0; i < len(signatures); i++ {
//		ob := signatures[i].(map[string]interface{})
//		sign := Signature{
//			ValidatorAddress: ob["validator_address"].(string),
//			TimeStamp:        ob["timestamp"].(string),
//		}
//		if ob["signature"] != nil {
//			sign.Sign = ob["signature"].(string)
//		}
//		block.Signatures = append(block.Signatures, sign)
//	}
//	block.Txs = []string{}
//	// txs
//	txsJson := jsonData.Get("block").Get("data").Get("txs")
//	txlist := txsJson.MustArray()
//	for i := 0; i < len(txlist); i++ {
//		base64decode, err := base64.StdEncoding.DecodeString(txlist[i].(string))
//		if err == nil {
//			tx := tenderTypes.Tx(base64decode)
//			block.Txs = append(block.Txs, strings.ToUpper(hex.EncodeToString(tx.Hash())))
//		}
//	}
//	return
//}

/**

*/
func (this *BlockClient) GetSyncInfo() (blockData *coretypes.SyncInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return nil, err
	}
	nodeStatus, err := node.Status(context.Background())
	if err != nil {
		log.WithError(err).Error("node.Status")
		return nil, err
	}
	return &nodeStatus.SyncInfo, nil
}


func (this *BlockClient) SearchBlock(ctx context.Context, minHeight, maxHeight int64) (result *coretypes.ResultBlockchainInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"min": minHeight, "max": maxHeight})
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return nil, err
	}
	if minHeight <= core.InitHeight || maxHeight <= core.InitHeight {
		return nil, trerr.HistoryBlockNotQuery.GetError()
	}
	return node.BlockchainInfo(ctx, minHeight, maxHeight)
}

/**

*/
func (this *BlockClient) Block(height int64) (blockData *coretypes.ResultBlock, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("height", height)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return nil, err
	}
	if height <= core.InitHeight {
		return nil, trerr.HistoryBlockNotQuery.GetError()
	}
	
	return node.Block(context.Background(), &height)
}

/**
ctxClient
*/
func (this *BlockClient) Find(height int64) (blockData *Block, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("height", height)
	blockData = &Block{}
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return nil, err
	}
	
	if height == 0 {
		nodeStatus, err := node.Status(context.Background())
		if err != nil {
			log.WithError(err).Error("node.Status")
			return nil, err
		}
		height = nodeStatus.SyncInfo.LatestBlockHeight
	}

	
	blockInfo, err := node.Block(context.Background(), &height)
	if err != nil {
		log.WithError(err).Error("node.Block")
		return nil, err
	}
	blockData.Height = blockInfo.Block.Height
	blockData.Datahash = blockInfo.Block.DataHash.String()
	blockData.ChainId = blockInfo.Block.ChainID
	blockData.Time = blockInfo.Block.Time.String()
	blockData.Apphash = blockInfo.Block.AppHash.String()
	blockData.ConsensusHash = blockInfo.Block.ConsensusHash.String()
	blockData.EvidenceHash = blockInfo.Block.EvidenceHash.String()
	blockData.LastCommitHash = blockInfo.Block.LastCommitHash.String()
	blockData.LastResultsHash = blockInfo.Block.LastResultsHash.String()
	blockData.ValidatorsHash = blockInfo.Block.ValidatorsHash.String()
	blockData.NextValidatorsHash = blockInfo.Block.NextValidatorsHash.String()
	blockData.ProposerAddress = blockInfo.Block.ProposerAddress.String()
	blockData.LastResultsHash = blockInfo.Block.LastResultsHash.String()
	blockData.LastBlockId = blockInfo.Block.LastBlockID.Hash.String()
	blockData.BlockId = blockInfo.BlockID.Hash.String()
	for _, s := range blockInfo.Block.LastCommit.Signatures {
		signature := new(Signature)
		signature.ValidatorAddress = s.ValidatorAddress.String()
		signature.Sign = string(s.Signature)
		signature.TimeStamp = s.Timestamp.String()
		blockData.Signatures = append(blockData.Signatures, *signature)
	}
	for i := 0; i < len(blockInfo.Block.Txs); i++ {
		resTx, err := node.Tx(context.Background(), blockInfo.Block.Txs[i].Hash(), true)
		if err != nil {
			log.WithError(err).Error("node.Tx")
			return nil, err
		}
		blockData.Txs = append(blockData.Txs, strings.ToUpper(hex.EncodeToString(resTx.Hash)))
	}

	return
}


func (this *BlockClient) FindTxProps(height string) (resp []copyTypes.TxPropValue, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("height", height)
	var reponse string
	reponse, err = GetRequest(this.ServerUrl, "/copyright/list/"+height+"/debug")
	if err != nil {
		log.WithError(err).Error("GetRequest")
		return
	}
	message := new(copyTypes.TxPropsResponse)
	err = clientCtx.LegacyAmino.UnmarshalJSON([]byte(reponse), &message)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, err
	}
	resp = message.Props
	return
}


func (this *BlockClient) FindMsgs(height string) (resp *copyTypes.MessageResp, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("height", height)
	var reponse string
	reponse, err = GetRequest(this.ServerUrl, "/copyright/list/"+height)
	if err != nil {
		log.WithError(err).Error("GetRequest")
		return
	}
	message := new(copyTypes.MsgsResponse)
	err = clientCtx.LegacyAmino.UnmarshalJSON([]byte(reponse), &message)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, err
	}
	resp = &message.Message
	return
}


func (this *BlockClient) FindBlockResults(height *int64) (events []abci.Event, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return nil, err
	}
	blockResults, err := node.BlockResults(context.Background(), height)
	if err != nil {
		log.WithError(err).Error("node.BlockResults")
		return nil, err
	}
	for _, endEvent := range blockResults.EndBlockEvents {
		blockResults.BeginBlockEvents = append(blockResults.BeginBlockEvents, endEvent)
	}
	return blockResults.BeginBlockEvents, nil
}


func (this *BlockClient) FindRelationData(height int64) (rds *copyTypes.BlockRDS, notFound bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("height", height)
	notFound = false
	params := copyTypes.QueryBlockRDSParams{height}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, notFound, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+copyTypes.QueryBlockRDS, bz)
	if err != nil {
		if strings.Contains(err.Error(), copyTypes.RDSNotFoundErr.Error()) {
			notFound = true
			err = nil 
		} else {
			log.WithError(err).Error("QueryWithData")
		}
		return nil, notFound, err
	}
	var data copyTypes.BlockRDS
	err = util.Json.Unmarshal(resBytes, &data)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return &data, notFound, err
}
