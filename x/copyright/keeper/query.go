package keeper

import (
	// this line is used by starport scaffolding # 1

	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
	"strings"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)
		switch path[0] {
		case types.QueryCopyrightDetail: 
			return queryCopyrightDetail(ctx, req, k, legacyQuerierCdc)
		case types.QueryOriginHash: //hash
			return queryOriginDatahashDetail(ctx, req, k, legacyQuerierCdc)
		case types.QueryPubCount: 
			return queryPubCount(ctx, req, k, legacyQuerierCdc)
		case types.QueryBonusExtrainfor: 
			return queryCopyrightBonusInfor(ctx, req, k, legacyQuerierCdc)
		case types.QueryCopyrightExist: 
			return queryCopyrightExist(ctx, req, k, legacyQuerierCdc)
		case types.QueryCopyrightParty: 
			return queryCopyrightParty(ctx, req, k, legacyQuerierCdc)
		case types.QueryCopyrightPublishId: //id
			return queryCopyrightPublishId(ctx, k)
		case types.QueryCopyrightPartyExist: 
			return queryCopyrightPartyExist(ctx, req, k, legacyQuerierCdc)
		case types.QueryAccountSpaceMiner: 
			return queryAccountSpaceMiner(ctx, req, k, legacyQuerierCdc)
		case types.QuerySpaceFee: 
			return querySpacefee(ctx, k)
		case types.QueryTotalSpaceInfor: 
			return spaceTotal(ctx, k)
		case types.QuerySpaceAward: 
			return spaceAward(ctx, k)
		case types.QueryHasMinerBonus: 
			return hasMinerBonus(ctx, k)
		case types.QuerySpaceAmount: 
			return spaceAmount(ctx, k)
		case types.QueryDeflationMinerInfor: 
			return queryDeflationMinerInfor(ctx, k)
		case types.QueryDeflationRateInfor: 
			return queryDeflationRateInfor(ctx, k)
		case types.QueryNftInfor: //nft
			return queryNftInfor(ctx, req, k, legacyQuerierCdc)
		case types.QueryBlockRDS: //RDS
			return queryBlockRDS(ctx, req, k, legacyQuerierCdc)
		case types.QueryValidatorByConsAddress: //POS，
			return queryValidatorByConsAddress(ctx, req, k, legacyQuerierCdc)
		case types.QueryValidatorDelegationDetail: //，、、
			return queryValidatorDelegationDetail(ctx, req, k, legacyQuerierCdc)
		case types.QuerySigningInfo: 
			return queryValidatorSigningInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryDelegationShares: 
			return queryDelegationShares(ctx, req, k, legacyQuerierCdc)
		case types.QueryDelegation: 
			return queryDelegation(ctx, req, k, legacyQuerierCdc)
		case types.QueryDelegationPreview: //POS，
			return queryDelegationPreview(ctx, req, k, legacyQuerierCdc)
		case types.QueryUnbondingDelegationPreview: //POS，
			return queryUnbondingDelegationPreview(ctx, req, k, legacyQuerierCdc)
		case types.QueryTotalShares: 
			return queryTotalShares(ctx, k, legacyQuerierCdc)
		case types.QueryDelegationByConsAddress: 
			return queryDelegationByConsAddress(ctx, req, k, legacyQuerierCdc)
		case types.QueryInviteRecord: 
			return queryInviteRecord(ctx, req, k, legacyQuerierCdc)
		case types.QueryInviteStatistics: 
			return queryInviteRewardStatistics(ctx, req, k, legacyQuerierCdc)
		case types.QueryMortgAmount: 
			return queryMortgAmount(ctx, k)
		case types.QueryMiningStage: 
			return []byte(strconv.FormatInt(core.MortgageRate, 10)), nil
		case types.QueryMortgMiningInfor: 
			return queryMortgMinerInfor(ctx, k)
		case types.QueryCopyrightComplain: 
			return queryCopyrightComplain(ctx, req, k, legacyQuerierCdc)
		case types.QueryRewardInfo: 
			return queryRewardInfo(ctx, req, k, legacyQuerierCdc)
		case types.QuerySpaceMinerRewardInfo: 
			return querySpaceMinerRewardInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryValidatorInfo: 
			return queryValidatorInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryDelegationFreeze: 
			return queryDelegationFreeze(ctx, req, k, legacyQuerierCdc)
		case types.QueryCopyrightExport: 
			return queryCopyrightExport(ctx, req, k, legacyQuerierCdc)
		case types.QueryParams: 
			return queryParams(ctx, k, legacyQuerierCdc)
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}
func queryParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}


func queryCopyrightExport(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryAccountParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	genesisState := types.DefaultGenesis()
	switch params.Account {
	case "accountSpace":
		res := k.ExportAccountSpace(ctx)
		genesisState.AccountSpace = res
	case "deflationInfor":
		res := k.ExportDeflationInfor(ctx)
		genesisState.DeflationInfor = res
	case "inviteRelation":
		res := k.ExportInviteRelation(ctx)
		genesisState.InviteRelation = res
	case "inviteRecording":
		res := k.ExportInviteRecording(ctx)
		genesisState.InviteRecords = res
	case "inviteReward":
		res := k.ExportInviteReward(ctx)
		genesisState.InviteReward = res
	case "inviteStatistics":
		res := k.ExportInviteStatistics(ctx)
		genesisState.InvitesStatistics = res
	case "copyrightParty":
		res := k.ExportCopyrightParty(ctx)
		genesisState.CopyrightPart = res
	case "copyrightPublishId":
		res := k.ExportCopyrightPublishId(ctx)
		genesisState.CpyrightPublishId = res
	case "copyright":
		res := k.ExportCopyright(ctx)
		genesisState.Copyright = res
	case "copyrightExtra":
		res := k.ExportCopyrightExtraNew(ctx)
		genesisState.CopyrightExtra = res
	case "copyrightIp":
		res := k.ExportCopyrightIp(ctx)
		genesisState.CopyrightIp = res
	case "copyrightOriginHash":
		res := k.ExportCopyrightOriginHash(ctx)
		genesisState.CopyrightOriginHash = res
	case "copyrightBonusAddress":
		res := k.ExportCopyrightBonusAddress(ctx)
		genesisState.CopyrightBonus = res
	case "copyrightNft":
		res := k.ExportCopyrightNft(ctx)
		genesisState.NftInfo = res
	case "copyrightVote":
		res := k.ExportCopyrightVote(ctx)
		genesisState.CopyrightVote = res
	case "copyrightVoteList":
		res := k.ExportCopyrightVoteList(ctx)
		genesisState.CopyrightVoteList = res
	case "copyrightApproveResult":
		res := k.ExportCopyrightApproveResult(ctx)
		genesisState.ApproveResult = res
	case "copyrightVoteRedeem":
		res := k.ExportCopyrightVoteRedeem(ctx)
		genesisState.CopyrightVoteRedeem = res
	}
	res, err := util.Json.Marshal(genesisState)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return res, nil
}


func queryValidatorInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params stakingTypes.QueryValidatorsParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	
	validators := k.stakingKeeper.GetAllValidators(ctx)
	validatorsFilter := stakingTypes.Validators{}
	filteredVals := make([]types.ValidatorInfo, 0)
	if params.Status != "" {
		for _, val := range validators {
			if strings.EqualFold(val.GetStatus().String(), params.Status) {
				validatorsFilter = append(validatorsFilter, val)
			}
		}
	} else {
		validatorsFilter = validators
	}
	for _, val := range validatorsFilter {
		
		validatorInfo := valInfo(val)
		consAddr, err := val.GetConsAddr()
		validatorInfo.ConsAddress = consAddr.String()
		if err != nil {
			continue
		}
		
		validatorSignInfo, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
		if !found {
			continue
		}
		validatorInfo.IndexOffset = validatorSignInfo.IndexOffset
		validatorInfo.MissedBlocksCounter = validatorSignInfo.MissedBlocksCounter
		filteredVals = append(filteredVals, validatorInfo)
	}
	res, err := util.Json.Marshal(filteredVals)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return res, nil
}


func valInfo(val stakingTypes.Validator) types.ValidatorInfo {
	validatorInfo := types.ValidatorInfo{
		OperatorAddress:   val.OperatorAddress,
		Jailed:            val.Jailed,
		Status:            int(val.Status),
		Tokens:            val.Tokens.String(),
		DelegatorShares:   types.RemoveDecLastZero(val.DelegatorShares),
		Moniker:           val.Description.Moniker,
		Identity:          val.Description.Identity,
		Website:           val.Description.Website,
		SecurityContact:   val.Description.SecurityContact,
		Details:           val.Description.Details,
		UnbondingHeight:   val.UnbondingHeight,
		UnbondingTime:     val.UnbondingTime.Unix(),
		Rate:              types.RemoveDecLastZero(val.Commission.Rate),
		MaxRate:           types.RemoveDecLastZero(val.Commission.MaxRate),
		MaxChangeRate:     types.RemoveDecLastZero(val.Commission.MaxChangeRate),
		MinSelfDelegation: val.MinSelfDelegation.String(),
	}
	return validatorInfo
}


func querySpaceMinerRewardInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryAccountParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	//accountBonusMap := k.QueryDeflationAccountBonus(ctx)
	accountBonusString, _, err := k.QueryAccountMinerBonusAmount(ctx, params.Account)
	if err != nil && err != types.NoSpaceMinerRewardErr {
		return nil, err
	}
	realCoinBonus := types.NewRealCoinFromStr(sdk.DefaultBondDenom, accountBonusString)
	res, err := util.Json.Marshal(realCoinBonus)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return res, nil
}


func queryRewardInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryAccountParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	settlement, err := k.QueryRewardInfo(ctx, params.Account)
	if err != nil {
		return nil, err
	}
	reward := make(map[string]string)
	reward["actually"] = "0"
	reward["deserve"] = "0"
	if settlement != nil {
		log.WithFields(logrus.Fields{
			"Capacity expansion": settlement.ExpansionRewardSpace,
			"Invitation reward":  settlement.InviteRewardSpace,
		}).Debug("Account invitation reward information")
		
		deserveSpace := settlement.InviteRewardSpace.Add(settlement.ExpansionRewardSpace)
		reward["deserve"] = deserveSpace.StringFixed(4)

		accountMiner := k.QueryAccountSpaceMinerInfor(ctx, params.Account)
		
		sp := accountMiner.BuySpace.Mul(decimal.RequireFromString(InviteSpaceRateKey)).Sub(accountMiner.RewardSpace)
		
		if sp.LessThan(deserveSpace) {
			deserveSpace = sp
		}
		if sp.IsNegative() {
			deserveSpace = decimal.Zero
		}
		reward["actually"] = deserveSpace.StringFixed(4)
	}
	res, err := util.Json.Marshal(reward)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}


func queryPubCount(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryPubCountParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	copyrightSpaceInfor := k.GetPubCount(ctx, params.DayString)
	res, err := util.Json.Marshal(copyrightSpaceInfor)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}


func queryCopyrightBonusInfor(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryResourceAndHashRelationParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, err
	}
	res := keeper.GetCopyrightBonusInfo(ctx, params.Hash, params.Account)
	return res, nil
}


func queryMortgAmount(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	mortgageAmount := keeper.GetMiningAmount(ctx)
	return []byte(mortgageAmount.Amount), nil
}


func queryMortgMinerInfor(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	mortgageAmount := keeper.QueryMortgMinerInfor(ctx)
	return util.Json.Marshal(mortgageAmount)
}


func queryCopyrightComplain(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryCopyrightComplainParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	complainInfor, err := keeper.GetCopyrightComplainInfor(ctx, params.ComplainId)
	if err != nil {
		return nil, err
	}
	complainBytes, err := util.Json.Marshal(complainInfor)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return complainBytes, nil
}


func queryInviteRecord(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryInviteRecordParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	inviteAddress := params.InviteAddress
	addr, err := sdk.AccAddressFromBech32(inviteAddress)
	if err != nil {
		log.WithError(err).WithField("inviteAddress", inviteAddress).Error("AccAddressFromBech32")
		return nil, err
	}
	record, err := keeper.GetInviteRecording(ctx, addr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if record == nil {
		return nil, nil
	}
	recordByte, err := util.Json.Marshal(record)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return recordByte, nil
}


func queryInviteRewardStatistics(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryAccountParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	address, err := sdk.AccAddressFromBech32(params.Account)
	if err != nil {
		log.WithError(err).WithField("Account", params.Account).Error("AccAddressFromBech32")
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	res, err := keeper.GetInviteRewardStatistics(ctx, address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	resByte, err := util.Json.Marshal(res)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return resByte, nil
}


func queryValidatorByConsAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryValidatorByConsAddrParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	
	validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, params.ValidatorConsAddress)
	if !found {
		return nil, stakingTypes.ErrNoValidatorFound
	}
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, validator)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}


func queryDelegationByConsAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryDelegatorByConsAddrParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	
	validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, params.ValidatorConsAddress)
	if !found {
		return nil, stakingTypes.ErrNoValidatorFound
	}
	
	delegation, found := k.stakingKeeper.GetDelegation(ctx, params.DelegatorAddr, validator.GetOperator())
	if !found {
		return nil, stakingTypes.ErrNoDelegation
	}
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, delegation)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}


func queryTotalShares(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	/*validators := k.stakingKeeper.GetAllValidators(ctx)
	total := sdk.NewDec(0)
	for _, val := range validators {
		total = total.Add(val.DelegatorShares)
	}*/
	totalString := k.GetAllDelegatorShares(ctx)
	return []byte(totalString), nil
}


func queryUnbondingDelegationPreview(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryUnbondingDelegationPreviewParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	validator, found := k.stakingKeeper.GetValidator(ctx, params.ValidatorAddr)
	if !found {
		return nil, stakingTypes.ErrNoValidatorFound
	}
	delegation, found := k.stakingKeeper.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return nil, stakingTypes.ErrNoDelegatorForAddress
	}
	
	shares, err := k.stakingKeeper.ValidateUnbondAmount(
		ctx, params.DelegatorAddr, params.ValidatorAddr, params.Amount.Amount,
	)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"delAddr": params.DelegatorAddr.String(),
			"valAddr": params.ValidatorAddr.String(),
			"amt":     params.Amount.Amount.String(),
		}).Error("ValidateUnbondAmount")
		return nil, err
	}
	var removeToken sdk.Int
	source_amount := types.MustParseLedgerDec(validator.TokensFromShares(delegation.Shares)) 
	if delegation.Shares.LT(shares) {
		return nil, stakingTypes.ErrNotEnoughDelegationShares
	}
	balanceShares := delegation.Shares.Sub(shares) 

	validator, removeToken = validator.RemoveDelShares(shares) 
	//fmt.Println("token:",validator.Tokens)

	resp := types.UnbondingDelegationPreviewResponse{
		Shares:        types.MustParseLedgerDec(shares),              
		Amount:        types.MustParseLedgerDec(removeToken.ToDec()), 
		SourceAmount:  source_amount,
		SourceShares:  types.MustParseLedgerDec(delegation.Shares),
		BalanceShares: types.MustParseLedgerDec(balanceShares), 
	}

	if validator.Tokens.IsZero() {
		resp.BalanceAmount = "0" 
	} else {
		balanceAmount := validator.TokensFromShares(delegation.Shares).Sub(sdk.NewDecFromInt(removeToken)) 
		resp.BalanceAmount = types.MustParseLedgerDec(balanceAmount)                                       
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}


func queryDelegationPreview(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryDelegationPreviewParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	validator, found := k.stakingKeeper.GetValidator(ctx, params.ValidatorAddr)
	if !found {
		return nil, stakingTypes.ErrNoValidatorFound
	}

	delegation, found := k.stakingKeeper.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		delegation = stakingTypes.NewDelegation(params.DelegatorAddr, params.ValidatorAddr, sdk.ZeroDec())
	}

	var shares, balance_amount sdk.Dec

	validator, shares = validator.AddTokensFromDel(params.Amount.TruncateInt())

	balance_shares := delegation.Shares.Add(shares)
	balance_amount = validator.TokensFromSharesTruncated(balance_shares)

	resp := types.DelegationPreviewResponse{
		Shares:        types.MustParseLedgerDec(shares), 
		SourceAmount:  types.MustParseLedgerDec(validator.TokensFromShares(delegation.Shares)),
		SourceShares:  types.MustParseLedgerDec(delegation.Shares),
		BalanceAmount: types.MustParseLedgerDec(balance_amount), 
		BalanceShares: types.MustParseLedgerDec(balance_shares), 
	}
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

//pos
func queryDelegationShares(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryDelegatorParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	totalSharesStr, _ := k.GetAccountDelegatorShares(ctx, params.DelegatorAddr)
	return []byte(totalSharesStr), nil
}

//pos
func queryDelegation(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryDelegatorParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	totalSharesStr, totalBalanceStr := k.GetAccountDelegatorShares(ctx, params.DelegatorAddr)
	resp := types.ValidatorsDelegationResp{
		Shares:  totalSharesStr,
		Balance: totalBalanceStr,
	}
	return util.Json.Marshal(resp)
}

func queryDelegationFreeze(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryDelegatorParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	freezeCoin, err := k.stakingKeeper.GetDelegationFreeze(ctx, params.DelegatorAddr)
	if err != nil {
		return nil, err
	}
	realCoin := types.MustParseLedgerDec(freezeCoin)
	return []byte(realCoin), nil
}


func queryValidatorSigningInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryValidatorSigningInfoParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	validatorSignInfo, found := k.slashingKeeper.GetValidatorSigningInfo(ctx, params.ConsAddress)
	if !found {
		return nil, slashingTypes.ErrNoSigningInfoFound 
	}
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, validatorSignInfo)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

/**
，、、
*/
func queryValidatorDelegationDetail(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryBondsParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegationDetail := types.NewDelegationDetail(params.DelegatorAddr, params.ValidatorAddr)

	validator, found := k.stakingKeeper.GetValidator(ctx, params.ValidatorAddr)
	if !found {
		return nil, stakingTypes.ErrNoValidatorFound
	}

	undelegationAmount := sdk.NewDec(0)               
	delegationShareNumber := sdk.NewDec(0)            
	validatorShareNumber := validator.DelegatorShares 

	
	delegations, found := k.stakingKeeper.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if found {
		delegationShareNumber = delegations.Shares
	}
	/*delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)
	for i := 0; i < len(delegations); i++ {
		if delegations[i].DelegatorAddress.Equals(params.DelegatorAddr) {
			delegationShareNumber = delegationShareNumber.Add(delegations[i].Shares)
		}
	}*/

	
	undelegations := k.stakingKeeper.GetAllUnbondingDelegations(ctx, params.DelegatorAddr)
	for i := 0; i < len(undelegations); i++ {
		if undelegations[i].ValidatorAddress == params.ValidatorAddr.String() {
			for j := 0; j < len(undelegations[i].Entries); j++ {
				undelegationAmount = undelegationAmount.Add(undelegations[i].Entries[j].Balance.ToDec())
			}
		}
	}

	if validator.Tokens.IsZero() {
		delegationDetail.DelegationAmount = types.MustParseLedgerDec(sdk.NewDec(0))
		delegationDetail.ValidatorShareNumber = types.MustParseLedgerDec(validatorShareNumber)
		delegationDetail.UnbindingDelegationAmount = types.MustParseLedgerDec(undelegationAmount)

		res, err := codec.MarshalJSONIndent(legacyQuerierCdc, delegationDetail)
		if err != nil {
			log.WithError(err).Error("MarshalJSONIndent1")
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return res, nil
	}

	delegationDetail.DelegationShareNumber = types.MustParseLedgerDec(delegationShareNumber)
	delegationDetail.DelegationAmount = types.MustParseLedgerDec(validator.TokensFromSharesTruncated(delegationShareNumber))
	delegationDetail.ValidatorShareNumber = types.MustParseLedgerDec(validatorShareNumber)
	delegationDetail.UnbindingDelegationAmount = types.MustParseLedgerDec(undelegationAmount)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, delegationDetail)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent2")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}


func queryCopyrightExist(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryCopyrightDetailParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	_, err := k.GetCopyright(ctx, params.Hash)
	
	if err != nil {
		return []byte("non-existent"), nil
	}
	return []byte("exist"), nil
}


func queryCopyrightDetail(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryCopyrightDetailParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	dataBytes, err := k.GetCopyright(ctx, params.Hash)
	if err != nil {
		return nil, err
	}
	return dataBytes, nil
}


func queryOriginDatahashDetail(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryCopyrightDetailParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	dataBytes, err := k.QueryCopyrightOriginDataHash(ctx, params.Hash)
	if err != nil {
		return nil, err
	}
	return dataBytes, nil
}


func queryCopyrightPartyExist(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryCopyrightPartyParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	bz, err := k.GetCopyrightParty(ctx, params.Creator)
	
	if err != nil || len(bz) == 0 {
		return []byte("non-existent"), nil
	}
	return []byte("exist"), nil
}

//Id
func queryCopyrightPublishId(ctx sdk.Context, k Keeper) ([]byte, error) {
	//MAP,,.json,
	return k.QueryPublisherIdMap(ctx), nil
}


func queryCopyrightParty(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryCopyrightPartyParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	dataBytes, err := k.GetCopyrightParty(ctx, params.Creator)
	if err != nil {
		return nil, err
	}
	return dataBytes, nil
}


func queryBlockRDS(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryBlockRDSParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	dataBytes, err := k.GetBlockRDS(ctx, params.Height)
	if err != nil {
		return nil, err
	}
	return dataBytes, nil
}


func queryAccountSpaceMiner(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAccountSpaceMinerParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	accountSpaceMiner := k.QueryAccountSpaceMinerInfor(ctx, params.Account)
	return util.Json.Marshal(accountSpaceMiner)
}

func queryDeflationRateInfor(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	deflationRate := keeper.QueryDeflationRate(ctx)
	rateAndVoteIndex := RateAndVoteIndex{
		DeflationRate: deflationRate,
	}
	return util.Json.Marshal(rateAndVoteIndex)
}


func querySpacefee(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	spaceTotal := keeper.SpaceFeeEstimate(ctx)
	return []byte(spaceTotal), nil
}


func spaceTotal(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	spaceTotal := keeper.QueryDeflatinSpaceTotal(ctx)
	return []byte(spaceTotal), nil
}


func spaceAward(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	spaceTotal := keeper.QuerySpaceMinerBonusAmount(ctx)
	return []byte(spaceTotal.String()), nil
}


func hasMinerBonus(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	hasMinerBonus := keeper.QueryHasMinerBonusAmount(ctx)
	return []byte(hasMinerBonus.String()), nil
}


func spaceAmount(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	spaceAmount := keeper.QuerySpaceMinerAmount(ctx)
	return util.Json.Marshal(spaceAmount)
}

func queryDeflationMinerInfor(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	deflationMinerInfor, err := keeper.QueryDeflationMinerInfor(ctx)
	
	res, err := util.Json.Marshal(deflationMinerInfor)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryNftInfor(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryNftInforParams
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	deflationMinerInfor := keeper.QueryNftInfor(ctx, params.TokenId)
	
	res, err := util.Json.Marshal(deflationMinerInfor)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
