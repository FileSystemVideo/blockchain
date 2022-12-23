package client

import (
	"fs.video/blockchain/core"
	"fs.video/blockchain/x/copyright/types"
)

//deprecated
type AuthorizeClient struct {
	TxClient  *TxClient
	ServerUrl string
}


func (this *AuthorizeClient) QueryMortgageRate() (rate string, blockNum int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, blockNum, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryMiningStage, []byte(""))
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	rate = string(resBytes)
	return
}
