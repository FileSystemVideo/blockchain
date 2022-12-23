package client

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
)

type CopyrightPartyClient struct {
	TxClient  *TxClient
	ServerUrl string
	logPrefix string
}


func (this *CopyrightPartyClient) IsRegisted(bech32AccountAddr string) (exist bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryCopyrightPartyParams{bech32AccountAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return exist, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightPartyExist, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return exist, err
	}
	if string(resBytes) == "exist" {
		exist = true
	} else if string(resBytes) == "non-existent" {
		exist = false
	}
	return
}


func (this *CopyrightPartyClient) Find(bech32AccountAddr string) (data *types.CopyrightPartyData, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryCopyrightPartyParams{bech32AccountAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightParty, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	data = &types.CopyrightPartyData{}
	err = util.Json.Unmarshal(resBytes, data)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}

//id
func (this *CopyrightPartyClient) QueryPublisherMap() (publisherMap map[string]string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightPublishId, []byte(""))
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	if len(resBytes) == 0 {
		return
	}
	err = util.Json.Unmarshal(resBytes, &publisherMap)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (this *CopyrightPartyClient) Register(data types.CopyrightPartyData, privateKey string) (txRes *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgRegisterCopyrightParty(data)
	if err != nil {
		log.WithError(err).Error("NewMsgRegisterCopyrightParty")
		return
	}

	
	resp, err := this.TxClient.SignAndSendMsg(data.Creator.String(), privateKey, data.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}
