package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	
	QueryCopyrightDetail = "copyright_detail"
	
	QueryPubCount = "pub_count"
	
	QueryBonusExtrainfor = "copyright_bonus_extrainfor"
	
	QueryCopyrightExist = "copyright_exist"
	//hash
	QueryOriginHash = "copyright_origin_hash"
	
	QueryCopyrightParty = "copyright_party"
	
	QueryCopyrightPublishId = "copyright_publish_id"
	
	QueryCopyrightPartyExist = "copyright_party_exist"
	
	QueryBlockRDS = "copyright_block_relation"
	
	QueryAccountSpaceMiner = "account_space_miner"
	
	QueryTotalSpaceInfor = "total_space_infor"
	
	QueryDeflationRateInfor = "deflation_rate_infor"
	
	QuerySpaceAmount = "space_amount"
	
	QuerySpaceAward = "space_award"
	
	QueryHasMinerBonus = "has_miner_bonus"
	
	QuerySpaceFee = "query_space_fee"
	
	QueryDeflationMinerInfor = "deflation_miner_infor"
	//Nft
	QueryNftInfor = "nft_infor"
	//，、、
	QueryValidatorDelegationDetail = "validatorDelegationsDetail"
	
	QuerySigningInfo = "signingInfo"
	
	QueryDelegationShares = "queryDelegationShares"
	
	QueryDelegation = "queryDelegation"
	//POS，
	QueryDelegationPreview = "delegationPreview"
	//POS，
	QueryUnbondingDelegationPreview = "unbondingDelegationPreview"
	
	QueryTotalShares = "totalShares"
	
	QueryDelegationByConsAddress = "delegationByConsAddress"
	
	QueryValidatorByConsAddress = "validatorByConsAddress"
	
	QueryInviteRecord = "inviteRecord"
	
	QueryInviteStatistics = "inviteStatistics"
	
	QueryMortgAmount = "mortg_amount"
	
	QueryMiningStage = "mining_stage"
	
	QueryMortgMiningInfor = "mortg_miner_infor"
	
	QueryCopyrightComplain = "copyright_complain"
	
	QueryRewardInfo = "reward_info"
	
	QuerySpaceMinerRewardInfo = "space_miner_reward_info"
	
	QueryValidatorInfo = "validator_info"
	
	QueryDelegationFreeze = "delegation_freeze"
	
	QueryCopyrightExport = "copyright_export"
	
	QueryParams = "copyright_params"
)


type QueryValidatorParams struct {
	Status string `json:"status"`
}


type QueryAccountParams struct {
	Account string `json:"account"`
}


type QueryInviteRecordParams struct {
	InviteAddress string `json:"invite_address"`
}


type QueryCopyrightDetailParams struct {
	Hash string `json:"hash"`
}


type QueryCopyrightPartyParams struct {
	Creator string `json:"creator"` 
}


type QueryBlockRDSParams struct {
	Height int64 `json:"height"` 
}


type QueryAccountSpaceMinerParams struct {
	Account string `json:"account"` 
}

//nft
type QueryNftInforParams struct {
	TokenId string `json:"token_id"` //nftid
}


type QueryCopyrightComplainParams struct {
	ComplainId string `json:"complain_id"` //id
}


type QueryAuthorizeAccountParams struct {
	Pubkey string `json:"pubkey"` 
}


type QueryAuthorizeValidatorParams struct {
	ConsAddr string `json:"cons_addr"` 
}


type QueryPubCountParams struct {
	DayString string `json:"day_string"` 
}


type QueryResourceAndHashRelationParams struct {
	Account string `json:"account"` 
	Hash    string `json:"hash"`    //hash
}

type QueryValidatorByConsAddrParams struct {
	ValidatorConsAddress sdk.ConsAddress
}

type QueryDelegatorByConsAddrParams struct {
	DelegatorAddr        sdk.AccAddress
	ValidatorConsAddress sdk.ConsAddress
}

func NewQueryDelegatorByConsAddrParams(delegatorAddr sdk.AccAddress, validatorConsAddress sdk.ConsAddress) QueryDelegatorByConsAddrParams {
	return QueryDelegatorByConsAddrParams{
		DelegatorAddr:        delegatorAddr,
		ValidatorConsAddress: validatorConsAddress,
	}
}

type QueryDelegationPreviewParams struct {
	DelegatorAddr sdk.AccAddress
	ValidatorAddr sdk.ValAddress
	Amount        sdk.Dec
}

func NewQueryDelegationPreviewParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, addAmount sdk.Dec) QueryDelegationPreviewParams {
	return QueryDelegationPreviewParams{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
		Amount:        addAmount,
	}
}

type QueryUnbondingDelegationPreviewParams struct {
	DelegatorAddr sdk.AccAddress
	ValidatorAddr sdk.ValAddress
	Amount        sdk.Coin
}

func NewQueryUnbondingDelegationPreviewParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, amount sdk.Coin) QueryUnbondingDelegationPreviewParams {
	return QueryUnbondingDelegationPreviewParams{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
		Amount:        amount,
	}
}

type QueryDelegatorParams struct {
	DelegatorAddr sdk.AccAddress
}

type QueryValidatorSigningInfoParams struct {
	ConsAddress sdk.ConsAddress
}

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

type QueryBondsParams struct {
	DelegatorAddr sdk.AccAddress
	ValidatorAddr sdk.ValAddress
}

func NewQueryBondsParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) QueryBondsParams {
	return QueryBondsParams{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
	}
}
