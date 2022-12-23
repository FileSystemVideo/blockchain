package ante

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
)

// MempoolFeeDecorator will check if the transaction's fee is at least as large
// as the local validator's minimum gasFee (defined in validator config).
// If fee is too low, decorator returns error and tx is rejected from mempool.
// Note this only applies when ctx.CheckTx = true
// If fee is high enough or not CheckTx, then call next AnteHandler
// CONTRACT: Tx must implement FeeTx to use MempoolFeeDecorator
type MempoolFeeDecorator struct{}

func NewMempoolFeeDecorator() MempoolFeeDecorator {
	return MempoolFeeDecorator{}
}

func (mfd MempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()

	
	if len(feeCoins) > 1 { 
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "There can only be one handling fee currency")
	}

	for _, coin := range feeCoins {
		if coin.GetDenom() != sdk.DefaultBondDenom { //fsv
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Fee coin only supports fsv")
		}
	}

	gas := feeTx.GetGas()

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() && !simulate {
		minGasPrices := ctx.MinGasPrices()
		msg := feeTx.GetMsgs()
		for _, m := range msg {
			for _, acc := range config.BlackList {
				for _, signer := range m.GetSigners() {
					if signer.String() == acc {
						return ctx, errors.New("Illegal account")
					}
				}
			}
			for _, acc := range config.WhiteList {
				for _, signer := range m.GetSigners() {
					if signer.String() == acc {
						switch proto.MessageName(m) {
						case proto.MessageName((*stakingTypes.MsgDelegate)(nil)):
							return next(ctx, tx, simulate)
						case proto.MessageName((*stakingTypes.MsgUndelegate)(nil)):
							return next(ctx, tx, simulate)
						case proto.MessageName((*stakingTypes.MsgBeginRedelegate)(nil)):
							return next(ctx, tx, simulate)
						case proto.MessageName((*stakingTypes.MsgCreateValidator)(nil)):
							return next(ctx, tx, simulate)
						case proto.MessageName((*stakingTypes.MsgEditValidator)(nil)):
							return next(ctx, tx, simulate)
						case proto.MessageName((*distributionTypes.MsgWithdrawDelegatorReward)(nil)):
							return next(ctx, tx, simulate)
						case proto.MessageName((*distributionTypes.MsgWithdrawValidatorCommission)(nil)):
							return next(ctx, tx, simulate)
						case proto.MessageName((*govTypes.MsgSubmitProposal)(nil)):
							return next(ctx, tx, simulate)
						case proto.MessageName((*govTypes.MsgDeposit)(nil)):
							return next(ctx, tx, simulate)
						case proto.MessageName((*govTypes.MsgVote)(nil)):
							return next(ctx, tx, simulate)
						default:
							return ctx, errors.New("Illegal account")
						}
					}
				}
			}
		}
		if len(msg) == 1 {
			if proto.MessageName(msg[0]) == "copyrightBonusRear" {
				return next(ctx, tx, simulate)
			}
			if proto.MessageName(msg[0]) == "copyrightBonusRearV2" {
				return next(ctx, tx, simulate)
			}
			if proto.MessageName(msg[0]) == "crossChainIn" {
				return next(ctx, tx, simulate)
			}
			if proto.MessageName(msg[0]) == "crossChainOut" {
				return next(ctx, tx, simulate)
			}
			if proto.MessageName(msg[0]) == "copyright.v1beta1.MsgTransfer" {
				return next(ctx, tx, simulate)
			}
		}
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdk.NewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	return next(ctx, tx, simulate)
}

// DeductFeeDecorator deducts fees from the first signer of the tx
// If the first signer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type DeductFeeDecorator struct {
	ak             AccountKeeper
	bankKeeper     types.BankKeeper
	feegrantKeeper FeegrantKeeper
}

func NewDeductFeeDecorator(ak AccountKeeper, bk types.BankKeeper, fk FeegrantKeeper) DeductFeeDecorator {
	return DeductFeeDecorator{
		ak:             ak,
		bankKeeper:     bk,
		feegrantKeeper: fk,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := dfd.ak.GetModuleAddress(types.FeeCollectorName); addr == nil {
		return ctx, fmt.Errorf("Fee collector module account (%s) has not been set", types.FeeCollectorName)
	}

	fee := feeTx.GetFee()
	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()

	deductFeesFrom := feePayer

	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if dfd.feegrantKeeper == nil {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, tx.GetMsgs())

			if err != nil {
				return ctx, sdkerrors.Wrapf(err, "%s not allowed to pay fees from %s", feeGranter, feePayer)
			}
		}

		deductFeesFrom = feeGranter
	}
	//ctx.Logger().Error("",feeTx.GetFee().String())
	deductFeesFromAcc := dfd.ak.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		deductFeesFromAcc = types.NewBaseAccountWithAddress(feePayer)
		//return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", deductFeesFrom)
	}
	// deduct the fees
	if !feeTx.GetFee().IsZero() {
		msg := feeTx.GetMsgs()
		if len(msg) == 1 {
			if proto.MessageName(msg[0]) == "copyright.v1beta1.MsgTransfer" {
				return next(ctx, tx, simulate)
			}
		}
		err = DeductFees(dfd.bankKeeper, ctx, deductFeesFromAcc, feeTx.GetFee())
		if err != nil {
			return ctx, err
		}
	}

	events := sdk.Events{sdk.NewEvent(sdk.EventTypeTx,
		sdk.NewAttribute(sdk.AttributeKeyFee, feeTx.GetFee().String()),
	)}
	ctx.EventManager().EmitEvents(events)

	return next(ctx, tx, simulate)
}

// DeductFees deducts fees from the given account.
func DeductFees(bankKeeper types.BankKeeper, ctx sdk.Context, acc types.AccountI, fees sdk.Coins) error {
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, err.Error())
	}

	return nil
}
