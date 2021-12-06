package keeper

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

func (k Keeper) ExportAccountInvite(ctx sdk.Context) []types.AccountInvite {
	store := ctx.KVStore(k.storeKey)
	inviteRelationKey := strings.ToLower(InviteRelationKey)
	spaceMinerByte := []byte(inviteRelationKey)
	iterator := sdk.KVStorePrefixIterator(store, spaceMinerByte)
	defer iterator.Close()
	accountInviteArray := make([]types.AccountInvite, 0)
	for ; iterator.Valid(); iterator.Next() {
		ctx.Logger().Info("invite infor", string(iterator.Key()), string(iterator.Value()))
		account := string(iterator.Key())
		supAccount := string(iterator.Value())
		account = strings.Replace(account, inviteRelationKey, "", -1)
		accountInvite := types.AccountInvite{
			Account:    account,
			SupAccount: supAccount,
		}
		accountInviteArray = append(accountInviteArray, accountInvite)
	}
	return accountInviteArray
}

func (k Keeper) ExportAccountRecord(ctx sdk.Context) []types.InviteRecords {
	store := ctx.KVStore(k.storeKey)
	inviteRecordingKey := strings.ToLower(InviteRecordingKey)
	spaceMinerByte := []byte(inviteRecordingKey)
	iterator := sdk.KVStorePrefixIterator(store, spaceMinerByte)
	defer iterator.Close()
	accountInviteArray := make([]types.InviteRecords, 0)
	for ; iterator.Valid(); iterator.Next() {
		ctx.Logger().Info("invite infor ", string(iterator.Key()), string(iterator.Value()))
		account := string(iterator.Key())
		account = strings.Replace(account, inviteRecordingKey, "", -1)

		var inviteRecording []types.InviteRecording
		err := util.Json.Unmarshal(iterator.Value(), &inviteRecording)
		if err != nil {
			panic(err)
		}
		accountInviteRecordArray := make([]types.AccountInviteRecord, 0)
		for i := 0; i < len(inviteRecording); i++ {
			inviteTime :=  strconv.FormatInt(inviteRecording[i].InviteTime,10)
			accountInviteRecord := types.AccountInviteRecord{
				Address:    inviteRecording[i].Address,
				InviteTime:inviteTime ,
			}
			accountInviteRecordArray = append(accountInviteRecordArray, accountInviteRecord)
		}
		inviteRecords := types.InviteRecords{}
		inviteRecords.Account = account
		inviteRecords.AccountInviteRecord = accountInviteRecordArray
		accountInviteArray = append(accountInviteArray, inviteRecords)
	}
	return accountInviteArray
}

func (k Keeper) ExportAccountSpace(ctx sdk.Context) []types.AccountSpace {
	store := ctx.KVStore(k.storeKey)
	spaceMinerByte := []byte(spaceMinerKey)
	iterator := sdk.KVStorePrefixIterator(store, spaceMinerByte)
	defer iterator.Close()
	accountSpaceArray := make([]types.AccountSpace, 0)
	for ; iterator.Valid(); iterator.Next() {
		var accountSpaceMiner AccountSpaceMiner

		ctx.Logger().Info("space infor", string(iterator.Key()), string(iterator.Value()))
		key := string(iterator.Key())
		if strings.Contains(key, spaceMinerAccountKey) || strings.Contains(key, spaceMinerAmountKey) || strings.Contains(key, spaceMinerBonusKey) {
			continue
		}
		err := util.Json.Unmarshal(iterator.Value(), &accountSpaceMiner)
		if err != nil {
			var accountMap map[string]AccountSpaceMiner
			err = util.Json.Unmarshal(iterator.Value(), &accountMap)
			if err != nil {
				panic(err)
			}
			continue
		}
		accountSpace := types.AccountSpace{
			Account:    accountSpaceMiner.Account,
			SpaceTotal: accountSpaceMiner.SpaceTotal.String(),
		}
		accountSpaceArray = append(accountSpaceArray, accountSpace)
	}
	return accountSpaceArray
}

func (k Keeper) ExportDeflationMinerInfor(ctx sdk.Context) types.DeflationInfor {
	store := ctx.KVStore(k.storeKey)
	deflationMinerKeyByte := []byte(deflationMinerKey)
	iterator := sdk.KVStorePrefixIterator(store, deflationMinerKeyByte)
	defer iterator.Close()
	deflationInfor := types.DeflationInfor{}
	for ; iterator.Valid(); iterator.Next() {
		var deflationMinerInfor DeflationMinerInfor
		err := util.Json.Unmarshal(iterator.Value(), &deflationMinerInfor)
		if err != nil {
			logs.Error("export data error")
			panic(err)
		}
		deflationInfor.MinerTotalAmount = deflationMinerInfor.MinerTotalAmount.String()
		deflationInfor.HasMinerAmount = deflationMinerInfor.HasMinerAmount.String()
		deflationInfor.RemainMinerAmount = deflationMinerInfor.RemainMinerAmount.String()
		deflationInfor.DayMinerAmount = deflationMinerInfor.DayMinerAmount.String()
		deflationInfor.DayMinerRemain = strconv.FormatInt(deflationMinerInfor.DayMinerRemain, 10)
	}
	spaceMinerAmountKeyByte := []byte(spaceMinerAmountKey)
	iterator = sdk.KVStorePrefixIterator(store, spaceMinerAmountKeyByte)
	for ; iterator.Valid(); iterator.Next() {
		deflationInfor.SpaceMinerAmount = string(iterator.Value())
	}
	spaceMinerBonusKeyByte := []byte(spaceMinerBonusKey)
	iterator = sdk.KVStorePrefixIterator(store, spaceMinerBonusKeyByte)
	for ; iterator.Valid(); iterator.Next() {
		deflationInfor.SpaceMinerBonus = string(iterator.Value())
	}
	deflationSpaceKeyByte := []byte(deflationSpaceTotalKey)
	iterator = sdk.KVStorePrefixIterator(store, deflationSpaceKeyByte)
	for ; iterator.Valid(); iterator.Next() {
		deflationInfor.DeflationSpaceTotal = string(iterator.Value())
	}
	return deflationInfor
}

func (k Keeper) InitDeflationMinerInfor(ctx sdk.Context, deflationInfor types.DeflationInfor) {
	//store := ctx.KVStore(k.storeKey)
	store := k.KVHelper(ctx)
	deflationMinerInfor := DeflationMinerInfor{}
	deflationMinerInfor.MinerTotalAmount = decimal.RequireFromString(deflationInfor.MinerTotalAmount)
	deflationMinerInfor.HasMinerAmount = decimal.RequireFromString(deflationInfor.HasMinerAmount)
	deflationMinerInfor.RemainMinerAmount = decimal.RequireFromString(deflationInfor.RemainMinerAmount)
	deflationMinerInfor.DayMinerAmount = decimal.RequireFromString(deflationInfor.DayMinerAmount)
	deflationMinerInfor.DayMinerRemain, _ = strconv.ParseInt(deflationInfor.DayMinerRemain, 10, 64)
	deflationMinerInfor.MinerBlockNum = 129600
	store.Set(deflationMinerKey, deflationMinerInfor)
	store.Set(spaceMinerAmountKey, deflationInfor.SpaceMinerAmount)
	store.Set(spaceMinerBonusKey, deflationInfor.SpaceMinerBonus)
	store.Set(deflationSpaceTotalKey, deflationInfor.DeflationSpaceTotal)
}

func (k Keeper) InitAccountInvite(ctx sdk.Context, accountInvite types.AccountInvite) {
	//store := ctx.KVStore(k.storeKey)
	store := k.KVHelper(ctx)
	store.Set(InviteRelationKey+accountInvite.Account, accountInvite.SupAccount)
}

func (k Keeper) InitAccountSpace(ctx sdk.Context, accountSpace types.AccountSpace) {
	//store := ctx.KVStore(k.storeKey)
	store := k.KVHelper(ctx)
	accountSpaceMiner := AccountSpaceMiner{
		Account:    accountSpace.Account,
		SpaceTotal: decimal.RequireFromString(accountSpace.SpaceTotal),
	}
	store.Set(spaceMinerKey+accountSpace.Account, accountSpaceMiner)
}
