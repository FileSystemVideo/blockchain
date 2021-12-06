package keeper

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	slashingKeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/tendermint/tendermint/libs/log"
	"strconv"
	"strings"
)

type (
	Keeper struct {
		distributionKeeper distKeeper.Keeper
		stakingKeeper      stakingKeeper.Keeper
		slashingKeeper     slashingKeeper.Keeper
		CoinKeeper         bankkeeper.Keeper
		cdc                codec.Marshaler
		storeKey           sdk.StoreKey
		memKey             sdk.StoreKey
	}
)

func NewKeeper(cdc codec.Marshaler, storeKey, memKey sdk.StoreKey, coinKeeper bankkeeper.Keeper, stakingKeeper stakingKeeper.Keeper, distributionKeeper distKeeper.Keeper, slashingKeeper slashingKeeper.Keeper) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		CoinKeeper:         coinKeeper,
		stakingKeeper:      stakingKeeper,
		distributionKeeper: distributionKeeper,
		slashingKeeper:     slashingKeeper,
	}
}

func (k Keeper) KVHelper(ctx sdk.Context) storeHelper {
	store := ctx.KVStore(k.storeKey)
	return storeHelper{
		store,
	}
}

func (k Keeper) PutBlockRelationMsg(ctx sdk.Context, copyright types.CopyrightData) error {
	store := k.KVHelper(ctx)
	return store.Set(types.CopyrightDetailKey+copyright.DataHash, copyright)
}

func GetTxHashAccount(txHash string) (string, error) {
	sum := sha256.Sum256([]byte(txHash))
	return bech32.ConvertAndEncode(sdk.Bech32MainPrefix, sum[:20])
}

func (k Keeper) GetCopyright(ctx sdk.Context, hash string) (data []byte, err error) {
	store := k.KVHelper(ctx)
	storeKey := types.CopyrightDetailKey + hash
	if store.Has(storeKey) {
		bz := store.Get(storeKey)
		return bz, err
	} else {
		return nil, types.CopyrightNotFoundErr
	}
	return
}

func (k Keeper) UpdateCopyrightCreator(ctx sdk.Context, copyright types.CopyrightData) error {
	store := k.KVHelper(ctx)
	partyKey := types.CopyrightPartyKey + copyright.Creator.String()
	if !store.Has(partyKey) {
		return types.CopyrightPartyNotFoundErr
	}
	storeKey := types.CopyrightDetailKey + copyright.DataHash
	return store.Set(storeKey, copyright)
}

func (k Keeper) UpdateCopyrightStatus(ctx sdk.Context, datahash string, status int) error {
	store := k.KVHelper(ctx)
	var copyrightData types.CopyrightData
	storeKey := types.CopyrightDetailKey + datahash
	copyrightBytes := store.Get(storeKey)
	if copyrightBytes == nil {
		return types.CopyrightNotFoundErr
	}
	err := util.Json.Unmarshal(copyrightBytes, &copyrightData)
	if err != nil {
		return err
	}
	copyrightData.ApproveStatus = status
	copyrightBytes, err = util.Json.Marshal(copyrightData)
	store.Set(storeKey, copyrightBytes)
	return k.UpdateCopyrightOriginDataHash(ctx, copyrightData.OriginDataHash, status)
}

func (k Keeper) SetCopyright(ctx sdk.Context, copyright types.CopyrightData) error {
	store := k.KVHelper(ctx)

	partyKey := types.CopyrightPartyKey + copyright.Creator.String()
	if !store.Has(partyKey) {
		return types.CopyrightPartyNotFoundErr
	}
	storeKey := types.CopyrightDetailKey + copyright.DataHash
	err := store.Set(storeKey, copyright)
	if err != nil {
		return err
	}
	err = k.SetCopyrightIp(ctx, copyright.DataHash, copyright.Ip, copyright.Creator)
	if err != nil {
		return err
	}
	k.SetCopyrightOriginDataHash(ctx, copyright.OriginDataHash, copyright.Creator)
	publishTime := ctx.BlockTime()
	k.UploadSpace(ctx, fmt.Sprintf("%d-%d-%d", publishTime.Year(), publishTime.Month(), publishTime.Day()))
	flag := k.UpdateAccountSpaceUsed(ctx, copyright.Creator, copyright.Size)
	if !flag {
		return types.AccountSpaceNotEnough
	}
	err = k.SetCopyrightVoteForResult(ctx, copyright.DataHash)
	if err != nil {
		return err
	}
	return k.BuildNft(ctx, copyright.DataHash, copyright.Creator.String(), copyright.Name, copyright.Name)

}

func (k Keeper) EditorCopyright(ctx sdk.Context, copyright types.EditorCopyrightData) error {
	store := k.KVHelper(ctx)
	storeKey := types.CopyrightDetailKey + copyright.DataHash
	var copyrightData types.CopyrightData
	err := store.GetUnmarshal(storeKey, &copyrightData)
	if err != nil {
		return err
	}
	if !copyrightData.Creator.Equals(copyright.Creator) {
		return types.CopyrightAccountError
	}
	copyrightData.Ip = copyright.Ip
	copyrightData.Name = copyright.Name
	copyrightData.Price = copyright.Price
	copyrightData.ChargeRate = copyright.ChargeRate
	err = store.Set(storeKey, copyrightData)
	if err != nil {
		return err
	}
	err = k.SetCopyrightIp(ctx, copyright.DataHash, copyright.Ip, copyright.Creator)
	if err != nil {
		return err
	}
	k.EditNft(ctx, copyright.DataHash, copyright.Name, copyright.Name)
	return nil
}

func (k Keeper) DeleteCopyright(ctx sdk.Context, copyright types.DeleteCopyrightData) error {
	store := k.KVHelper(ctx)
	var copyrightData types.CopyrightData

	err := store.GetUnmarshal(types.CopyrightDetailKey+copyright.DataHash, &copyrightData)
	if err != nil {
		return err
	}

	if !copyrightData.Creator.Equals(copyright.Creator) {
		return types.CopyrightAccountError
	}

	store.Delete(types.CopyrightDetailKey + copyright.DataHash)
	store.Delete(types.CopyrightOriginHashKey + copyrightData.OriginDataHash)

	if copyrightData.ApproveStatus != 2 {
		flag := k.UpdateAccountSpace(ctx, copyright.Creator, copyrightData.Size)
		if !flag {
			return types.AccountSpaceReturnError
		}
	}
	k.DeleteNft(ctx, copyright.DataHash)
	/*pubAccountString, err := GetTxHashAccount(copyright.DataHash)
	if err != nil {
		return err
	}
	pubAccount, err := sdk.AccAddressFromBech32(pubAccountString)
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, copyrightData.PublishPrice.TruncateInt()))
	err = k.CoinKeeper.SendCoins(ctx, pubAccount, copyright.Creator, coins)
	if err != nil {
		return err
	}*/
	return nil
}

func (k Keeper) DeleteCopyrightInfor(ctx sdk.Context, datahash string) error {
	store := k.KVHelper(ctx)
	storeKey := types.CopyrightDetailKey + datahash
	copyrightBytes := store.Get(storeKey)
	var copyrightData types.CopyrightData
	err := util.Json.Unmarshal(copyrightBytes, &copyrightData)
	if err != nil {
		return err
	}
	store.Delete(storeKey)
	originHashKey := types.CopyrightOriginHashKey + copyrightData.OriginDataHash
	store.Delete(originHashKey)

	k.DeleteNft(ctx, datahash)
	return nil
}

func (k Keeper) CopyrightBonus(ctx sdk.Context, copyright types.CopyrightBonusData) error {
	txhash := formatTxHash(ctx.TxBytes())
	copyrightByte, err := k.GetCopyright(ctx, copyright.DataHash)
	if err != nil {
		return err
	}
	var copyrightData types.CopyrightData
	err = util.Json.Unmarshal(copyrightByte, &copyrightData)
	if err != nil {
		return err
	}
	blockHeight := ctx.BlockHeight()
	copyrightBonusByte := k.GetCopyrightBonusInfo(ctx, copyright.DataHash, copyright.Downer.String())
	resourceRelation := CopyrightExtrainfor{}
	if copyrightBonusByte != nil {
		err := util.Json.Unmarshal(copyrightBonusByte, &resourceRelation)
		if err != nil {
			return errors.New("format mortgage error")
		}
		//if resourceRelation.Species == "buy" || resourceRelation.Height > blockHeight {
		if resourceRelation.Height > blockHeight {
			return types.CopyrightMortgageErr
		}
	}

	fee := "0"
	if copyright.Fee.Amount != nil && copyright.Fee.Amount.IsValid() {
		logs.Info("current fee", copyright.Fee.Amount[0].Amount.String())
		fee = types.MustParseLedgerInt(copyright.Fee.Amount[0].Amount)
	}
	k.SetCopyrightDownRelation(ctx, copyright.DataHash, copyrightData.Price, copyright.Downer, "buy", ctx.BlockHeight()+config.MortgageHeight)
	ledgeCoin := types.MustRealCoin2LedgerCoin(copyrightData.Price)
	bonusCoins := sdk.NewCoins(ledgeCoin)
	//err = k.CoinKeeper.SendCoins(ctx, copyright.Downer, copyright.HashAccount, bonusCoins)
	err = k.CoinKeeper.SendCoinsFromAccountToModule(ctx, copyright.Downer, config.KeyCopyrightBonus, bonusCoins)
	if err != nil {
		return err
	}
	if copyright.BonusType == config.BonusTypeFront {
		dealBonusAuthrizeLogic(ctx, k, copyright.OfferAccountShare, txhash, fee, config.KeyCopyrightBonus, ctx.BlockHeight(), copyrightData)
	}

	feeCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, fee)
	bonusTradeInfor := types.TradeInfor{
		From:           copyright.Downer.String(),
		To:             config.ContractAddressBonus.String(),
		Txhash:         txhash,
		TradeType:      types.TradeTypeCopyrightBuy,
		Amount:         copyrightData.Price,
		Fee:            feeCoin,
		BlockNum:       ctx.BlockHeight(),
		TimeStemp:      ctx.BlockTime().Unix(),
		FromFsvBalance: k.GetBalance(ctx, config.MainToken, copyright.Downer),
		FromTipBalance: k.GetBalance(ctx, config.InviteToken, copyright.Downer),
		ToFsvBalance:   k.GetBalance(ctx, config.MainToken, config.ContractAddressBonus),
		ToTipBalance:   k.GetBalance(ctx, config.InviteToken, config.ContractAddressBonus),
	}

	buySourceInfor := types.BuySourceInfor{
		Txhash:    txhash,
		DataHash:  copyrightData.DataHash,
		Creator:   copyrightData.Creator.String(),
		Purchaser: copyright.Downer.String(),
		Price:     copyrightData.Price,
		Fee:       feeCoin,
		Remark:    "buy",
		TimeStemp: ctx.BlockTime().Unix(),
	}
	err = k.AddBlockRDS(ctx, types.NewBlockRD(bonusTradeInfor))
	if err != nil {
		logs.Error("save bonus infor error", err)
		return err
	}
	err = k.AddBlockRDS(ctx, types.NewBlockRD(buySourceInfor))
	if err != nil {
		logs.Error("save bonus infor error", err)
		return err
	}
	return nil
}

type CopyrightIpinfor struct {
	Creator sdk.AccAddress `json:"creator"`
	Ip      string         `json:"ip"`
}

type CopyrightOrigininfor struct {
	Creator        sdk.AccAddress `json:"creator"`
	OriginDataHash string         `json:"origin_data_hash"`
	Status         int            `json:"status"`
}

type CopyrightCountInfor struct {
	DateStr string `json:"date_str"`
	Count   int64  `json:"count"`
}

type CopyrightExtrainfor struct {
	Downer  sdk.AccAddress `json:"downer"`
	Price   types.RealCoin `json:"price"`
	Species string         `json:"species"`
	Height  int64          `json:"height"`
}

func (k Keeper) SetCopyrightIp(ctx sdk.Context, datahash, ip string, createAccount sdk.AccAddress) error {
	store := k.KVHelper(ctx)
	var copyrightIp CopyrightIpinfor
	copyrightIp.Ip = ip
	copyrightIp.Creator = createAccount
	return store.Set(types.CopyrightIpKey+datahash, copyrightIp)
}

func (k Keeper) SetCopyrightOriginDataHash(ctx sdk.Context, datahash string, createAccount sdk.AccAddress) error {
	store := k.KVHelper(ctx) //ctx.KVStore(k.storeKey)
	var copyrightOrigin CopyrightOrigininfor
	copyrightOrigin.OriginDataHash = datahash
	copyrightOrigin.Creator = createAccount
	copyrightOrigin.Status = 0
	return store.Set(types.CopyrightOriginHashKey+datahash, copyrightOrigin)
}

func (k Keeper) UpdateCopyrightOriginDataHash(ctx sdk.Context, datahash string, status int) error {
	//store := ctx.KVStore(k.storeKey)
	store := k.KVHelper(ctx)
	key := types.CopyrightOriginHashKey + datahash
	var copyrightOrigin CopyrightOrigininfor
	if store.Has(key) {
		err := store.GetUnmarshal(key, &copyrightOrigin)
		if err != nil {
			return err
		}
		copyrightOrigin.Status = status
		err = store.Set(key, copyrightOrigin)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) QueryCopyrightOriginDataHash(ctx sdk.Context, datahash string) ([]byte, error) {
	store := k.KVHelper(ctx)
	key := types.CopyrightOriginHashKey + strings.ToLower(datahash)
	copyrightOriginByte := store.Get(key)
	return copyrightOriginByte, nil
}

func (k Keeper) SetCopyrightAuthor(ctx sdk.Context, datahash string, account sdk.AccAddress) error {

	store := k.KVHelper(ctx) //ctx.KVStore(k.storeKey)
	key := types.CopyrightDetailKey + datahash
	var copyrightInfor types.CopyrightData
	err := store.GetUnmarshal(key, &copyrightInfor)
	if err != nil {
		return err
	}
	copyrightInfor.Creator = account
	return store.Set(key, copyrightInfor)
}


func (k Keeper) UploadSpace(ctx sdk.Context, dateStr string) bool {
	var copyrightCountInfor CopyrightCountInfor
	store := k.KVHelper(ctx)
	key := types.CopyrightCountKey + dateStr
	if !k.IsDataCountPresent(ctx, dateStr) {
		copyrightCountInfor.DateStr = dateStr
		copyrightCountInfor.Count = 1
		err := store.Set(key, copyrightCountInfor)
		if err != nil {
			return false
		}
		return true
	}

	err := store.GetUnmarshal(key, &copyrightCountInfor)
	if err != nil {
		return false
	}
	copyrightCountInfor.Count += 1

	err = store.Set(key, copyrightCountInfor)
	if err != nil {
		return false
	}
	return true
}

func (k Keeper) IsDataCountPresent(ctx sdk.Context, dateStr string) bool {
	store := k.KVHelper(ctx)
	key := types.CopyrightCountKey + dateStr
	return store.Has(key)
}


func (k Keeper) GetPubCount(ctx sdk.Context, dateStr string) CopyrightCountInfor {
	var copyrightCountInfor CopyrightCountInfor
	store := k.KVHelper(ctx)
	key := types.CopyrightCountKey + dateStr
	if !k.IsDataCountPresent(ctx, dateStr) {
		copyrightCountInfor.DateStr = dateStr
		copyrightCountInfor.Count = 0
		store.Set(key, copyrightCountInfor)
		return copyrightCountInfor
	}
	store.GetUnmarshal(key, &copyrightCountInfor)
	return copyrightCountInfor
}

func (k Keeper) GetCopyrightBonusInfo(ctx sdk.Context, datahash string, address string) []byte {
	store := k.KVHelper(ctx)
	key := address + "_" + datahash
	bz := store.Get(key)
	return bz
}


func (k Keeper) SetCopyrightDownRelation(ctx sdk.Context, datahash string, price types.RealCoin, downer sdk.AccAddress, species string, heigth int64) {
	downAddress := downer.String()
	store := k.KVHelper(ctx)
	key := downAddress + "_" + datahash
	copyrightExtrainInfor := CopyrightExtrainfor{}
	copyrightExtrainInfor.Price = price
	copyrightExtrainInfor.Downer = downer
	copyrightExtrainInfor.Species = species
	copyrightExtrainInfor.Height = heigth
	store.Set(key, copyrightExtrainInfor)
}

func (k Keeper) GetCopyrightParty(ctx sdk.Context, account string) (data []byte, err error) {
	store := k.KVHelper(ctx)

	storeKey := types.CopyrightPartyKey + account
	if !store.Has(storeKey) {
		return nil, types.CopyrightPartyNotFoundErr
	}
	return store.Get(storeKey), nil
}

func (k Keeper) QueryPublisherIdMap(ctx sdk.Context) []byte {
	store := k.KVHelper(ctx)
	bz := store.Get(types.CopyrightPublishIdKey)
	return bz
}

func (k Keeper) SetPublisherId(ctx sdk.Context, publishIdString string) error {
	store := k.KVHelper(ctx)
	var publisherIdMap map[string]string
	store.GetUnmarshal(types.CopyrightPublishIdKey, &publisherIdMap)
	if publisherIdMap == nil {
		publisherIdMap = make(map[string]string)
	}
	publisherIdMap[publishIdString] = ""
	/*publisherIdMapByte, err := util.Json.Marshal(publisherIdMap)
	if err != nil {
		return err
	}*/
	return store.Set(types.CopyrightPublishIdKey, publisherIdMap)
}



func (k Keeper) SetCopyrightParty(ctx sdk.Context, copyrightParty types.CopyrightPartyData) error {
	store := k.KVHelper(ctx)
	copyrightBytes, err := util.Json.Marshal(copyrightParty)
	if err != nil {
		return err
	}
	accountSpace := k.QueryAccountSpaceMinerInfor(ctx, copyrightParty.Creator.String())
	if accountSpace.SpaceTotal.Sign() <= 0 {
		return types.SpaceFirstError
	}
	storeKey := types.CopyrightPartyKey + copyrightParty.Creator.String()
	store.Set(storeKey, copyrightBytes)
	return k.SetPublisherId(ctx, copyrightParty.Id)
}

func (k Keeper) IsDataHashPresent(ctx sdk.Context, datahash string) bool {
	store := k.KVHelper(ctx)
	key := types.CopyrightDetailKey + datahash
	return store.Has(key)
}

func (k Keeper) GetBlockRDS(ctx sdk.Context, height int64) (data []byte, err error) {
	store := k.KVHelper(ctx)
	heightStr := strconv.FormatInt(height, 10)
	storeKey := types.BlockRelationDataKey + heightStr
	if store.Has(storeKey) {
		bz := store.Get(storeKey)
		return bz, err
	} else {
		return nil, types.RDSNotFoundErr
	}
	return
}

func (k Keeper) AddBlockRDS(ctx sdk.Context, list ...types.BlockRD) error {

	store := k.KVHelper(ctx)
	height := strconv.FormatInt(ctx.BlockHeight(), 10)
	//fmt.Println("AddBlockRelationData()---------------------", height)
	storeKey := types.BlockRelationDataKey + height

	var finalList types.BlockRDS

	if store.Has(storeKey) {
		bLockRdListBytes := store.Get(storeKey)
		err := util.Json.Unmarshal(bLockRdListBytes, &finalList)
		if err != nil {
			return err
		}
		finalList.Add(list...)
	} else {
		finalList.Add(list...)
	}
	return store.Set(storeKey, finalList)
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetBalance(ctx sdk.Context, denom string, address sdk.AccAddress) string {
	balance := k.CoinKeeper.GetBalance(ctx, address, denom)
	return types.MustParseLedgerCoin(balance)
}

func (k Keeper) GetCrossChainOutFeeRatio(ctx sdk.Context) (string, error) {
	store := k.KVHelper(ctx)
	key := types.CopyrightCrossChainOutFeeRatioKey

	if store.Has(key) {
		ratioByte := store.Get(key)
		ratioString := string(ratioByte)
		ratioDec, err := sdk.NewDecFromStr(ratioString)
		if err != nil {
			return "", err
		}
		if ratioDec.LT(sdk.NewDec(0)) {
			return "", errors.New("Fee Set Error")
		}
		return ratioString, nil
	} else {
		return config.CrossChainOutFeeRatio, nil
	}
}


