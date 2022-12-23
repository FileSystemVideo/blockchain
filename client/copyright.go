package client

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/client/rest"
	"fs.video/blockchain/x/copyright/keeper"
	"fs.video/blockchain/x/copyright/types"
	"github.com/sirupsen/logrus"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type CopyrightClient struct {
	TxClient      *TxClient
	AccountClient *AccountClient
	ServerUrl     string
	logPrefix     string
}

type AccountHashRelaiont struct {
	Status   int    `json:"status"`  
	Message  string `json:"message"` 
	Relation string `json:"data"`    // purchaser  author  none
}

func (this *CopyrightClient) QueryResourceAndHashRelation(account, hash string) (data *AccountHashRelaiont, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": account, "hash": hash})
	params := types.QueryResourceAndHashRelationParams{Account: account, Hash: hash}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryBonusExtrainfor, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	resourceRelation := &keeper.CopyrightExtrainfor{}
	if resBytes != nil {
		err = util.Json.Unmarshal(resBytes, resourceRelation)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return nil, err
		}
	}

	data = &AccountHashRelaiont{}
	

	/*if resourceRelation.Species == "buy" {
		data.Status = 1
		data.Relation = "purchaser"
	} else if resourceRelation.Species == "mortgage" && (!resourceRelation.Downer.Empty() && height < resourceRelation.Height) { 
		data.Status = 1
		data.Relation = "mortgage"
	} */
	
	//if height < resourceRelation.Height {
	if resourceRelation.Species == "buy" {
		data.Status = 1
		data.Relation = "purchaser"
	} else {
		copyrightInfor, err := this.Find(hash)
		if err != nil {
			return data, err
		}
		if copyrightInfor.DataHash == "" {
			data.Message = rest.DataHashNotExist
			return data, err
		}
		if copyrightInfor.Creator.String() == account {
			data.Relation = "author"
		} else {
			data.Relation = "none"
		}
		data.Status = 1
	}
	return
}

func (this *CopyrightClient) QueryPubCount(dayString string) (data *keeper.CopyrightCountInfor, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("day", dayString)
	params := types.QueryPubCountParams{dayString}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryPubCount, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	data = &keeper.CopyrightCountInfor{}
	err = util.Json.Unmarshal(resBytes, data)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (this *CopyrightClient) Find(hash string) (data *types.CopyrightData, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("hash", hash)
	params := types.QueryCopyrightDetailParams{hash}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightDetail, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	data = &types.CopyrightData{}
	err = util.Json.Unmarshal(resBytes, data)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (this *CopyrightClient) Create(copyData types.CopyrightData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgCreateCopyright(copyData)
	if err != nil {
		log.WithError(err).Error("NewMsgCreateCopyright")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(copyData.Creator.String(), privateKey, copyData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *CopyrightClient) Vote(copyData types.CopyrightVoteData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgCopyrightVote(copyData)
	if err != nil {
		log.WithError(err).Error("NewMsgCopyrightVote")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(copyData.Address, privateKey, copyData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *CopyrightClient) Editor(copyData types.EditorCopyrightData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgEditorCopyright(copyData)
	if err != nil {
		log.WithError(err).Error("NewMsgEditorCopyright")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(copyData.Creator.String(), privateKey, copyData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *CopyrightClient) Delete(copyData types.DeleteCopyrightData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgDeleteCopyright(copyData)
	if err != nil {
		log.WithError(err).Error("NewMsgDeleteCopyright")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(copyData.Creator.String(), privateKey, copyData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *CopyrightClient) Bonus(copyData types.CopyrightBonusData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgCopyrightBonusV2(copyData)
	if err != nil {
		log.WithError(err).Error("NewMsgCopyrightBonus")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(copyData.Downer.String(), privateKey, copyData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *CopyrightClient) BonusRear(copyData types.CopyrightBonusRearData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgCopyrightBonusRearV2(copyData)
	if err != nil {
		log.WithError(err).Error("NewMsgCopyrightBonusRear")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(copyData.BonusAddress, privateKey, copyData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *CopyrightClient) QueryCopyrightExport(typesStr, rpcUrl string) (data *types.GenesisState, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	rpcClient, err := rpchttp.New(rpcUrl, "/websocket")
	if err != nil {
		panic("start ctx client error.")
	}

	clientCtx = clientCtx.WithNodeURI(rpcUrl).
		WithClient(rpcClient)
	params := types.QueryAccountParams{typesStr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightExport, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	data = &types.GenesisState{}
	err = util.Json.Unmarshal(resBytes, data)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}
