package keeper

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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
		authKeeper         types.AccountKeeper
		cdc                codec.Codec
		storeKey           sdk.StoreKey
		paramstore         paramtypes.Subspace
	}
)

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, coinKeeper bankkeeper.Keeper, stakingKeeper stakingKeeper.Keeper, distributionKeeper distKeeper.Keeper, slashingKeeper slashingKeeper.Keeper, authKeeper types.AccountKeeper, ps paramtypes.Subspace) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		CoinKeeper:         coinKeeper,
		stakingKeeper:      stakingKeeper,
		distributionKeeper: distributionKeeper,
		slashingKeeper:     slashingKeeper,
		authKeeper:         authKeeper,
		paramstore:         ps.WithKeyTable(types.ParamKeyTable()),
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

//txhash
func GetTxHashAccount(txHash string) (string, error) {
	sum := sha256.Sum256([]byte(txHash))
	return bech32.ConvertAndEncode(sdk.Bech32MainPrefix, sum[:20])
}


func (k Keeper) GetCopyright(ctx sdk.Context, hash string) (data []byte, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	storeKey := types.CopyrightDetailKey + hash
	if store.Has(storeKey) {
		bz := store.Get(storeKey)
		return bz, err
	} else {
		log.WithError(err).Error("Copyright Not Found")
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
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var copyrightData types.CopyrightData
	storeKey := types.CopyrightDetailKey + datahash
	copyrightBytes := store.Get(storeKey)
	if copyrightBytes == nil {
		return types.CopyrightNotFoundErr
	}
	err := util.Json.Unmarshal(copyrightBytes, &copyrightData)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	copyrightData.ApproveStatus = status
	copyrightBytes, err = util.Json.Marshal(copyrightData)
	err = store.Set(storeKey, copyrightBytes)
	if err != nil {
		log.WithError(err).Error("store.Set")
	}
	return k.UpdateCopyrightOriginDataHash(ctx, copyrightData.OriginDataHash, status)
}


func (k Keeper) SetCopyright(ctx sdk.Context, copyright types.CopyrightData) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
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
	//ip
	err = k.SetCopyrightIp(ctx, copyright.DataHash, copyright.Ip, copyright.Creator)
	if err != nil {
		log.WithError(err).Error("SetCopyrightIp")
		return err
	}
	//hash
	err = k.SetCopyrightOriginDataHash(ctx, copyright.OriginDataHash, copyright.Creator)
	if err != nil {
		log.WithError(err).Error("SetCopyrightOriginDataHash")
		return err
	}
	
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
	//nft
	return k.BuildNft(ctx, copyright.DataHash, copyright.Creator.String(), copyright.Name, copyright.Name)

}


func (k Keeper) EditorCopyright(ctx sdk.Context, copyright types.EditorCopyrightData) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	storeKey := types.CopyrightDetailKey + copyright.DataHash
	var copyrightData types.CopyrightData
	err := store.GetUnmarshal(storeKey, &copyrightData)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
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
		log.WithError(err).Error("store.Set")
		return err
	}
	//ip
	err = k.SetCopyrightIp(ctx, copyright.DataHash, copyright.Ip, copyright.Creator)
	if err != nil {
		log.WithError(err).Error("SetCopyrightIp")
		return err
	}
	//nft
	err = k.EditNft(ctx, copyright.DataHash, copyright.Name, copyright.Name)
	if err != nil {
		log.WithError(err).Error("EditNft")
		return err
	}
	return nil
}


func (k Keeper) DeleteCopyright(ctx sdk.Context, copyright types.DeleteCopyrightData) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var copyrightData types.CopyrightData

	err := store.GetUnmarshal(types.CopyrightDetailKey+copyright.DataHash, &copyrightData)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return err
	}

	
	if !copyrightData.Creator.Equals(copyright.Creator) {
		return types.CopyrightAccountError
	}

	store.Delete(types.CopyrightDetailKey + copyright.DataHash)
	//hash
	store.Delete(types.CopyrightOriginHashKey + copyrightData.OriginDataHash)

	if copyrightData.ApproveStatus != 2 { 
		
		flag := k.UpdateAccountSpace(ctx, copyright.Creator, copyrightData.Size)
		if !flag {
			return types.AccountSpaceReturnError
		}
	}
	//nft
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

// hash
func (k Keeper) DeleteCopyrightInfor(ctx sdk.Context, datahash string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	storeKey := types.CopyrightDetailKey + datahash
	copyrightBytes := store.Get(storeKey)
	var copyrightData types.CopyrightData
	err := util.Json.Unmarshal(copyrightBytes, &copyrightData)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	store.Delete(storeKey)
	//hash
	originHashKey := types.CopyrightOriginHashKey + copyrightData.OriginDataHash
	store.Delete(originHashKey)

	//nft
	k.DeleteNft(ctx, datahash)
	return nil
}


func (k Keeper) CopyrightBonus(ctx sdk.Context, copyright types.CopyrightBonusData) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	txhash := formatTxHash(ctx.TxBytes())
	copyrightByte, err := k.GetCopyright(ctx, copyright.DataHash)
	if err != nil {
		return err
	}
	var copyrightData types.CopyrightData
	err = util.Json.Unmarshal(copyrightByte, &copyrightData)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	blockHeight := ctx.BlockHeight()
	
	copyrightBonusByte := k.GetCopyrightBonusInfo(ctx, copyright.DataHash, copyright.Downer.String())
	resourceRelation := CopyrightExtrainfor{}
	if copyrightBonusByte != nil {
		err := util.Json.Unmarshal(copyrightBonusByte, &resourceRelation)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return errors.New("format mortgage error")
		}
		
		//if resourceRelation.Species == "buy" || resourceRelation.Height > blockHeight {
		if resourceRelation.Height > blockHeight {
			return types.CopyrightMortgageErr
		}
	}

	fee := "0"
	if copyright.Fee.Amount != nil && copyright.Fee.Amount.IsValid() {
		log.WithField("Current handling fee amount", copyright.Fee.Amount[0].Amount.String()).Debug("CopyrightBonus")
		fee = types.MustParseLedgerInt(copyright.Fee.Amount[0].Amount)
	}
	
	k.SetCopyrightDownRelation(ctx, copyright.DataHash, copyrightData.Price, copyright.Downer, "buy", ctx.BlockHeight()+core.MortgageHeight)
	ledgeCoin := types.MustRealCoin2LedgerCoin(copyrightData.Price)
	bonusCoins := sdk.NewCoins(ledgeCoin)
	//err = k.CoinKeeper.SendCoins(ctx, copyright.Downer, copyright.HashAccount, bonusCoins)
	err = k.CoinKeeper.SendCoinsFromAccountToModule(ctx, copyright.Downer, core.KeyCopyrightBonus, bonusCoins)
	if err != nil {
		return err
	}
	if copyright.BonusType == core.BonusTypeFront {
		err = dealBonusAuthrizeLogic(ctx, k, copyright.OfferAccountShare, txhash, fee, core.KeyCopyrightBonus, ctx.BlockHeight(), copyrightData)
		if err != nil {
			return err
		}
	}

	feeCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, fee)
	
	bonusTradeInfor := types.TradeInfor{
		From:           copyright.Downer.String(),
		To:             core.ContractAddressBonus.String(),
		Txhash:         txhash,
		TradeType:      core.TradeTypeCopyrightBuy,
		Amount:         copyrightData.Price,
		Fee:            feeCoin,
		BlockNum:       ctx.BlockHeight(),
		TimeStemp:      ctx.BlockTime().Unix(),
		FromFsvBalance: k.GetBalance(ctx, core.MainToken, copyright.Downer),
		FromTipBalance: k.GetBalance(ctx, core.InviteToken, copyright.Downer),
		ToFsvBalance:   k.GetBalance(ctx, core.MainToken, core.ContractAddressBonus),
		ToTipBalance:   k.GetBalance(ctx, core.InviteToken, core.ContractAddressBonus),
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
		log.WithError(err).Error("AddBlockRDS-bonusTradeInfor")
		return err
	}
	err = k.AddBlockRDS(ctx, types.NewBlockRD(buySourceInfor))
	if err != nil {
		log.WithError(err).Error("AddBlockRDS-buySourceInfor")
		return err
	}
	return nil
}

//Ip
type CopyrightIpinfor struct {
	Creator sdk.AccAddress `json:"creator"`
	Ip      string         `json:"ip"`
}

//hash
type CopyrightOrigininfor struct {
	Creator        sdk.AccAddress `json:"creator"`
	OriginDataHash string         `json:"origin_data_hash"`
	Status         int            `json:"status"` // 0  1  2(,)
}


type CopyrightCountInfor struct {
	DateStr string `json:"date_str"`
	Count   int64  `json:"count"`
}

type CopyrightExtrainfor struct {
	Downer  sdk.AccAddress `json:"downer"`
	Price   types.RealCoin `json:"price"`
	Species string         `json:"species"` //:buy   mortgage 
	Height  int64          `json:"height"`  
}
type CopyrightBonusInfo struct {
	Downer string `json:"downer"`
	Height int64  `json:"height"` 
}

//ip
func (k Keeper) SetCopyrightIp(ctx sdk.Context, datahash, ip string, createAccount sdk.AccAddress) error {
	store := k.KVHelper(ctx)
	var copyrightIp CopyrightIpinfor
	copyrightIp.Ip = ip
	copyrightIp.Creator = createAccount
	return store.Set(types.CopyrightIpKey+datahash, copyrightIp)
}

//hash
func (k Keeper) SetCopyrightOriginDataHash(ctx sdk.Context, datahash string, createAccount sdk.AccAddress) error {
	store := k.KVHelper(ctx) //ctx.KVStore(k.storeKey)
	var copyrightOrigin CopyrightOrigininfor
	copyrightOrigin.OriginDataHash = datahash
	copyrightOrigin.Creator = createAccount
	copyrightOrigin.Status = 0
	return store.Set(types.CopyrightOriginHashKey+datahash, copyrightOrigin)
}

//hash
func (k Keeper) UpdateCopyrightOriginDataHash(ctx sdk.Context, datahash string, status int) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	key := types.CopyrightOriginHashKey + datahash
	var copyrightOrigin CopyrightOrigininfor
	if store.Has(key) {
		err := store.GetUnmarshal(key, &copyrightOrigin)
		if err != nil {
			log.WithError(err).Error("GetUnmarshal")
			return err
		}
		copyrightOrigin.Status = status
		err = store.Set(key, copyrightOrigin)
		if err != nil {
			log.WithError(err).Error("store.Set")
			return err
		}
	}
	return nil
}

//hash
func (k Keeper) QueryCopyrightOriginDataHash(ctx sdk.Context, datahash string) ([]byte, error) {
	store := k.KVHelper(ctx)
	key := types.CopyrightOriginHashKey + strings.ToLower(datahash)
	copyrightOriginByte := store.Get(key)
	return copyrightOriginByte, nil
}

//(nft)
func (k Keeper) SetCopyrightAuthor(ctx sdk.Context, datahash string, account sdk.AccAddress) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx) //ctx.KVStore(k.storeKey)
	key := types.CopyrightDetailKey + datahash
	var copyrightInfor types.CopyrightData
	err := store.GetUnmarshal(key, &copyrightInfor)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return err
	}
	copyrightInfor.Creator = account
	return store.Set(key, copyrightInfor)
}

/**

*/
func (k Keeper) UploadSpace(ctx sdk.Context, dateStr string) bool {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var copyrightCountInfor CopyrightCountInfor
	store := k.KVHelper(ctx)
	key := types.CopyrightCountKey + dateStr
	if !k.IsDataCountPresent(ctx, dateStr) {
		copyrightCountInfor.DateStr = dateStr
		copyrightCountInfor.Count = 1
		err := store.Set(key, copyrightCountInfor)
		if err != nil {
			log.WithError(err).Error("store.Set")
			return false
		}
		return true
	}
	
	err := store.GetUnmarshal(key, &copyrightCountInfor)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return false
	}
	copyrightCountInfor.Count += 1
	
	err = store.Set(key, copyrightCountInfor)
	if err != nil {
		log.WithError(err).Error("store.Set")
		return false
	}
	return true
}

func (k Keeper) IsDataCountPresent(ctx sdk.Context, dateStr string) bool {
	store := k.KVHelper(ctx)
	key := types.CopyrightCountKey + dateStr
	return store.Has(key)
}

/**

*/
func (k Keeper) GetPubCount(ctx sdk.Context, dateStr string) CopyrightCountInfor {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var copyrightCountInfor CopyrightCountInfor
	store := k.KVHelper(ctx)
	key := types.CopyrightCountKey + dateStr
	if !k.IsDataCountPresent(ctx, dateStr) {
		copyrightCountInfor.DateStr = dateStr
		copyrightCountInfor.Count = 0
		err := store.Set(key, copyrightCountInfor)
		if err != nil {
			log.WithError(err).Error("store.Set")
			return copyrightCountInfor
		}
		return copyrightCountInfor
	}
	
	err := store.GetUnmarshal(key, &copyrightCountInfor)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
	}
	return copyrightCountInfor
}


func (k Keeper) GetCopyrightBonusInfo(ctx sdk.Context, datahash string, address string) []byte {
	store := k.KVHelper(ctx)
	key := types.CopyrightRelationKey + address + "_" + datahash
	bz := store.Get(key)
	return bz
}

/*
	
*/
func (k Keeper) SetCopyrightDownRelation(ctx sdk.Context, datahash string, price types.RealCoin, downer sdk.AccAddress, species string, heigth int64) {
	downAddress := downer.String() 
	store := k.KVHelper(ctx)
	key := types.CopyrightRelationKey + downAddress + "_" + datahash
	
	copyrightExtrainInfor := CopyrightExtrainfor{}
	copyrightExtrainInfor.Price = price
	copyrightExtrainInfor.Downer = downer
	copyrightExtrainInfor.Species = species 
	copyrightExtrainInfor.Height = heigth   
	store.Set(key, copyrightExtrainInfor)
}


func (k Keeper) SetBonusAddress(ctx sdk.Context, bonusAddress, downer string, heigth int64) error {
	store := k.KVHelper(ctx)
	copyrightBonusInfo := CopyrightBonusInfo{}
	copyrightBonusInfo.Downer = downer
	copyrightBonusInfo.Height = heigth
	return store.Set(types.CopyrightBonusAddressKey+bonusAddress, copyrightBonusInfo)
}


func (k Keeper) IsExistBonusAddress(ctx sdk.Context, bonusAddress string) bool {
	store := k.KVHelper(ctx)
	return store.Has(types.CopyrightBonusAddressKey + bonusAddress)
}


func (k Keeper) DeleteBonusAddress(ctx sdk.Context, bonusAddress string) {
	store := k.KVHelper(ctx)
	store.Delete(types.CopyrightBonusAddressKey + bonusAddress)
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

//id
func (k Keeper) SetPublisherId(ctx sdk.Context, publishIdString string) error {
	store := k.KVHelper(ctx)
	var publisherIdMap map[string]string
	store.GetUnmarshal(types.CopyrightPublishIdKey, &publisherIdMap)
	if publisherIdMap == nil {
		publisherIdMap = make(map[string]string)
	}
	publisherIdMap[publishIdString] = ""
	return store.Set(types.CopyrightPublishIdKey, publisherIdMap)
}


func (k Keeper) SetCopyrightParty(ctx sdk.Context, copyrightParty types.CopyrightPartyData) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	copyrightBytes, err := util.Json.Marshal(copyrightParty)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return err
	}
	accountSpace := k.QueryAccountSpaceMinerInfor(ctx, copyrightParty.Creator.String())
	//0.,
	if accountSpace.SpaceTotal.Sign() <= 0 {
		return types.SpaceFirstError
	}
	storeKey := types.CopyrightPartyKey + copyrightParty.Creator.String()
	err = store.Set(storeKey, copyrightBytes)
	if err != nil {
		log.WithError(err).Error("store.Set")
		return err
	}
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
}


func (k Keeper) AddBlockRDS(ctx sdk.Context, list ...types.BlockRD) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	
	height := strconv.FormatInt(ctx.BlockHeight(), 10)
	//KEY
	storeKey := types.BlockRelationDataKey + height

	var finalList types.BlockRDS

	
	if store.Has(storeKey) {
		bLockRdListBytes := store.Get(storeKey)
		err := util.Json.Unmarshal(bLockRdListBytes, &finalList)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
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
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	key := types.CopyrightCrossChainOutFeeRatioKey
	if store.Has(key) {
		ratioByte := store.Get(key)
		ratioString := string(ratioByte)
		ratioDec, err := sdk.NewDecFromStr(ratioString)
		if err != nil {
			log.WithError(err).Error("NewDecFromStr")
			return "", err
		}
		if ratioDec.LT(sdk.NewDec(0)) {
			return "", errors.New("Fee Set Error")
		}
		return ratioString, nil
	} else {
		return core.CrossChainOutFeeRatio, nil
	}
}
