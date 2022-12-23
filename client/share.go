package client

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
)


type ShareClient struct {
	TxClient  *TxClient
	ServerUrl string
	logPrefix string
}


func (this *ShareClient) QueryInviteRecord(inviteAddress string) (record *[]types.InviteRecording, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryInviteRecordParams{InviteAddress: inviteAddress}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).WithField("address", inviteAddress).Error("MarshalJSON")
		return
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryInviteRecord, bz)
	if err != nil {
		log.WithError(err).WithField("address", inviteAddress).Error("QueryWithData")
		return
	}
	if resBytes != nil {
		record = &[]types.InviteRecording{}
		err = util.Json.Unmarshal(resBytes, record)
		if err != nil {
			log.WithError(err).WithField("address", inviteAddress).Error("Unmarshal")
		}
	}
	return
}


func (this *ShareClient) QueryInviteStatistics(account string) (statistics *types.InviteRewardStatistics, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryAccountParams{Account: account}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).WithField("acc", account).Error("MarshalJSON")
		return
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryInviteStatistics, bz)
	if err != nil {
		log.WithError(err).WithField("acc", account).Error("QueryWithData")
		return
	}
	if resBytes != nil {
		statistics = &types.InviteRewardStatistics{}
		err = util.Json.Unmarshal(resBytes, statistics)
		if err != nil {
			log.WithError(err).WithField("acc", account).Error("Unmarshal")
		}
	}
	return
}


func (this *ShareClient) InviteReward(data types.InviteRewardData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	msg, err := types.NewMsgInviteReward(data)
	if err != nil {
		log.WithError(err).Error("NewMsgInviteReward")
		return
	}
	
	resp, err = this.TxClient.SignAndSendMsg(data.Address, privateKey, data.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *ShareClient) QueryInviteReward(account string) (reward map[string]string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryAccountParams{Account: account}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).WithField("acc", account).Error("MarshalJSON")
		return
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryRewardInfo, bz)
	if err != nil {
		log.WithError(err).WithField("acc", account).Error("QueryWithData")
		return
	}
	if resBytes != nil {
		err = util.Json.Unmarshal(resBytes, &reward)
		if err != nil {
			log.WithError(err).WithField("acc", account).Error("Unmarshal")
		}
	}
	return
}
