package keeper

import (
	"fmt"
	"fs.video/blockchain/x/smart/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	stakingKeeper stakingKeeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	evmKeeper     types.EVMKeeper
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	stakingKeeper stakingKeeper.Keeper,
	evmKeeper types.EVMKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramstore:    ps,
		accountKeeper: ak,
		bankKeeper:    bk,
		stakingKeeper: stakingKeeper,
		evmKeeper:     evmKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KVHelper(ctx sdk.Context) StoreHelper {
	store := ctx.KVStore(k.storeKey)
	return StoreHelper{
		store,
	}
}


func (k Keeper) createValidator(ctx sdk.Context, delegatorAddress sdk.AccAddress, validatorAddress sdk.ValAddress, msg types.MsgCreateSmartValidator, delegation sdk.Coin) error {
	pk, err := ParseBech32ValConsPubkey(msg.PubKey)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	if _, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return stakingTypes.ErrValidatorPubKeyExists
	}
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if delegation.Denom != bondDenom {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", delegation.Denom, bondDenom,
		)
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return err
	}

	validator, err := stakingTypes.NewValidator(validatorAddress, pk, msg.Description)
	if err != nil {
		return err
	}
	commission := stakingTypes.NewCommissionWithTime(
		msg.Commission.Rate, msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	)
	validator, err = validator.SetInitialCommission(commission)
	
	validator.MinSelfDelegation = msg.MinSelfDelegation
	k.stakingKeeper.SetValidator(ctx, validator)
	k.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	k.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
	k.stakingKeeper.AfterValidatorCreated(ctx, validator.GetOperator())
	
	_, err = k.stakingKeeper.Delegate(ctx, delegatorAddress, delegation.Amount, stakingTypes.Unbonded, validator, true)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeCreateValidator,
			sdk.NewAttribute(stakingTypes.AttributeKeyValidator, validator.GetOperator().String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, delegation.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, stakingTypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})
	return nil
}


func (k Keeper) QueryContractIsExist(ctx sdk.Context, address common.Address) bool {
	acct := k.evmKeeper.GetAccountWithoutBalance(ctx, address)
	var code []byte
	if acct != nil && acct.IsContract() {
		code = k.evmKeeper.GetCode(ctx, common.BytesToHash(acct.CodeHash))
	}
	if len(code) == 0 {
		return false
	}
	return true
}
