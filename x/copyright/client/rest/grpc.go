package rest

import (
	"errors"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/keeper"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/shopspring/decimal"
)


func grpcQueryBalance(cliCtx *client.Context, address sdk.AccAddress, denom string) (coin sdk.Coin, err error) {
	params := bankTypes.QueryBalanceRequest{address.String(), denom}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return coin, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/bank/balance", bz)
	if err != nil {
		logs.Error("QueryWithData error:", err.Error())
		return coin, errors.New(QueryChainInforError)
	}

	err = cliCtx.LegacyAmino.UnmarshalJSON(resBytes, &coin)
	if err != nil {
		return coin, errors.New(UnmarshalError + "2")
	}
	return coin, nil
}


func grpcQueryCopyrightPartyExist(cliCtx *client.Context, account string) (exist bool, err error) {
	exist = false
	params := types.QueryCopyrightPartyParams{account}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return exist, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightPartyExist, bz)
	if err != nil {
		return exist, err
	}
	if string(resBytes) == "exist" {
		exist = true
	} else if string(resBytes) == "non-existent" {
		exist = false
	}
	return exist, nil
}


func grpcQueryPublisherExist(cliCtx *client.Context, publisherId string) (exist bool, err error) {
	exist = false
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightPublishId, []byte(""))
	if err != nil {
		logs.Error("QueryWithData  error :", err.Error())
		return exist, err
	}
	var publisherIdMap map[string]string
	if resBytes != nil {
		err = util.Json.Unmarshal(resBytes, &publisherIdMap)
		if err != nil {
			logs.Error("Unmarshal error:", err.Error())
			return exist, err
		}
		if _, ok := publisherIdMap[publisherId]; ok {
			return true, nil
		}
	}
	return exist, nil
}


func grpcQueryCopyrightExist(cliCtx *client.Context, hash string) (exist bool, err error) {
	exist = false
	params := types.QueryCopyrightDetailParams{hash}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return exist, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightExist, bz)
	if err != nil {
		logs.Error("QueryWithData error:", err.Error())
		return exist, errors.New(QueryChainInforError)
	}
	if string(resBytes) == "exist" {
		exist = true
	} else if string(resBytes) == "non-existent" {
		exist = false
	}
	return exist, nil
}


func grpcQueryCopyrightOrigin(cliCtx *client.Context, hash string) (exist bool, err error) {
	exist = false
	params := types.QueryCopyrightDetailParams{hash}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return exist, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryOriginHash, bz)
	if err != nil {
		logs.Error("QueryWithData error:", err.Error())
		return exist, errors.New(QueryChainInforError)
	}

	if resBytes != nil {
		var originHash keeper.CopyrightOrigininfor
		err = util.Json.Unmarshal(resBytes, &originHash)
		if err != nil {
			return exist, errors.New(QueryChainInforError)
		}
		if originHash.Status == 1 || originHash.Status == 0 {
			return true, nil
		}
	}
	return exist, nil
}

func grpcQueryCopyrightAndAccount(cliCtx *client.Context, account, hash string) (data keeper.CopyrightExtrainfor, err error) {
	params := types.QueryResourceAndHashRelationParams{Account: account, Hash: hash}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryBonusExtrainfor, bz)
	if err != nil {
		return
	}
	resourceRelation := &keeper.CopyrightExtrainfor{}
	if resBytes != nil {
		err = util.Json.Unmarshal(resBytes, resourceRelation)
		if err != nil {
			return
		}
	}
	return
}


func grpcQueryCopyright(cliCtx *client.Context, hash string) (copyrightData types.CopyrightData, height int64, err error) {
	params := types.QueryCopyrightDetailParams{hash}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return
	}
	resBytes, height, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightDetail, bz)
	if err != nil {
		logs.Error("QueryWithData error:", err.Error())
		return copyrightData, height, errors.New(QueryChainInforError)
	}
	err = util.Json.Unmarshal(resBytes, &copyrightData)
	if err != nil {
		logs.Error("Unmarshal error", err)
		return copyrightData, height, err
	}
	return
}


func grpcQueryAuthorizeAccount(cliCtx *client.Context, hash string) (authorizeData types.AuthorizeAccountData, err error) {
	params := types.QueryAuthorizeAccountParams{hash}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryAuthorizePubkey, bz)
	if err != nil {
		logs.Error("QueryWithData error:", err.Error())
		return authorizeData, errors.New(QueryChainInforError)
	}
	err = util.Json.Unmarshal(resBytes, &authorizeData)
	if err != nil {
		logs.Error("Unmarshal error", err)
		return authorizeData, err
	}
	return
}


func grpcQueryCopyrightComplain(cliCtx *client.Context, complainId string) (copyrightData keeper.CopyRightComplain, err error) {
	params := types.QueryCopyrightComplainParams{complainId}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryCopyrightComplain, bz)
	if err != nil {
		logs.Error("QueryWithData error:", err.Error())
		return copyrightData, errors.New(QueryChainInforError)
	}
	err = util.Json.Unmarshal(resBytes, &copyrightData)
	if err != nil {
		logs.Error("Unmarshal error", err)
		return copyrightData, err
	}
	return
}


func grpcQueryAccountSpace(cliCtx *client.Context, size int64, account string) (exist bool, err error) {
	exist = false
	params := types.QueryAccountSpaceMinerParams{account}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return exist, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryAccountSpaceMiner, bz)
	if err != nil {
		return exist, errors.New(QueryChainInforError)
	}
	var space keeper.AccountSpaceMiner
	err = util.Json.Unmarshal(resBytes, &space)
	if err != nil {
		return exist, err
	}
	spaceLeft := space.SpaceTotal.Sub(space.UsedSpace).Sub(space.LockedSpace)
	if spaceLeft.LessThan(decimal.NewFromFloat(float64(size))) {
		return exist, errors.New(SpaceNotEnough)
	}
	return true, nil
}


func grpcQueryDelegationVote(cliCtx *client.Context, delegatorAddr sdk.AccAddress) (exist bool, err error) {
	exist = false
	params := types.QueryDelegatorParams{DelegatorAddr: delegatorAddr}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return exist, err
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryDelegationShares, bz)
	if err != nil {
		return exist, err
	}

	delegation := sdk.MustNewDecFromStr(string(resBytes))
	if delegation.IsPositive() {
		return true, nil
	}
	return
}


func grpcQueryAccountVoteEnough(cliCtx *client.Context, delegatorAddr sdk.AccAddress, power string) (exist bool, err error) {
	exist = false
	var dataVoteMap map[string]keeper.DataHashVote
	err = util.Json.Unmarshal([]byte(power), &dataVoteMap)
	if err != nil {
		return
	}
	var totalVoteDecimal decimal.Decimal
	for _, voteInfor := range dataVoteMap {
		totalVoteDecimal = totalVoteDecimal.Add(voteInfor.Power)
	}
	params := types.QueryDelegatorParams{DelegatorAddr: delegatorAddr}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return exist, err
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryDelegationShares, bz)
	if err != nil {
		return exist, err
	}

	delegation := sdk.MustNewDecFromStr(string(resBytes))
	if !delegation.IsPositive() {
		return
	}
	totalPower := decimal.RequireFromString(string(resBytes))

	freezeBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryDelegationFreeze, bz)
	if err != nil {
		return exist, err
	}
	freezeDecimal := decimal.RequireFromString(string(freezeBytes))
	if totalPower.Sub(freezeDecimal).GreaterThanOrEqual(totalVoteDecimal) {
		return true, nil
	}
	return false, errors.New(CopyrightVoteNotEnough)
}


func grpcQueryDelegationVoteExist(cliCtx *client.Context, deflationVoteAddr sdk.AccAddress) (exist bool, err error) {
	exist = false
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryDeflationVoteInfor, []byte(""))
	if err != nil {
		return exist, err
	}
	var deflationVoteList keeper.DeflationVoteAccountArray
	err = util.Json.Unmarshal(resBytes, &deflationVoteList)
	if err != nil {
		return exist, err
	}
	if len(deflationVoteList.DeflationVoteAccountArray) > 0 {
		for i := 0; i < len(deflationVoteList.DeflationVoteAccountArray); i++ {
			deflationVoteAccount := deflationVoteList.DeflationVoteAccountArray[i]
			if deflationVoteAccount.VoteAccount.Equals(deflationVoteAddr) {
				return true, nil
			}
		}
	}
	return
}


func grpcQueryMortgageAmount(cliCtx *client.Context) (exist bool, err error) {
	exist = false
	resBytes, height, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryMortgAmount, []byte(""))
	if err != nil {
		return exist, err
	}

	if height < config.MortgageStartHeight {
		return exist, errors.New(MortgageStartHeight)
	}

	minerAmount := decimal.RequireFromString(string(resBytes))
	if config.MinerUpperLimitStandand.LessThan(minerAmount) {
		return exist, types.MortgMinerHasFinish
	}
	return true, nil
}


func grpcQueryPubCount(cliCtx *client.Context, dayString string) (pubCount int64, err error) {
	params := types.QueryPubCountParams{dayString}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return pubCount, err
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/copyright/"+types.QueryPubCount, bz)
	if err != nil {
		return pubCount, err
	}
	data := &keeper.CopyrightCountInfor{}
	err = util.Json.Unmarshal(resBytes, data)
	if err != nil {
		return pubCount, err
	}
	pubCount = data.Count
	return pubCount, nil
}
