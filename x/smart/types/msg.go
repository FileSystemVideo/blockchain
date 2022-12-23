package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	_ sdk.Msg = &MsgCreateSmartValidator{}
)

const (
	TypeMsgCreateSmartValidator = "create_smart_validator"
)

func NewMsgCreateSmartValidator(
	valAddr sdk.ValAddress, pubKey string, //nolint:interfacer
	selfDelegation sdk.Coin, description stakingtypes.Description, commission stakingtypes.CommissionRates, minSelfDelegation sdk.Int,
) (*MsgCreateSmartValidator, error) {
	return &MsgCreateSmartValidator{
		Description:       description,
		DelegatorAddress:  sdk.AccAddress(valAddr).String(),
		ValidatorAddress:  valAddr.String(),
		PubKey:            pubKey,
		Value:             selfDelegation,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateSmartValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateSmartValidator) Type() string { return TypeMsgCreateSmartValidator }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgCreateSmartValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(delAddr.Bytes(), addr.Bytes()) {
		addrs = append(addrs, sdk.AccAddress(addr))
	}

	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateSmartValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateSmartValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}
	if delAddr.Empty() {
		return stakingtypes.ErrEmptyDelegatorAddr
	}

	if msg.ValidatorAddress == "" {
		return stakingtypes.ErrEmptyValidatorAddr
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	if !sdk.AccAddress(valAddr).Equals(delAddr) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator address is invalid")
	}

	if msg.PubKey == "" {
		return stakingtypes.ErrEmptyValidatorPubKey
	}

	if !msg.Value.IsValid() || !msg.Value.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid delegation amount")
	}

	if msg.Description == (stakingtypes.Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.Commission == (stakingtypes.CommissionRates{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}

	if err := msg.Commission.Validate(); err != nil {
		return err
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return stakingtypes.ErrSelfDelegationBelowMinimum
	}
	return nil
}
