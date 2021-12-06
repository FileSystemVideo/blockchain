package keeper

import (
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
)


func (k Keeper) culFee(coins sdk.Coins) (sdk.Coins, error) {
	isFsv := false
	isTip := false
	for _, coin := range coins {
		if coin.Denom == config.MainToken {
			isFsv = true
		}
		if coin.Denom == config.InviteToken {
			isTip = true
		}
	}
	fee := types.NewLedgerCoins(config.CopyrightInviteFee)
	for _, coin := range coins {
		if coin.Denom == config.MainToken {

			minTransfer := types.NewLedgerCoin(config.MinFsvTransfer)
			if coin.IsLT(minTransfer) && !coin.IsEqual(fee[0]) {
				return fee, types.FeeInvalidErr
			}
			if !isTip {
				if coin.IsLT(minTransfer) {
					return fee, types.FeeInvalidErr
				}
			}
			if coin.IsGTE(minTransfer) {

				rate, err := sdk.NewDecFromStr(config.MortgageFee)
				if err != nil {
					return fee, err
				}
				culFee, _ := sdk.NewDecCoinsFromCoins(coin).MulDec(rate).TruncateDecimal()
				if isTip {
					if culFee[0].IsGTE(fee[0]) {
						fee = culFee
					}
				} else {
					fee = culFee
				}
			}
		}
	}

	if !isFsv {
		return fee, types.FeeInvalidErr
	}
	return fee, nil
}


func (k Keeper) Transfer(ctx sdk.Context, fromAddress, toAddress sdk.AccAddress, coins sdk.Coins) error {

	lockFlag := types2.JudgeLockedAccount(fromAddress.String())
	if lockFlag {
		return sdkerrors.ErrLockedAccount
	}
	coins = coins.Sort()
	if fromAddress.String() == config.CrossChainAccount && toAddress.String() == config.CrossChainAutoDump {

	} else if fromAddress.String() == config.CrossChainAutoDump && toAddress.String() == config.CrossChainAccount {

	} else {

		fee, err := k.culFee(coins)
		if err != nil {
			return err
		}

		coins = coins.Sub(fee)
		err = k.CoinKeeper.SendCoinsFromAccountToModule(ctx, fromAddress, authtypes.FeeCollectorName, fee)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeTransferFee,
				sdk.NewAttribute(types.AttributeKeyTransferFee, fee.String()),
			),
		})
	}

	err := k.CoinKeeper.SendCoins(ctx, fromAddress, toAddress, coins)
	if err != nil {
		return err
	}

	for i := 0; i < coins.Len(); i++ {
		coin := coins[i]
		if coin.Denom == config.InviteToken {
			flag := k.IsInvite(ctx, toAddress)
			if !flag {

				err = k.InviteRelation(ctx, fromAddress, toAddress)
				if err != nil {
					return err
				}

				err = k.InviteRecording(ctx, fromAddress, toAddress, ctx.BlockTime().Unix())
				if err != nil {
					return err
				}
				err = k.transferInviteReward(ctx, fromAddress, toAddress)
				if err != nil {
					return err
				}
			} else {
				logs.Debug("invite infor")
			}

			flag = k.IsInvite(ctx, fromAddress)
			if !flag {

				store := k.KVHelper(ctx)
				key := InviteRelationKey + fromAddress.String()
				store.Set(key, "")
			}

		}
	}
	return nil
}


func (k Keeper) transferInviteReward(ctx sdk.Context, fromAddress, toAddress sdk.AccAddress) error {

	accountMiner := k.QueryAccountSpaceMinerInfor(ctx, toAddress.String())
	err := k.InviteReward(ctx, accountMiner.BuySpace, toAddress, 1)
	if err != nil {
		return err
	}

	toSettlement, err := k.QueryRewardInfo(ctx, toAddress.String())
	if err != nil {
		return err
	}

	if toSettlement != nil {
		sp := toSettlement.InviteRewardSpace.Add(toSettlement.ExpansionRewardSpace)
		err = k.InviteReward(ctx, sp.Round(4), toAddress, 2)
		if err != nil {
			return err
		}
	}
	return nil
}
