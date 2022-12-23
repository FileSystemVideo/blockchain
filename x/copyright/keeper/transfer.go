package keeper

import (
	"fs.video/blockchain/core"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sirupsen/logrus"
)


func (k Keeper) culFee(coins sdk.Coins) (sdk.Coins, error) {
	isFsv := false
	isTip := false
	for _, coin := range coins {
		if coin.Denom == core.MainToken {
			isFsv = true
		}
		if coin.Denom == core.InviteToken {
			isTip = true
		}
	}
	fee := types.NewLedgerCoins(core.CopyrightInviteFee)
	for _, coin := range coins {
		if coin.Denom == core.MainToken {
			
			minTransfer := types.NewLedgerCoin(core.MinFsvTransfer)
			if coin.IsLT(minTransfer) && !coin.IsEqual(fee[0]) { //FSVTIP，
				return fee, types.FeeInvalidErr
			}
			if !isTip { //TIP,FSV
				if coin.IsLT(minTransfer) {
					return fee, types.FeeInvalidErr
				}
			}
			if coin.IsGTE(minTransfer) { //FSV
				
				rate, err := sdk.NewDecFromStr(core.MortgageFee.String())
				if err != nil {
					return fee, err
				}
				culFee, _ := sdk.NewDecCoinsFromCoins(coin).MulDec(rate).TruncateDecimal()
				if isTip { //TIP ，2
					if culFee[0].IsGTE(fee[0]) {
						fee = culFee
					}
				} else {
					fee = culFee
				}
			}
		}
	}
	//fsv
	if !isFsv {
		return fee, types.FeeInvalidErr
	}
	return fee, nil
}


func (k Keeper) Transfer(ctx sdk.Context, fromAddress, toAddress sdk.AccAddress, coins sdk.Coins) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)

	coins = coins.Sort()
	if fromAddress.String() == core.CrossChainAccount && toAddress.String() == core.CrossChainAutoDump {

	} else if fromAddress.String() == core.CrossChainAutoDump && toAddress.String() == core.CrossChainAccount {

	} else {
		
		fee, err := k.culFee(coins)
		if err != nil {
			log.WithError(err).Error("culFee")
			return err
		}
		
		coins = coins.Sub(fee)
		err = k.CoinKeeper.SendCoinsFromAccountToModule(ctx, fromAddress, authtypes.FeeCollectorName, fee)
		if err != nil {
			log.WithError(err).Error("SendCoinsFromAccountToModule")
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
		log.WithError(err).Error("SendCoins")
		return err
	}
	
	for i := 0; i < coins.Len(); i++ {
		coin := coins[i]
		if coin.Denom == core.InviteToken {
			flag := k.IsInvite(ctx, toAddress)
			if !flag {
				
				err = k.InviteRelation(ctx, fromAddress, toAddress)
				if err != nil {
					log.WithError(err).WithFields(logrus.Fields{
						"fromAddress": fromAddress,
						"toAddress":   toAddress,
					}).Error("InviteRelation")
					return err
				}
				
				err = k.InviteRecording(ctx, fromAddress, toAddress, ctx.BlockTime().Unix())
				if err != nil {
					log.WithError(err).WithFields(logrus.Fields{
						"fromAddress": fromAddress,
						"toAddress":   toAddress,
					}).Error("InviteRecording")
					return err
				}
				err = k.transferInviteReward(ctx, fromAddress, toAddress)
				if err != nil {
					log.WithError(err).WithFields(logrus.Fields{
						"fromAddress": fromAddress,
						"toAddress":   toAddress,
					}).Error("transferInviteReward")
					return err
				}
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
