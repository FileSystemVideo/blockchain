package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/blockchainchat module sentinel errors
var (
	CopyrightNotFoundErr         = sdkerrors.Register(ModuleName, 101, "copyright not exist")
	CopyrightPartyNotFoundErr    = sdkerrors.Register(ModuleName, 103, "Copyright information not found")
	RDSNotFoundErr               = sdkerrors.Register(ModuleName, 104, "RDS information not found")
	DeflationVoteOptionErr       = sdkerrors.Register(ModuleName, 105, "Illegal deflation voting option")
	DeflationVoted               = sdkerrors.Register(ModuleName, 106, "Already voted")
	TokenidNotExist              = sdkerrors.Register(ModuleName, 107, "NFT does not exist")
	TokenidHasNoRight            = sdkerrors.Register(ModuleName, 108, "You do not have permission to operate on the current NFT")
	TokenidFormatErr             = sdkerrors.Register(ModuleName, 109, "NFT data formatting failed")
	MortgMinerHasFinish          = sdkerrors.Register(ModuleName, 110, "mortg miner has finish")
	ErrDataHashDoesNotExist      = sdkerrors.Register(ModuleName, 111, "datahash does not exist")
	AccuseAccountInvalid         = sdkerrors.Register(ModuleName, 112, "complain account is invalid")
	ComplainIdInvalid            = sdkerrors.Register(ModuleName, 113, "param complainId is empty")
	ErrComplainStatusInvalid     = sdkerrors.Register(ModuleName, 114, "current complain status invalid")
	ErrResponseStatusInvalid     = sdkerrors.Register(ModuleName, 115, "current complain response status invalid")
	ErrComplainDoesNotExist      = sdkerrors.Register(ModuleName, 116, "complain not exist")
	ErrComplainFinished          = sdkerrors.Register(ModuleName, 126, "current complain has finished")
	ErrAccountHasNoVoteRight     = sdkerrors.Register(ModuleName, 127, "current account has no vote right")
	ErrInvalidContributionValue  = sdkerrors.Register(ModuleName, 128, "Invalid download contribution value")
	ErrComplainVoteStatusInvalid = sdkerrors.Register(ModuleName, 129, "param voteStatus is invalid")
	HasVoted                     = sdkerrors.Register(ModuleName, 130, "current account has vote for this complain")
	AccountSpaceReturnError      = sdkerrors.Register(ModuleName, 134, "Storage space return failed")
	AccountSpaceNotEnough        = sdkerrors.Register(ModuleName, 135, "miner space not enough")
	NoSpaceMinerRewardErr        = sdkerrors.Register(ModuleName, 136, "There is no mining income to be claimed")
	OnlyMainTokenErr             = sdkerrors.Register(ModuleName, 137, "Please use the main currency")
	TokenIdEmpty                 = sdkerrors.Register(ModuleName, 138, "token id not empty")
	TokenidHasExist              = sdkerrors.Register(ModuleName, 142, "tokenId has exist")
	SpaceSettlementErr           = sdkerrors.Register(ModuleName, 143, "Reward space settlement failed")
	CopyrightMortgageErr         = sdkerrors.Register(ModuleName, 144, "The copyright has not expired, and the purchase is invalid")
	UnAuthorizedAccountError     = sdkerrors.Register(ModuleName, 145, "Unauthorized account number")
	CopyrightAccountError        = sdkerrors.Register(ModuleName, 146, "Inconsistent ownership of copyright")
	SpaceFirstError              = sdkerrors.Register(ModuleName, 147, "Please redeem the space first")
	CopyrightComplainErr         = sdkerrors.Register(ModuleName, 148, "current datahash is appealing")
	CopyrightVoteInvalidErr      = sdkerrors.Register(ModuleName, 149, "Copyright voting is invalid")
	CopyrightVoteNotEnoughErr    = sdkerrors.Register(ModuleName, 150, "copyright vote power not enough")
	CopyrightApproveHasFinished  = sdkerrors.Register(ModuleName, 151, "Copyright audit voting has ended")
	FeeInvalidErr                = sdkerrors.Register(ModuleName, 152, "fee amount is invalid")
	QuantityIllegalErr           = sdkerrors.Register(ModuleName, 153, "amount quantity illegal")
	CopyrightPirceErr            = sdkerrors.Register(ModuleName, 154, "copyright price illegal")
	CopyrightChargeRateErr       = sdkerrors.Register(ModuleName, 155, "copyright charge rate illegal")
	AmountDecimalMax18Err        = sdkerrors.Register(ModuleName, 156, "decimal point cannot exceed 18 digits")
)
