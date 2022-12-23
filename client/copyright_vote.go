package client

import (
	"errors"
	"fs.video/blockchain/x/copyright/types"
)

type CopyrightVoteClient struct {
	TxClient  *TxClient
	ServerUrl string
}


func (this *CopyrightVoteClient) CompyrightVote(copyrightVoteData types.CopyrightVoteData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	
	msg, err := types.NewMsgCopyrightVote(copyrightVoteData)
	if err != nil {
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(copyrightVoteData.Address, privateKey, copyrightVoteData.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}
