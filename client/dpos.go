package client

import (
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/util/spaceutil"
	"fs.video/blockchain/x/copyright/client/rest"
	"fs.video/blockchain/x/copyright/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"strings"
	"time"

	//"fs.video/blockchain/x/copyright/types"
	"errors"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	//sdk "github.com/cosmos/cosmos-sdk/types"
)

/**
dpos
*/

type DposClient struct {
	TxClient  *TxClient
	ServerUrl string
	logPrefix string
}


func (this *DposClient) GetValidatorStatus(status stakingTypes.BondStatus, jailed bool) int {
	if jailed {
		return 3
	} else {
		switch status.String() {
		case "BOND_STATUS_UNBONDED":
			return 0 
		case "BOND_STATUS_UNBONDING":
			return 1 
		case "BOND_STATUS_BONDED":
			return 2 
		}
	}
	return 0
}


func (this *DposClient) FindValidatorRegisterLimit() (limit *types.ValidatorRegisterLimit) {
	diskSize, _ := spaceutil.GetDiskInfo()
	//fmt.Println(":", (diskSize / 1024 / 1024 / 1024), "G")
	//fmt.Println(":", (diskFree / 1024 / 1024 / 1024), "G")
	limit = &types.ValidatorRegisterLimit{
		Status: "0",
	}
	diskSpace := (diskSize / 1024 / 1024 / 1024)
	/*diskMortgInt := sdk.FormatIntMul2Int(sdk.NewInt(int64(diskSpace)))
	mortgAmount := sdk.NewCoin(sdk.DefaultBondDenom,diskMortgInt)
	sdk.FormatCoin2Stand(mortgAmount)*/
	//staking
	limit.Status = "1"
	limit.MortgAmount = sdk.NewInt(int64(diskSpace)).Mul(sdk.NewInt(int64(core.DiskSpaceMortgRate))).String()
	return limit
}

//pos
func (this *DposClient) MinerReportForm(ben32DelegatorAddr, ben32ValidatorAddr string) (delegation *types.PosReportFormResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(ben32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return
	}
	_, err = sdk.ValAddressFromBech32(ben32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		return
	}
	totalShares, err := this.FindTotalShares()
	if err != nil {
		return
	}
	delegationDetail, err := this.FindDelegationDetail(ben32DelegatorAddr, ben32ValidatorAddr)
	if err != nil {
		if strings.Contains(err.Error(), "validator does not exist") {
			reportForm := types.PosReportFormResponse{
				MortgAmount:         "0",         
				Shares:              "0",         
				UnbondAmount:        "0",         
				ValidatorShares:     "0",         
				PosRewardUnreceived: "0",         //pos
				TotalShares:         totalShares, 
				AccountTotalShares:  "0",         
			}
			return &reportForm, nil
		} else {
			return nil, err
		}
	}
	
	accountTotalShares, err := this.FindDelegationShares(ben32DelegatorAddr)
	if err != nil {
		return
	}
	
	validatorReward, delegationReward, err := this.RewardsPreview(ben32DelegatorAddr, ben32ValidatorAddr)
	if err != nil {
		if !strings.Contains(err.Error(), distributionTypes.ErrNoDelegationExists.Error()) {
			log.WithError(err).Error("RewardsPreview")
			return
		}
	}

	posRewardUnreceived := sdk.MustNewDecFromStr(validatorReward.Amount).Add(sdk.MustNewDecFromStr(delegationReward.Amount))
	reportForm := &types.PosReportFormResponse{
		MortgAmount:         delegationDetail.DelegationAmount,          
		Shares:              delegationDetail.DelegationShareNumber,     
		UnbondAmount:        delegationDetail.UnbindingDelegationAmount, 
		ValidatorShares:     delegationDetail.ValidatorShareNumber,      
		PosRewardUnreceived: posRewardUnreceived.String(),               //pos
		TotalShares:         totalShares,                                
		AccountTotalShares:  accountTotalShares,                         
	}
	return reportForm, nil
}

/**

minSelfDelegation : 
comm :
*/
func (this *DposClient) RegisterValidator(bech32DelegatorAddr, bech32ValidatorAddr string, bech32ValidatorPubkey cryptotypes.PubKey, selfDelegation sdk.Coin, desc stakingTypes.Description, commission stakingTypes.CommissionRates, minSelfDelegation sdk.Int, privateKey string, fee decimal.Decimal) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	//validatorPubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, bech32ValidatorPubkey)
	//validatorPubkey, err := legacybech32.UnmarshalPubKey(legacybech32.ConsPK, bech32ValidatorPubkey)
	//if err != nil {
	//	log.WithError(err).Error("GetPubKeyFromBech32")
	//	err = errors.New(rest.ParseAccountError)
	//	return
	//}
	
	if selfDelegation.Amount.LT(sdk.NewInt(core.ValidatorRegisterMinAmount)) {
		err = errors.New(rest.ValidatorMortgMinAmount)
		return
	}
	_, err = desc.EnsureLength()
	if err != nil {
		log.WithError(err).Error("EnsureLength")
		err = errors.New(rest.ValidatorDescriptionError)
		return
	}
	
	balance, err := this.TxClient.Balance(bech32DelegatorAddr, core.MainToken)
	if err != nil {
		return
	}
	balanceCoin := types.MustRealCoin2LedgerCoin(*balance) 
	if balanceCoin.Amount.LTE(selfDelegation.Amount) {
		err = errors.New(rest.AccountInsufficient) 
		return
	}
	
	existed, err := this.FindValidatorExisted(bech32ValidatorAddr)
	if err != nil {
		err = errors.New(rest.QueryChainInforError)
		return
	}

	
	if existed {
		err = errors.New(rest.ValidatorExist)
		return
	}

	
	//if minSelfDelegation.IsZero() {
	//	err = errors.New(rest.MinSelfDelegationZeroError)
	//	return
	//}

	
	err = commission.Validate()
	if err != nil {
		log.WithError(err).Error("commission.Validate")
		return
	}
	//commission := stakingTypes.NewCommissionRates(sdk.MustNewDecFromStr("0.20"), sdk.OneDec(), sdk.MustNewDecFromStr("0.10")) 
	minSelfDelegation = types.MustRealInt2LedgerInt(minSelfDelegation) 
	msg, err := stakingTypes.NewMsgCreateValidator(validatorAddr, bech32ValidatorPubkey, selfDelegation, desc, commission, minSelfDelegation)
	if err != nil {
		log.WithError(err).Error("NewMsgCreateValidator")
		return
	}

	resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, types.NewLedgerFee(fee), "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

/**

*/
func (this *DposClient) EditorValidator(bech32ValidatorAccAddr string, desc stakingTypes.Description, newRate sdk.Dec, minSelfDelegation sdk.Int, privateKey string, fee decimal.Decimal) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	/*_, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, bech32ValidatorOprateAddr)
	if err != nil {
		err = errors.New(rest.ParseAccountError)
		return
	}*/
	accAddr, err := sdk.AccAddressFromBech32(bech32ValidatorAccAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return
	}
	validatorAddress := sdk.ValAddress(accAddr).String()
	//validatorAddress := sdk.ValAddress(accAddr).String()
	
	validatorInfor, err := this.FindValidatorByValAddress(validatorAddress)
	if err != nil {
		err = errors.New(rest.QueryChainInforError)
		return
	}
	
	if validatorInfor.GetOperator().Empty() {
		err = errors.New(rest.ValidatorExist)
		return
	}

	
	//if minSelfDelegation.IsZero() {
	//	err = errors.New(rest.MinSelfDelegationZeroError)
	//	return
	//}
	minSelfDelegation = types.MustRealInt2LedgerInt(minSelfDelegation) 
	if !minSelfDelegation.GTE(validatorInfor.MinSelfDelegation) {
		return nil, stakingTypes.ErrMinSelfDelegationDecreased
	}

	if minSelfDelegation.GT(validatorInfor.Tokens) {
		return nil, stakingTypes.ErrSelfDelegationBelowMinimum
	}

	
	//err = validatorInfor.Commission.ValidateNewRate(validatorInfor.Commission.Rate, time.Now())
	if time.Now().Sub(validatorInfor.Commission.UpdateTime).Hours() < 24 {
		err = errors.New(rest.ValidatorInfoError)
		return
	}
	//log.Debug( "", validatorInfor.OperatorAddress)
	//log.Debug( "", validatorInfor.OperatorAddress)
	/*if validatorInfor.OperatorAddress == bech32ValidatorOprateAddr {
		err = errors.New(rest.NoRightOprateValidator)
		return
	}*/
	_, err = desc.EnsureLength()
	if err != nil {
		log.WithError(err).Error("EnsureLength")
		err = errors.New(rest.ValidatorDescriptionError)
		return
	}
	
	balance, err := this.TxClient.Balance(bech32ValidatorAccAddr, core.MainToken)
	if err != nil {
		return
	}
	decimalBalance, err := decimal.NewFromString(balance.Amount)
	if decimalBalance.LessThan(fee) {
		err = errors.New(rest.AccountInsufficient) 
		return
	}

	msg := stakingTypes.NewMsgEditValidator(validatorInfor.GetOperator(), desc, &newRate, &minSelfDelegation)
	if err != nil {
		log.WithError(err).Error("NewMsgEditValidator")
		return
	}
	resp, err = this.TxClient.SignAndSendMsg(bech32ValidatorAccAddr, privateKey, types.NewLedgerFee(fee), "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

/**

*/
func (this *DposClient) Delegation(bech32DelegatorAddr, bech32ValidatorAddr string, amount sdk.Coin, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	msg := stakingTypes.NewMsgDelegate(delegatorAddr, validatorAddr, amount)
	fee := types.NewLedgerFee(core.ChainDefaultFee)
	resp, err = this.TxClient.SignAndSendMsg(msg.DelegatorAddress, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

/**

*/
func (this *DposClient) FindDelegation(delegatorAddr, validatorAddr string) (delegation *stakingTypes.DelegationResponse, notFound bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	notFound = false
	params := stakingTypes.QueryDelegatorValidatorRequest{DelegatorAddr: delegatorAddr, ValidatorAddr: validatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, notFound, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegation, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		if strings.Contains(err.Error(), stakingTypes.ErrNoDelegation.Error()) {
			notFound = true
		}
		return nil, notFound, err
	}
	delegation = &stakingTypes.DelegationResponse{}
	err = util.Json.Unmarshal(resBytes, delegation)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}

/**

*/
func (this *DposClient) FindDelegationByConsAddress(bech32DelegatorAddr string, bech32ConsAddr string) (delegation *stakingTypes.Delegation, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	consAddr, err := sdk.ConsAddressFromBech32(bech32ConsAddr)
	if err != nil {
		log.WithError(err).Error("ConsAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryDelegatorByConsAddrParams{DelegatorAddr: delegatorAddr, ValidatorConsAddress: consAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryDelegationByConsAddress, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	delegation = &stakingTypes.Delegation{}
	err = util.Json.Unmarshal(resBytes, delegation)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (bc *DposClient) FindValidatorExisted(bech32ValidatorAddr string) (existed bool, err error) {
	_, err = bc.FindValidatorByValAddress(bech32ValidatorAddr)
	if err != nil {
		
		if strings.Contains(err.Error(), stakingTypes.ErrNoValidatorFound.Error()) {
			err = errors.New(rest.QueryChainInforError)
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

/**
 
*/
func (this *DposClient) FindValidatorByValAddress(bech32ValidatorAddr string) (validator *stakingTypes.Validator, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := stakingTypes.QueryValidatorParams{ValidatorAddr: validatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryValidator, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	validator = &stakingTypes.Validator{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, validator)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}

/**
 
*/
func (this *DposClient) FindValidatorByValAddressWithHeight(bech32ValidatorAddr string, height int64) (validator *stakingTypes.Validator, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := stakingTypes.QueryValidatorParams{ValidatorAddr: validatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	clientCtx = clientCtx.WithHeight(height)
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryValidator, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	validator = &stakingTypes.Validator{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, validator)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}

/**
 
*/
func (this *DposClient) FindValidatorByConsAddress(bech32ConsAddr string) (validator *stakingTypes.Validator, notFound bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	notFound = false
	consAddress, err := sdk.ConsAddressFromBech32(bech32ConsAddr)
	if err != nil {
		log.WithError(err).Error("ConsAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryValidatorByConsAddrParams{ValidatorConsAddress: consAddress}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, notFound, err
	}

	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryValidatorByConsAddress, bz)
	if err != nil {
		
		if strings.Contains(err.Error(), stakingTypes.ErrNoValidatorFound.Error()) {
			notFound = true
			err = nil 
		} else {
			log.WithError(err).Error("QueryWithData")
		}
		return nil, notFound, err
	}
	validator = &stakingTypes.Validator{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, validator)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}

/**

*/
func (this *DposClient) ListValidator(status string, page, limit int) (validators *[]stakingTypes.Validator, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := stakingTypes.QueryValidatorsParams{Status: status, Page: page, Limit: limit}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryValidators, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	validators = &[]stakingTypes.Validator{}
	err = util.Json.Unmarshal(resBytes, validators)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}

/**

*/
func (this *DposClient) FindValidatorsName(addrs []sdk.ValAddress) (validators map[string]string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := stakingTypes.NewQueryValidatorsNameParams(addrs)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryValidatorsName, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	validators = make(map[string]string)
	err = util.Json.Unmarshal(resBytes, &validators)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}

/**

*/
func (this *DposClient) Validators(status string, page, limit int) (validators []types.ValidatorInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := stakingTypes.QueryValidatorsParams{Status: status, Page: page, Limit: limit}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryValidatorInfo, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	validators = []types.ValidatorInfo{}
	err = util.Json.Unmarshal(resBytes, &validators)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}

/**

*/
func (this *DposClient) FindDelegationDetail(bech32DelegatorAddr, bech32ValidatorAddr string) (delegation *types.DelegationDetail, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryBondsParams{DelegatorAddr: delegatorAddr, ValidatorAddr: validatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryValidatorDelegationDetail, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	delegation = &types.DelegationDetail{}
	err = util.Json.Unmarshal(resBytes, delegation)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (this *DposClient) ListUnbondingDelegation(bech32DelegatorAddr string) (list []stakingTypes.UnbondingDelegation, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryDelegatorParams{DelegatorAddr: delegatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegatorUnbondingDelegations, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	list = []stakingTypes.UnbondingDelegation{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &list)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}

/**

*/
func (this *DposClient) ListDelegationByValidator(bech32ValAddr string, page, limit int) (resObject stakingTypes.DelegationResponses, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	valAddr, err := sdk.ValAddressFromBech32(bech32ValAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := stakingTypes.NewQueryValidatorParams(valAddr, page, limit)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryValidatorDelegations, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	resObject = stakingTypes.DelegationResponses{}
	err = util.Json.Unmarshal(resBytes, &resObject)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (this *DposClient) ListBondingDelegation(bech32DelegatorAddr string) (list []stakingTypes.DelegationResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryDelegatorParams{DelegatorAddr: delegatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegatorDelegations, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	list = []stakingTypes.DelegationResponse{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &list)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}


func (this *DposClient) FindTotalShares() (total string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryTotalShares, []byte{})
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return total, err
	}
	total = string(resBytes)
	return
}


func (this *DposClient) FindDelegationShares(bech32DelegatorAddr string) (total string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryDelegatorParams{DelegatorAddr: delegatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return total, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryDelegationShares, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return total, err
	}
	total = string(resBytes)
	return
}


func (this *DposClient) FindAllDelegations(bech32DelegatorAddr string) (resp types.ValidatorsDelegationResp, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryDelegatorParams{DelegatorAddr: delegatorAddr}
	resp = types.ValidatorsDelegationResp{}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return resp, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryDelegation, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = util.Json.Unmarshal(resBytes, &resp)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (this *DposClient) FindDelegationFreeze(bech32DelegatorAddr string) (total string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryDelegatorParams{DelegatorAddr: delegatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return total, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryDelegationFreeze, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return total, err
	}
	total = string(resBytes)
	return
}


func (this *DposClient) DelegationPreview(bech32DelegatorAddr string, bech32ValidatorAddr string, amount sdk.Dec) (resp *types.DelegationPreviewResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryDelegationPreviewParams{DelegatorAddr: delegatorAddr, ValidatorAddr: validatorAddr, Amount: amount}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryDelegationPreview, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	resp = &types.DelegationPreviewResponse{}
	err = util.Json.Unmarshal(resBytes, resp)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}

//  amount:
func (this *DposClient) UnbondDelegation(bech32DelegatorAddr string, bech32ValidatorAddr string, amount sdk.Coin, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	msg := stakingTypes.NewMsgUndelegate(delegatorAddr, validatorAddr, amount)
	fee := types.NewLedgerFee(core.ChainDefaultFee)
	resp, err = this.TxClient.SignAndSendMsg(msg.DelegatorAddress, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) UnbondDelegationAll(bech32DelegatorAddr, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := stakingTypes.NewQueryDelegatorParams(delegatorAddr)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegatorDelegations, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	delegationResp := stakingTypes.DelegationResponses{}
	err = util.Json.Unmarshal(resBytes, &delegationResp)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return
	}
	if len(delegationResp) <= 0 {
		err = errors.New("delegation does not exist")
		return
	}
	msgs := []sdk.Msg{}
	for _, delegation := range delegationResp {
		valAddr, err := sdk.ValAddressFromBech32(delegation.Delegation.ValidatorAddress)
		if err != nil {
			continue
		}
		msg := stakingTypes.NewMsgUndelegate(delegatorAddr, valAddr, delegation.Balance)
		msgs = append(msgs, msg)
	}

	fee := types.NewLedgerFee(core.ChainDefaultFee)
	resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msgs...)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

//   amount:
func (this *DposClient) UnbondDelegationPreview(bech32DelegatorAddr, bech32ValidatorAddr string, amount sdk.Coin) (res *types.UnbondingDelegationPreviewResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.NewQueryUnbondingDelegationPreviewParams(delegatorAddr, validatorAddr, amount)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryUnbondingDelegationPreview, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	res = &types.UnbondingDelegationPreviewResponse{}
	err = util.Json.Unmarshal(resBytes, res)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (this *DposClient) FindSlashingParams() (resp *slashingTypes.Params, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := clientCtx.QueryWithData("custom/slashing/"+slashingTypes.QueryParameters, []byte{})
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	resp = &slashingTypes.Params{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, resp)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}


func (this *DposClient) FindSigningInfo(ben32ConsAddr string) (resp *slashingTypes.ValidatorSigningInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	consAddress, err := sdk.ConsAddressFromBech32(ben32ConsAddr)
	if err != nil {
		log.WithError(err).Error("ConsAddressFromBech32")
		return nil, err
	}
	params := types.QueryValidatorSigningInfoParams{ConsAddress: consAddress}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QuerySigningInfo, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	resp = &slashingTypes.ValidatorSigningInfo{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, resp)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}


func (this *DposClient) UnjailValidator(bech32DelegatorAddr, bech32ValidatorAddr, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	
	validatorInfo, err := this.FindValidatorByValAddress(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("FindValidatorByValAddress")
		err = errors.New(rest.QueryChainInforError)
		return
	}
	if !validatorInfo.Jailed {
		err = errors.New(rest.ValidatorUnJail)
		return
	}

	accAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.QueryChainInforError)
		return
	}
	OperatorAddress := sdk.ValAddress(accAddr).String()
	if validatorInfo.OperatorAddress != OperatorAddress {
		err = errors.New(rest.NoRightOprateValidator)
		return
	}
	
	delegatorResponse, notFound, err := this.FindDelegation(bech32DelegatorAddr, bech32ValidatorAddr)
	if err != nil {
		if notFound {
			err = errors.New(rest.UnjailAccountNOtDelegation)
		} else {
			err = errors.New(rest.QueryChainInforError)
		}
		return
	}
	
	if delegatorResponse.Delegation.Shares.IsZero() {
		err = errors.New(rest.UnjailAccountNOtDelegation)
		return
	}

	tokens := validatorInfo.TokensFromShares(delegatorResponse.Delegation.Shares).TruncateInt()
	if tokens.LT(validatorInfo.MinSelfDelegation) {
		err = errors.New(rest.UnjailAccountNOtEnoughDelegation)
		return
	}

	
	msg := slashingTypes.NewMsgUnjail(validatorAddr)
	fee := types.NewLedgerFee(core.ChainDefaultFee)
	resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
	return nil, nil
}

/**

*/
func (this *DposClient) CommunityRewardPreview() (realCoins *types.RealCoins, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryCommunityPool, []byte{})
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	decCoins := sdk.DecCoins{}
	err = util.Json.Unmarshal(resBytes, &decCoins)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return nil, err
	}
	realCoins1 := types.MustLedgerDecCoins2RealCoins(decCoins)
	return &realCoins1, nil
}

/**

*/
func (this *DposClient) CommunityRewardDraw(bech32AccountAddr, privateKey string, coins sdk.Coins) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	accountAddr, err := sdk.AccAddressFromBech32(bech32AccountAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	msg, err := types.NewMsgDistributeCommunityReward(accountAddr, coins.String())
	if err != nil {
		log.WithError(err).Error("NewMsgDistributeCommunityReward")
		return
	}

	fee := types.NewLedgerFee(core.ChainDefaultFee)
	resp, err = this.TxClient.SignAndSendMsg(bech32AccountAddr, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

/**

delegatorReward 
validatorReward 
*/
func (this *DposClient) RewardsPreview(bech32DelegatorAddr string, bech32ValidatorAddr string) (delegatorReward types.RealCoin, validatorReward types.RealCoin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	/**
	1.
	2.
	3.1
	(ValidatorsDistributionTotal)
	3.2
	(ValidatorsDistributionReWards2)
	*/
	delegatorReward = types.RealCoin{Denom: core.MainToken, Amount: "0"}
	validatorReward = types.RealCoin{Denom: core.MainToken, Amount: "0"}
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return delegatorReward, validatorReward, errors.New(rest.ParseAccountError)
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		return delegatorReward, validatorReward, errors.New(rest.ParseAccountError)
	}
	params := distributionTypes.NewQueryDelegationRewardsParams(delegatorAddr, validatorAddr)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return delegatorReward, validatorReward, errors.New(rest.QueryChainInforError)
	}
	
	resBytes, _, err := clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryDelegationRewards, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData error1 | ", err.Error())
		return delegatorReward, validatorReward, err
	}

	delegatorRewardDecCoins := sdk.DecCoins{} 

	err = util.Json.Unmarshal(resBytes, &delegatorRewardDecCoins)
	if err != nil {
		log.WithError(err).Error("Unmarshal error1 | ", err.Error())
		return delegatorReward, validatorReward, errors.New(rest.QueryChainInforError)
	}

	if len(delegatorRewardDecCoins) != 0 {
		delegatorReward = types.MustLedgerDecCoins2RealCoins(delegatorRewardDecCoins)[0] 
	}

	bech32DelegatorValidatorAddr := sdk.ValAddress(delegatorAddr).String()

	
	if bech32DelegatorValidatorAddr == bech32ValidatorAddr {
		
		resBytes, _, err = clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryValidatorCommission, bz)
		if err != nil {
			log.WithError(err).Error("QueryWithData error2 | ", err.Error())
			return delegatorReward, validatorReward, errors.New(rest.QueryChainInforError)
		}
		validatorCommDecCoins := distributionTypes.ValidatorAccumulatedCommission{} 
		err = util.Json.Unmarshal(resBytes, &validatorCommDecCoins)
		if err != nil {
			log.WithError(err).Error("Unmarshal error2 | ", err.Error())
			return delegatorReward, validatorReward, errors.New(rest.QueryChainInforError)
		}
		if len(validatorCommDecCoins.Commission) > 0 {
			defaultMoney := validatorCommDecCoins.Commission[0]

			
			if defaultMoney.Amount.GTE(sdk.NewDec(core.MinLedgerAmountInt64)) { 
				validatorReward = types.MustLedgerDecCoin2RealCoin(defaultMoney) 
			}
		}
	}

	return delegatorReward, validatorReward, err
}


func (this *DposClient) DrawCommissionDelegationRewards(bech32DelegatorAddr, bech32ValidatorAddr, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	bech32DelegatorValidatorAddr := sdk.ValAddress(delegatorAddr).String()

	
	delegationReward, validatorReward, err := this.RewardsPreview(bech32DelegatorAddr, bech32ValidatorAddr)
	if err != nil {
		return
	}

	//fmt.Println("validatorReward:", validatorReward, "delegationReward:", delegationReward)

	if validatorReward.Amount == "0" && delegationReward.Amount == "0" {
		err = errors.New("There is no reward to receive") 
		return
	}
	msgs := []sdk.Msg{}

	if delegationReward.Amount != "0" { 
		msg1 := distributionTypes.NewMsgWithdrawDelegatorReward(delegatorAddr, validatorAddr)
		msgs = append(msgs, msg1)
	}
	
	if bech32DelegatorValidatorAddr == bech32ValidatorAddr && validatorReward.Amount != "0" { 
		
		msg2 := distributionTypes.NewMsgWithdrawValidatorCommission(validatorAddr)
		msgs = append(msgs, msg2) //tx
	}

	fee := types.NewLedgerFee(core.ChainDefaultFee)
	resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msgs...)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) RewardsPreviewAll(bech32DelegatorAddr string) (distributionTypes.QueryDelegatorTotalRewardsResponse, types.ValidatorCommissionResp, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	totalRewards := distributionTypes.QueryDelegatorTotalRewardsResponse{}
	var ValAddressArr types.ValidatorCommissionResp
	var valAddrcommission types.ValidatorCommission
	ValAddressArr.Total = sdk.NewDecCoin(core.MainToken, sdk.NewInt(0))
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return totalRewards, ValAddressArr, nil
	}
	params := distributionTypes.NewQueryDelegatorParams(delegatorAddr)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON error1 | ", err.Error())
		return totalRewards, ValAddressArr, nil
	}
	
	resBytes, _, err := clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryDelegatorTotalRewards, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData error1 | ", err.Error())
		return totalRewards, ValAddressArr, nil
	}

	err = util.Json.Unmarshal(resBytes, &totalRewards)
	if err != nil {
		log.WithError(err).Error("Unmarshal error1 | ", err.Error())
		return totalRewards, ValAddressArr, nil
	}
	for _, reward := range totalRewards.Rewards {
		validatorAddr, err := sdk.ValAddressFromBech32(reward.ValidatorAddress)
		if err != nil {
			log.WithError(err).Error("ValAddressFromBech32 error2 | ", err.Error())
			continue
		}
		params := distributionTypes.NewQueryDelegationRewardsParams(delegatorAddr, validatorAddr)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			log.WithError(err).Error("MarshalJSON error2 | ", err.Error())
			continue
		}
		bech32DelegatorValidatorAddr := sdk.ValAddress(delegatorAddr).String()
		
		if bech32DelegatorValidatorAddr == validatorAddr.String() {
			
			resBytes, _, err = clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryValidatorCommission, bz)
			if err != nil {
				log.WithError(err).Error("QueryWithData error2 | ", err.Error())
				continue
			}
			validatorComm := distributionTypes.ValidatorAccumulatedCommission{} 
			err = util.Json.Unmarshal(resBytes, &validatorComm)
			if err != nil {
				log.WithError(err).Error("Unmarshal error2 | ", err.Error())
				continue
			}
			if validatorComm.Commission.IsZero() {
				continue
			}
			valAddrcommission.ValidatorAddress = validatorAddr
			valAddrcommission.Reward = validatorComm.Commission[0]
			ValAddressArr.Total = ValAddressArr.Total.Add(validatorComm.Commission[0])
			ValAddressArr.ValidatorCommissions = append(ValAddressArr.ValidatorCommissions, valAddrcommission)
		}

	}
	return totalRewards, ValAddressArr, nil
}


func (this *DposClient) DrawCommissionDelegationRewardsAll(bech32DelegatorAddr, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	totalRewards, ValAddressArr, err := this.RewardsPreviewAll(bech32DelegatorAddr)
	if err != nil {
		err = errors.New(rest.ParseAccountError)
		return
	}
	if totalRewards.Total.IsZero() && len(ValAddressArr.ValidatorCommissions) <= 0 {
		err = errors.New("There is no reward to receive") 
		return
	}
	msgs := []sdk.Msg{}

	if !totalRewards.Total.IsZero() { 
		for _, reward := range totalRewards.Rewards {
			//log.Debug(":", reward.Reward.String())
			if reward.Reward.IsZero() {
				continue
			}
			if reward.Reward[0].Amount.LT(sdk.NewDec(1)) {
				continue
			}
			validatorAddr, err := sdk.ValAddressFromBech32(reward.ValidatorAddress)
			if err != nil {
				continue
			}
			msg1 := distributionTypes.NewMsgWithdrawDelegatorReward(delegatorAddr, validatorAddr)
			msgs = append(msgs, msg1)
		}
	}
	
	if len(ValAddressArr.ValidatorCommissions) > 0 { 
		
		for _, valAddress := range ValAddressArr.ValidatorCommissions {
			msg2 := distributionTypes.NewMsgWithdrawValidatorCommission(valAddress.ValidatorAddress)
			msgs = append(msgs, msg2) //tx
		}
	}

	fee := types.NewLedgerFee(core.ChainDefaultFee)
	resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msgs...)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) DrawCommissionRewards(bech32DelegatorAddr, bech32ValidatorAddr, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	msg := distributionTypes.NewMsgWithdrawValidatorCommission(validatorAddr)
	fee := types.NewLedgerFee(core.ChainDefaultFee)
	resp, err = this.TxClient.SignAndSendMsg(delegatorAddr.String(), privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) DrawDelegationRewards(bech32DelegatorAddr, bech32ValidatorAddr, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	msg := distributionTypes.NewMsgWithdrawDelegatorReward(delegatorAddr, validatorAddr)
	resp, err = this.TxClient.SignAndSendMsg(msg.DelegatorAddress, privateKey, types.NewLedgerFee(core.ChainDefaultFee), "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) FindFeePool() (sdk.DecCoins, error) {
	dc := sdk.DecCoins{}
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	resBytes, _, err := clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryCommunityPool, []byte{})
	if err != nil {
		log.WithError(err).Error("QueryWithData error2 | ", err.Error())
		return dc, err
	}
	err = util.Json.Unmarshal(resBytes, &dc)
	if err != nil {
		log.WithError(err).Error("Unmarshal error2 | ", err.Error())
		return dc, err
	}
	return dc, err
}
