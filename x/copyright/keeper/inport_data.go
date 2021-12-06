package keeper

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/export"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankType "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/shopspring/decimal"
	"strconv"
)

func (k Keeper) ImportData(ctx sdk.Context) {
	k.ImportAccountBalanceData(ctx)

	deflationSpaceTotal := k.ImportdeflationInforData(ctx)
	k.ImportAccountInviteData(ctx)
	k.ImportAccountSpaceData(ctx, deflationSpaceTotal)
	k.ImportInviteRecordData(ctx)
}

func (k Keeper) ImportInviteRecordData(ctx sdk.Context) {
	var inviteRecords []types.InviteRecords
	err := util.Json.Unmarshal([]byte(export.InviteRecordTotal), &inviteRecords)
	if err != nil {
		panic("unmarshal error" + err.Error())
	}
	store := k.KVHelper(ctx)
	for i := 0; i < len(inviteRecords); i++ {
		accountInviteRecords := inviteRecords[i]
		recordLen := len(accountInviteRecords.AccountInviteRecord)
		recordArray := make([]types.InviteRecording, 0)
		for j := 0; j < recordLen; j++ {
			record := accountInviteRecords.AccountInviteRecord[j]
			inviteTime, err := strconv.ParseInt(record.InviteTime, 10, 64)
			if err != nil {
				panic(err)
			}
			inviteRecord := types.InviteRecording{
				Address:    record.Address,
				InviteTime: inviteTime,
			}
			recordArray = append(recordArray, inviteRecord)
		}
		err = store.Set(InviteRecordingKey+accountInviteRecords.Account, recordArray)
		if err != nil {
			panic("save account error" + err.Error())
		}
	}
}

func (k Keeper) ImportdeflationInforData(ctx sdk.Context) decimal.Decimal {
	var deflationInfor types.DeflationInfor
	err := util.Json.Unmarshal([]byte(export.DeflationInfor), &deflationInfor)
	if err != nil {
		panic("unmarshal error" + err.Error())
	}
	store := k.KVHelper(ctx)

	newDeflationSpaceTotal := decimal.RequireFromString(deflationInfor.DeflationSpaceTotal).Mul(ByteToMb)
	store.Set(deflationSpaceTotalKey, newDeflationSpaceTotal.StringFixed(4))

	newSpaceMinerAmount := decimal.RequireFromString(deflationInfor.SpaceMinerAmount)
	store.Set(spaceMinerAmountKey, newSpaceMinerAmount)
	k.SetSpaceMinerBonusAmount(ctx, config.ChuangshiFee)
	return newDeflationSpaceTotal
}

func (k Keeper) ImportAccountBalanceData(ctx sdk.Context) {
	var accountBalance []bankType.Balance
	err := util.Json.Unmarshal([]byte(export.AccountBalance), &accountBalance)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(accountBalance); i++ {
		accountBal := accountBalance[i]
		accountAddress, err := sdk.AccAddressFromBech32(accountBal.Address)
		if err != nil {
			if accountBal.Address == "fsv000000000000000000000000000000000000000" {
				accountAddress = sdk.AccAddress([]byte(sdk.BlackHoleAddress))
			} else {
				panic("account error" + accountBal.Address + err.Error())
			}

		}
		err = k.CoinKeeper.SetBalances(ctx, accountAddress, accountBal.Coins)
		if err != nil {
			panic("account error" + err.Error())
		}
	}
}

func (k Keeper) ImportAccountInviteData(ctx sdk.Context) {
	var accountInvite []types.AccountInvite
	err := util.Json.Unmarshal([]byte(export.AccountInvite), &accountInvite)
	if err != nil {
		panic(err)
	}
	store := k.KVHelper(ctx)
	for i := 0; i < len(accountInvite); i++ {
		accountInv := accountInvite[i]
		store.Set(InviteRelationKey+accountInv.Account, accountInv.SupAccount)
	}
}

func (k Keeper) ImportAccountSpaceData(ctx sdk.Context, spaceTotal decimal.Decimal) {
	var accountSpaces []types.AccountSpace
	err := util.Json.Unmarshal([]byte(export.AccountSpace), &accountSpaces)
	if err != nil {
		panic(err)
	}
	store := k.KVHelper(ctx)

	for i := 0; i < len(accountSpaces); i++ {
		logs.Info("account infor", accountSpaces[i])
		accountSpace := accountSpaces[i]
		accountSpaceMiner := AccountSpaceMiner{}

		accountSpaceKey := spaceMinerKey + accountSpace.Account
		spaceTotal := decimal.RequireFromString(accountSpace.SpaceTotal).Mul(ByteToMb)
		if store.Has(accountSpaceKey) {
			err = store.GetUnmarshal(accountSpaceKey, &accountSpaceMiner)
			if err != nil {
				panic(err)
			}
			accountSpaceMiner.SpaceTotal = accountSpaceMiner.SpaceTotal.Add(spaceTotal)
			accountSpaceMiner.BuySpace = accountSpaceMiner.BuySpace.Add(spaceTotal)
		} else {
			accountSpaceMiner.Account = accountSpace.Account
			accountSpaceMiner.SpaceTotal = spaceTotal
			accountSpaceMiner.BuySpace = spaceTotal
			settlement := make(map[int64]Settlement)
			set := Settlement{
				Index:      1,
				IndexSpace: accountSpaceMiner.SpaceTotal,
			}
			settlement[1] = set
			accountSpaceMiner.Settlement = settlement
		}
		err = store.Set(spaceMinerKey+accountSpace.Account, accountSpaceMiner)
		if err != nil {
			panic("import account infor" + err.Error())
		}
		accAddress, err := sdk.AccAddressFromBech32(accountSpace.Account)
		if err != nil {
			logs.Error("import account infor err")
			panic(err)
		}
		k.InviteReward(ctx, accountSpaceMiner.SpaceTotal, accAddress, 1)
	}
}
