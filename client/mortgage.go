package client

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/keeper"
	"fs.video/blockchain/x/copyright/types"
)

type MortgageClient struct {
	TxClient  *TxClient
	ServerUrl string
}


func (this *MortgageClient) QueryMortgageRate() (rate string, blockNum int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, blockNum, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryMiningStage, []byte(""))
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	//log.Debug("find result:", "---", string(resBytes))
	rate = string(resBytes)
	return
}


func (this *MortgageClient) QueryMortgageMinerInfor() (rate string, blockNum int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, blockNum, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryMortgMiningInfor, []byte(""))
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	//log.Debug("find result:", "---", string(resBytes))
	rate = string(resBytes)
	return
}


func (this *MortgageClient) QueryMortgageAmount() (mortAmount string, blockNum int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, blockNum, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryMortgAmount, []byte(""))
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	//log.Debug("find result:", "---", string(resBytes))
	mortAmount = string(resBytes)
	return
}


func (this *MortgageClient) Mortgage(data types.MortgageData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	
	msg, err := types.NewMsgMsgMortgage(data)
	if err != nil {
		return
	}
	
	resp, err = this.TxClient.SignAndSendMsg(data.MortgageAccount.String(), privateKey, data.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *MortgageClient) QueryDeflationRate() (*keeper.RateAndVoteIndex, int64, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	bz := []byte("")
	resBytes, height, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryDeflationRateInfor, bz)
	if err != nil {
		log.Error("QueryWithData")
	}
	//log.Debug("find result:", "---", string(resBytes))
	rateVoteIndex := &keeper.RateAndVoteIndex{}
	err = util.Json.Unmarshal(resBytes, rateVoteIndex)
	return rateVoteIndex, height, err
}
