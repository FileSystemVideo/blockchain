package client

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/keeper"
	"fs.video/blockchain/x/copyright/types"
)

type NftClient struct {
	TxClient  *TxClient
	ServerUrl string
}

//nft
func (this *NftClient) QueryNftInfor(tokenId string) (data *keeper.NftInfor, blockNum int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryNftInforParams{tokenId}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return nil, 0, err
	}
	resBytes, blockNum, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryNftInfor, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, 0, err
	}
	data = &keeper.NftInfor{}
	err = util.Json.Unmarshal(resBytes, data)
	return
}


func (this *NftClient) Transfer(data types.NftTransferData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	
	msg, err := types.NewMsgNftTransfer(data)
	if err != nil {
		return
	}
	
	resp, err = this.TxClient.SignAndSendMsg(data.From.String(), privateKey, data.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}
