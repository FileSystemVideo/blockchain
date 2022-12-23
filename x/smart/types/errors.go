package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrAddressBookSet        = sdkerrors.Register(ModuleName, 201, "address set error")
	ErrDelegationCoin        = sdkerrors.Register(ModuleName, 203, "Invalid amount")
	ErrGatewayNumber         = sdkerrors.Register(ModuleName, 204, "Number Already registered")
	ErrGatewayDelegation     = sdkerrors.Register(ModuleName, 205, "Insufficient mortgage amount")
	ErrGatewayNotExist       = sdkerrors.Register(ModuleName, 206, "gateway not exist")
	ErrGatewayNumNotFound    = sdkerrors.Register(ModuleName, 207, "gateway number not found")
	ErrGatewayNumLength      = sdkerrors.Register(ModuleName, 208, "Illegal length of number segment")
	ErrContractNotFound      = sdkerrors.Register(ModuleName, 209, "contract not found")
	ErrEmptyProposalContract = sdkerrors.Register(ModuleName, 210, "invalid proposal contract")
)
