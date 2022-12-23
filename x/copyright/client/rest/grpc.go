package rest

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/keeper"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/http"
)

var grpcLogPrefix = "restGrpc"

//grpc
func grpcQueryBalance(cliCtx *client.Context, address sdk.AccAddress, denom string) (coin sdk.Coin, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithFields(logrus.Fields{"addr": address, "denom": denom})
	params := bankTypes.QueryBalanceRequest{address.String(), denom}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return coin, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/bank/balance", bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return coin, errors.New(QueryChainInforError)
	}
	err = cliCtx.LegacyAmino.UnmarshalJSON(resBytes, &coin)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return coin, errors.New(UnmarshalError + "2")
	}
	return coin, nil
}

//grpc
func grpcQueryCopyrightPartyExist(cliCtx *client.Context, account string) (exist bool, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("acc", account)
	exist = false
	params := types.QueryCopyrightPartyParams{account}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return exist, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightPartyExist, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return exist, err
	}
	if string(resBytes) == "exist" {
		exist = true
	} else if string(resBytes) == "non-existent" {
		exist = false
	}
	return exist, nil
}

//grpcid
func grpcQueryPublisherExist(cliCtx *client.Context, publisherId string) (exist bool, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("id", publisherId)
	exist = false
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightPublishId, []byte(""))
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return exist, err
	}
	var publisherIdMap map[string]string
	if resBytes != nil {
		err = util.Json.Unmarshal(resBytes, &publisherIdMap)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return exist, err
		}
		if _, ok := publisherIdMap[publisherId]; ok {
			return true, nil
		}
	}
	return exist, nil
}

//grpc
func grpcQueryCopyrightExist(cliCtx *client.Context, hash string) (exist bool, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("hash", hash)
	exist = false
	params := types.QueryCopyrightDetailParams{hash}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return exist, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightExist, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return exist, errors.New(QueryChainInforError)
	}
	if string(resBytes) == "exist" {
		exist = true
	} else if string(resBytes) == "non-existent" {
		exist = false
	}
	return exist, nil
}

//grpchash
func grpcQueryCopyrightOrigin(cliCtx *client.Context, hash string) (exist bool, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("hash", hash)
	exist = false
	params := types.QueryCopyrightDetailParams{hash}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return exist, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryOriginHash, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return exist, errors.New(QueryChainInforError)
	}
	if resBytes != nil {
		var originHash keeper.CopyrightOrigininfor
		err = util.Json.Unmarshal(resBytes, &originHash)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return exist, errors.New(QueryChainInforError)
		}
		if originHash.Status == 1 || originHash.Status == 0 { 
			return true, nil
		}
	}
	return exist, nil
}

func grpcQueryCopyrightAndAccount(cliCtx *client.Context, account, hash string) (data keeper.CopyrightExtrainfor, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithFields(logrus.Fields{"acc": account, "hash": hash})
	params := types.QueryResourceAndHashRelationParams{Account: account, Hash: hash}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryBonusExtrainfor, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	resourceRelation := &keeper.CopyrightExtrainfor{}
	if resBytes != nil {
		err = util.Json.Unmarshal(resBytes, resourceRelation)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return
		}
	}
	return
}

//grpc
func grpcQueryCopyright(cliCtx *client.Context, hash string) (copyrightData types.CopyrightData, height int64, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("hash", hash)
	params := types.QueryCopyrightDetailParams{hash}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, height, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightDetail, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return copyrightData, height, errors.New(QueryChainInforError)
	}
	err = util.Json.Unmarshal(resBytes, &copyrightData)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return copyrightData, height, err
	}
	return
}

//grpc
func grpcQueryCopyrightComplain(cliCtx *client.Context, complainId string) (copyrightData keeper.CopyRightComplain, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("id", complainId)
	params := types.QueryCopyrightComplainParams{complainId}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightComplain, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return copyrightData, errors.New(QueryChainInforError)
	}
	err = util.Json.Unmarshal(resBytes, &copyrightData)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return copyrightData, err
	}
	return
}

//grpc
func grpcQueryAccountSpace(cliCtx *client.Context, size int64, account string) (exist bool, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithFields(logrus.Fields{"acc": account, "size": size})
	exist = false
	params := types.QueryAccountSpaceMinerParams{account}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return exist, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryAccountSpaceMiner, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return exist, errors.New(QueryChainInforError)
	}
	var space keeper.AccountSpaceMiner
	err = util.Json.Unmarshal(resBytes, &space)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return exist, err
	}
	spaceLeft := space.SpaceTotal.Sub(space.UsedSpace).Sub(space.LockedSpace)
	if spaceLeft.LessThan(decimal.NewFromFloat(float64(size))) {
		log.Warn("SpaceNotEnough")
		return exist, errors.New(SpaceNotEnough)
	}
	return true, nil
}

//grpc
func grpcQueryDelegationVote(cliCtx *client.Context, delegatorAddr sdk.AccAddress) (exist bool, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("addr", delegatorAddr)
	exist = false
	params := types.QueryDelegatorParams{DelegatorAddr: delegatorAddr}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return exist, err
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryDelegationShares, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return exist, err
	}
	delegation := sdk.MustNewDecFromStr(string(resBytes))
	if delegation.IsPositive() {
		return true, nil
	}
	return
}

//grpc
func grpcQueryAccountVoteEnough(cliCtx *client.Context, delegatorAddr sdk.AccAddress, power string) (exist bool, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithFields(logrus.Fields{"addr": delegatorAddr, "power": power})
	exist = false
	var dataVoteMap map[string]keeper.DataHashVote
	err = util.Json.Unmarshal([]byte(power), &dataVoteMap)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return
	}
	var totalVoteDecimal decimal.Decimal
	for _, voteInfor := range dataVoteMap {
		totalVoteDecimal = totalVoteDecimal.Add(voteInfor.Power)
	}
	params := types.QueryDelegatorParams{DelegatorAddr: delegatorAddr}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return exist, err
	}
	
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryDelegationShares, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData1")
		return exist, err
	}

	delegation := sdk.MustNewDecFromStr(string(resBytes))
	if !delegation.IsPositive() { 
		return
	}
	totalPower := decimal.RequireFromString(string(resBytes))
	
	freezeBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryDelegationFreeze, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData2")
		return exist, err
	}
	freezeDecimal := decimal.RequireFromString(string(freezeBytes))
	if totalPower.Sub(freezeDecimal).GreaterThanOrEqual(totalVoteDecimal) {
		return true, nil
	}
	return false, errors.New(CopyrightVoteNotEnough)
}

//grpc
//func grpcQueryDelegationVoteExist(cliCtx *client.Context, deflationVoteAddr sdk.AccAddress) (exist bool, err error) {
//	log := util.BuildLog(util.GetFuncName(), util.LmChainRest).WithField("addr", deflationVoteAddr)
//	exist = false
//	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryDeflationVoteInfor, []byte(""))
//	if err != nil {
//		log.WithError(err).Error("QueryWithData")
//		return exist, err
//	}
//	var deflationVoteList keeper.DeflationVoteAccountArray
//	err = util.Json.Unmarshal(resBytes, &deflationVoteList)
//	if err != nil {
//		log.WithError(err).Error("Unmarshal")
//		return exist, err
//	}
//	if len(deflationVoteList.DeflationVoteAccountArray) > 0 {
//		for i := 0; i < len(deflationVoteList.DeflationVoteAccountArray); i++ {
//			deflationVoteAccount := deflationVoteList.DeflationVoteAccountArray[i]
//			if deflationVoteAccount.VoteAccount.Equals(deflationVoteAddr) {
//				return true, nil
//			}
//		}
//	}
//	return
//}

//grpc
func grpcQueryMortgageAmount(cliCtx *client.Context) (exist bool, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	exist = false
	resBytes, height, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryMortgAmount, []byte(""))
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return exist, err
	}
	if height < core.MortgageStartHeight {
		return exist, errors.New(MortgageStartHeight)
	}
	
	minerAmount := decimal.RequireFromString(string(resBytes))
	if core.MinerUpperLimitStandand.LessThan(minerAmount) {
		return exist, types.MortgMinerHasFinish
	}
	return true, nil
}

//grpc
func grpcQueryPubCount(cliCtx *client.Context, dayString string) (pubCount int64, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("day", dayString)
	params := types.QueryPubCountParams{dayString}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return pubCount, err
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryPubCount, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return pubCount, err
	}
	data := &keeper.CopyrightCountInfor{}
	err = util.Json.Unmarshal(resBytes, data)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return pubCount, err
	}
	pubCount = data.Count
	return pubCount, nil
}

//grpc
func grpcQueryMinerBonus(cliCtx *client.Context, account string) (bonusAmount string, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("acc", account)
	params := types.QueryAccountParams{account}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return bonusAmount, err
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QuerySpaceMinerRewardInfo, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return bonusAmount, err
	}
	data := &types.RealCoin{}
	err = util.Json.Unmarshal(resBytes, data)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return bonusAmount, err
	}
	bonusAmount = data.Amount
	return bonusAmount, nil
}


func QueryTransferRateHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := types.DataResponse{}
		res.Data = "0.002"
		SendReponse(w, clientCtx, res)
	}
}
