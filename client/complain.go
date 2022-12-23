package client

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/keeper"
	"fs.video/blockchain/x/copyright/types"
)

type ComplainClient struct {
	TxClient  *TxClient
	ServerUrl string
	logPrefix string
}


func (this *ComplainClient) CreateComplain(complainData types.CopyrightComplainData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgCopyrightComplain(complainData)
	if err != nil {
		log.WithError(err).Error("NewMsgCopyrightComplain")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(complainData.ComplainAccount.String(), privateKey, complainData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *ComplainClient) ComplainResponse(complainResponseData types.ComplainResponseData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgComplainResponse(complainResponseData)
	if err != nil {
		log.WithError(err).Error("NewMsgComplainResponse")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(complainResponseData.AccuseAccount.String(), privateKey, complainResponseData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *ComplainClient) ComplainVote(complainVoteData types.ComplainVoteData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgComplainVote(complainVoteData)
	if err != nil {
		log.WithError(err).Error("NewMsgComplainVote")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(complainVoteData.VoteAccount.String(), privateKey, complainVoteData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *ComplainClient) QueryComplainInfor(complainId string) (data *keeper.CopyRightComplain, blockNum int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryCopyrightComplainParams{complainId}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, 0, err
	}
	resBytes, blockNum, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightComplain, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, 0, err
	}
	data = &keeper.CopyRightComplain{}
	err = util.Json.Unmarshal(resBytes, data)
	return
}
