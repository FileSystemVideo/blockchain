package keeper

import (
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"strconv"
)

const (
	spaceMinerKey          = "space_miner"           
	spaceTotalIndexKey     = "space_total_index"     
	spaceMinerAccountKey   = "space_miner_account"   
	spaceMinerAmountKey    = "space_miner_amount"    
	spaceMinerBonusKey     = "space_miner_bonus_"    
	deflationMinerKey      = "deflation_miner_key_"  
	deflationSpaceTotalKey = "deflation_space_total" 
)

var ByteToMb = decimal.RequireFromString("1048576")
var returnDays = decimal.RequireFromString("120")


//var spaceMinerUpperLimitStandand = decimal.RequireFromString("100000000")


//var validatorMinerPerDayStandand = decimal.RequireFromString("100")

type AccountSpaceMiner struct {
	Account         string               `json:"account"`          
	SpaceTotal      decimal.Decimal      `json:"space_total"`      // b
	UsedSpace       decimal.Decimal      `json:"used_space"`       // b
	DeflationAmount types.RealCoin       `json:"deflation_amount"` 
	BuySpace        decimal.Decimal      `json:"buy_space"`        
	RewardSpace     decimal.Decimal      `json:"reward_space"`     
	Settlement      map[int64]Settlement `json:"settlement"`       
	SettlementEnd   Settlement           `json:"settlement_end"`   
	LockedSpace     decimal.Decimal      `json:"locked_space"`     // b ,
}

type Settlement struct {
	Index      int64           `json:"index"`       
	IndexSpace decimal.Decimal `json:"index_space"` 
}

type RateAndVoteIndex struct {
	DeflationRate      string `json:"deflation_rate"`       
	DeflationVoteIndex string `json:"deflation_vote_index"` 
}

type DeflationMinerInfor struct {
	DeflationStatus   int             `json:"deflation_status"`    
	MinerTotalAmount  decimal.Decimal `json:"miner_total_amount"`  
	HasMinerAmount    decimal.Decimal `json:"has_miner_amount"`    
	RemainMinerAmount decimal.Decimal `json:"remain_miner_amount"` 
	DayMinerAmount    decimal.Decimal `json:"day_miner_amount"`    
	DayMinerRemain    int64           `json:"day_miner_remain"`    
	MinerBlockNum     int64           `json:"miner_block_num"`     
}


type DeflationAccountMinerInfor struct {
	DeflationMinerInfor
	SpaceTotal        string          `json:"space_total"`         
	AccountSpaceTotal decimal.Decimal `json:"account_space_total"` 
	SpacePercent      string          `json:"space_percent"`       
	LeftBlockNum      int64           `json:"left_block_num"`      
	ExpectBonus       string          `json:"expect_bonus"`        
	LeftReachHeight   int64           `json:"left_reach_height"`   
	DestoryAddress    string          `json:"destory_address"`     
	DestoryAmount     string          `json:"destory_amount"`      
	SpaceMinerReward  string          `json:"space_miner_reward"`  
}


func (k Keeper) SpaceMinerBonusNew(ctx sdk.Context) decimal.Decimal {
	store := k.KVHelper(ctx)
	bz := store.Get(spaceMinerAccountKey)
	if len(bz) > 0 {
		bonusDecimal := k.QuerySpaceMinerBonusAmount(ctx)
		bonusDec := sdk.NewDec(bonusDecimal.IntPart())
		
		if bonusDec.IsPositive() {
			return bonusDecimal
		}
	}
	return decimal.Decimal{}
}

/*
	
*/
/*
	Account         string               `json:"account"`          
	SpaceTotal      decimal.Decimal      `json:"space_total"`      // b
	UsedSpace       decimal.Decimal      `json:"used_space"`       // b
	DeflationAmount types.RealCoin       `json:"deflation_amount"` 
	BuySpace        decimal.Decimal      `json:"buy_space"`        
	RewardSpace     decimal.Decimal      `json:"reward_space"`     
	Settlement      map[int64]Settlement `json:"settlement"`       
	SettlementEnd   Settlement           `json:"settlement_end"`   
*/
func (k Keeper) QueryAccountMinerBonusAmount(ctx sdk.Context, account string) (string, AccountSpaceMiner, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	
	accountSpaceMiner := k.QueryAccountSpaceMinerInfor(ctx, account)

	//nil,0,0,
	if (accountSpaceMiner.Settlement == nil || len(accountSpaceMiner.Settlement) == 0) && accountSpaceMiner.SettlementEnd.Index == 0 {
		return "0", accountSpaceMiner, types.NoSpaceMinerRewardErr
	}
	height := ctx.BlockHeight()
	
	if height <= core.MinerStartHeight {
		return "0", accountSpaceMiner, types.NoSpaceMinerRewardErr
	}
	
	index := (height-core.MinerStartHeight)/core.SpaceMinerBonusBlockNum + 1
	var indexDiff int
	var indexArray []int64
	indexEnd := accountSpaceMiner.SettlementEnd.Index
	
	indexDiff = int(index - indexEnd)
	if indexDiff <= 0 {
		return "0", accountSpaceMiner, types.NoSpaceMinerRewardErr
	}
	for i := 0; i < indexDiff; i++ {
		indexArray = append(indexArray, indexEnd)
		indexEnd = indexEnd + 1
	}
	log.WithField("Number of dividend periods", indexArray).Debug("QueryAccountMinerBonusAmount")
	
	var totalSpace decimal.Decimal   
	var accountBonus decimal.Decimal 
	var accountBonusString string
	
	bonusDecimal := k.QuerySpaceMinerBonusAmount(ctx)
	bonusDec := sdk.NewDec(bonusDecimal.IntPart())
	
	if !bonusDec.IsPositive() {
		return "0", accountSpaceMiner, nil
	}
	//---------------
	
	totalSpace = accountSpaceMiner.SettlementEnd.IndexSpace
	for _, index := range indexArray {
		
		spaceTotal := k.QuerySpaceTotalIndex(ctx, index)
		
		if spaceTotal.Sign() <= 0 {
			continue
		}
		
		bonusDec, err := sdk.NewDecFromStr(bonusDecimal.String())
		if err != nil {
			continue
		}
		spaceTotalDec, err := sdk.NewDecFromStr(spaceTotal.String())
		if err != nil {
			continue
		}
		perBonus := bonusDec.Quo(spaceTotalDec)
		log.WithFields(logrus.Fields{
			"Amount per share":     perBonus.String(),
			"Number of executions": index}).Debug("QueryAccountMinerBonusAmount")
		
		if _, ok := accountSpaceMiner.Settlement[index]; ok {
			totalSpace = accountSpaceMiner.Settlement[index].IndexSpace
		}
		
		accountBonus = accountBonus.Add(totalSpace.Mul(decimal.RequireFromString(perBonus.String())))
	}
	accountBonusString = util.DecimalStringFixed(accountBonus.String(), core.CoinPlaces)
	
	accountSpaceMiner.SettlementEnd.Index = index
	
	accountSpaceMiner.SettlementEnd.IndexSpace = accountSpaceMiner.SpaceTotal
	//Map
	accountSpaceMiner.Settlement = nil
	
	return accountBonusString, accountSpaceMiner, nil
}


func (k Keeper) SettlementAmount(ctx sdk.Context, account string) (string, error) {
	accountBonusString, accountSpaceMiner, err := k.QueryAccountMinerBonusAmount(ctx, account)
	if err != nil {
		return "0", err
	}
	
	k.SetAccountSpaceMinerInfor(ctx, accountSpaceMiner)
	return accountBonusString, nil
}


func (k Keeper) QuerySpaceTotalIndex(ctx sdk.Context, index int64) decimal.Decimal {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(spaceTotalIndexKey + strconv.Itoa(int(index)))
	if bz != nil {
		log.WithField("index", string(bz)).Debug("QuerySpaceTotalIndex")
		return decimal.RequireFromString(string(bz))
	}
	return decimal.Decimal{}
}


func (k Keeper) SetSpaceTotalIndex(ctx sdk.Context) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	spaceTotal := k.QueryDeflatinSpaceTotal(ctx)
	height := ctx.BlockHeight()
	index := (height - core.MinerStartHeight) / core.SpaceMinerBonusBlockNum
	err := store.Set(spaceTotalIndexKey+strconv.Itoa(int(index)), spaceTotal)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"index":      index,
			"spaceTotal": spaceTotal,
		}).Error("store.Set")
		panic(err)
	}
}


func (k Keeper) SpaceMinerRewardSettlement(ctx sdk.Context, account string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	accountBonusString, err := k.SettlementAmount(ctx, account)
	if err != nil {
		return err
	}
	realCoinBonus := types.NewRealCoinFromStr(sdk.DefaultBondDenom, accountBonusString)
	ledgeCoin := types.MustRealCoin2LedgerCoin(realCoinBonus)
	bonusCoins := sdk.NewCoins(ledgeCoin)
	if bonusCoins.IsZero() {
		return types.NoSpaceMinerRewardErr
	}
	accAddress, err := sdk.AccAddressFromBech32(account)
	if err != nil {
		log.WithError(err).WithField("account", account).Error("AccAddressFromBech32")
		return err
	}
	err = k.CoinKeeper.SendCoinsFromModuleToAccount(ctx, core.KeyCopyrighDeflation, accAddress, bonusCoins)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"accAddress": accAddress.String(),
			"bonusCoins": bonusCoins.String(),
		}).Error("SendCoinsFromModuleToAccount")
		return err
	}
	
	err = k.SetSpaceMinerBonusAmount(ctx, decimal.RequireFromString(accountBonusString))
	if err != nil {
		return err
	}
	
	feeRealCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, core.ChainDefaultFee.String())
	var tradeRd []types.BlockRD
	tradeInfor := types.TradeInfor{
		From:           core.ContractAddressDeflation.String(),
		To:             account,
		Amount:         realCoinBonus,
		Fee:            feeRealCoin,
		BlockNum:       ctx.BlockHeight(),
		TradeType:      core.TradeTypeSpaceMinerBonus,
		TimeStemp:      ctx.BlockTime().Unix(),
		Txhash:         formatTxHash(ctx.TxBytes()),
		FromFsvBalance: k.GetBalance(ctx, core.MainToken, core.ContractAddressDeflation),
		FromTipBalance: k.GetBalance(ctx, core.InviteToken, core.ContractAddressDeflation),
		ToFsvBalance:   k.GetBalance(ctx, core.MainToken, accAddress),
		ToTipBalance:   k.GetBalance(ctx, core.InviteToken, accAddress),
	}
	tradeRd = append(tradeRd, types.NewBlockRD(tradeInfor))
	err = k.AddBlockRDS(ctx, tradeRd...)
	if err != nil {
		log.WithError(err).Error("AddBlockRDS")
		return err
	}
	return nil
}


func (k Keeper) AddSpaceMiner(ctx sdk.Context, spaceMiner types.SpaceMinerData) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	
	accountMiner := k.QueryAccountSpaceMinerInfor(ctx, spaceMiner.AwardAccount.String())
	
	spaceMinerDecimal := spaceMiner.DeflationAmount.AmountDec()
	
	if spaceMiner.AwardAccount.Equals(spaceMiner.Creator) {
		
		log.WithField("Pre destruction data", accountMiner.DeflationAmount.String()).Debug("AddSpaceMiner")
		
		if accountMiner.DeflationAmount.Denom == "" {
			accountMiner.DeflationAmount.Denom = spaceMiner.DeflationAmount.Denom
			accountMiner.DeflationAmount.Amount = "0"
		}
		
		amountDecimal := accountMiner.DeflationAmount.AmountDec()
		
		amountDecimal = amountDecimal.Add(spaceMinerDecimal)
		
		accountMiner.DeflationAmount = types.NewRealCoinFromStr(sdk.DefaultBondDenom, amountDecimal.String())
		
		log.WithField("Post destruction data", accountMiner.DeflationAmount.String()).Debug("AddSpaceMiner")
	} else { 
		creatorAccountMiner := k.QueryAccountSpaceMinerInfor(ctx, spaceMiner.Creator.String())
		
		if creatorAccountMiner.DeflationAmount.Denom == "" {
			creatorAccountMiner.DeflationAmount.Denom = spaceMiner.DeflationAmount.Denom
			creatorAccountMiner.DeflationAmount.Amount = "0"
		}
		
		amountDecimal := creatorAccountMiner.DeflationAmount.AmountDec()
		
		amountDecimal = amountDecimal.Add(spaceMinerDecimal)
		
		creatorAccountMiner.DeflationAmount = types.NewRealCoinFromStr(sdk.DefaultBondDenom, amountDecimal.String())
		
		k.SetAccountSpaceMinerInfor(ctx, creatorAccountMiner)
	}

	
	amount := spaceMiner.DeflationAmount.AmountDec()
	
	rate := k.QueryDeflationRate(ctx)
	
	space := amount.Div(decimal.RequireFromString(rate)).Mul(ByteToMb) 
	
	accountMiner.SpaceTotal = accountMiner.SpaceTotal.Add(space).Round(4)
	
	accountMiner.BuySpace = accountMiner.BuySpace.Add(space).Round(4)

	
	accountMiner.Account = spaceMiner.AwardAccount.String()
	
	deflationCoin := types.MustRealCoin2LedgerCoin(spaceMiner.DeflationAmount)
	
	destoryCoins := sdk.NewCoins(deflationCoin)
	
	err := k.CoinKeeper.SendCoins(ctx, spaceMiner.Creator, core.ContractAddressDestory, destoryCoins)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"fromAddr": spaceMiner.Creator.String(),
			"toAddr":   core.ContractAddressDestory,
			"coins":    destoryCoins.String(),
		}).Error("SendCoins")
		return err
	}
	accountMiner = k.calSettlementMap(ctx, accountMiner)
	
	k.SetAccountSpaceMinerInfor(ctx, accountMiner)
	
	k.SetSpaceMinerAccount(ctx, spaceMiner.Creator.String())
	
	err = k.SetSpaceMinerAmount(ctx, spaceMinerDecimal)
	if err != nil {
		return err
	}
	
	k.SetDeflationSpaceTotal(ctx, space)
	
	txhash := formatTxHash(ctx.TxBytes())
	
	txSpaceMiner := types.SpaceMiner{
		Txhash: txhash,
		Space:  space.StringFixed(4),
	}
	return k.AddBlockRDS(ctx, types.NewBlockRD(txSpaceMiner))
	
	//return k.InviteReward(ctx, space.Round(4), spaceMiner.AwardAccount, 1)
}

func (k Keeper) SetAccountSpaceMinerInfor(ctx sdk.Context, accountSpace AccountSpaceMiner) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	key := spaceMinerKey + accountSpace.Account
	err := store.Set(key, accountSpace)
	if err != nil {
		log.WithError(err).Error("store.Set")
		panic(err)
	}
}


func (k Keeper) QueryDeflatinSpaceTotal(ctx sdk.Context) string {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(deflationSpaceTotalKey)
	log.WithField("Storage space size", string(bz)).Debug("QueryDeflatinSpaceTotal")
	if bz != nil {
		return string(bz)
	} else {
		return "0"
	}
}


func (k Keeper) SetDeflationSpaceTotal(ctx sdk.Context, space decimal.Decimal) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var amountSpaceMiner decimal.Decimal
	if !store.Has(deflationSpaceTotalKey) {
		amountSpaceMiner = space
	} else {
		amountSpaceMiner = decimal.RequireFromString(string(store.Get(deflationSpaceTotalKey)))
		amountSpaceMiner = amountSpaceMiner.Add(space)
	}
	err := store.Set(deflationSpaceTotalKey, amountSpaceMiner.StringFixed(4))
	if err != nil {
		log.WithError(err).Error("store.Set")
		panic(err)
	}
}


func (k Keeper) QuerySpaceMinerAmount(ctx sdk.Context) decimal.Decimal {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var amountSpaceMiner decimal.Decimal
	if !store.Has(spaceMinerAmountKey) {
		return amountSpaceMiner
	}
	err := store.GetUnmarshal(spaceMinerAmountKey, &amountSpaceMiner)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		panic(err)
	}
	return amountSpaceMiner
}


func (k Keeper) SetSpaceMinerAmount(ctx sdk.Context, amount decimal.Decimal) (err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var amountSpaceMiner decimal.Decimal
	
	if store.Has(spaceMinerAmountKey) {
		err := store.GetUnmarshal(spaceMinerAmountKey, &amountSpaceMiner)
		if err != nil {
			log.WithError(err).Error("GetUnmarshal")
			return err
		}
		amountSpaceMiner = amountSpaceMiner.Add(amount)
	} else {
		amountSpaceMiner = amount
	}
	err = store.Set(spaceMinerAmountKey, amountSpaceMiner)
	if err != nil {
		log.WithError(err).Error("store.Set")
	}
	return err
}


func (k Keeper) UpdateAccountSpace(ctx sdk.Context, address sdk.Address, size int64) bool {
	accountSpace := k.QueryAccountSpaceMinerInfor(ctx, address.String())
	addSpace := decimal.NewFromFloat(float64(size))
	if accountSpace.UsedSpace.GreaterThanOrEqual(addSpace) { 
		accountSpace.UsedSpace = accountSpace.UsedSpace.Sub(addSpace).Round(4)
		k.SetAccountSpaceMinerInfor(ctx, accountSpace)
	} else {
		return false
	}
	return true
}


func (k Keeper) SetSpaceMinerAccount(ctx sdk.Context, account string) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	if store.Has(spaceMinerAccountKey) {
		return
	}
	accountString := account
	err := store.Set(spaceMinerAccountKey, accountString)
	if err != nil {
		log.WithError(err).Error("store.Set")
		panic(err)
	}
}


func (k Keeper) QueryAccountSpaceMinerInfor(ctx sdk.Context, address string) AccountSpaceMiner {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var accountSpaceMiner AccountSpaceMiner
	bz := store.Get(spaceMinerKey + address)
	if len(bz) > 0 {
		err := util.Json.Unmarshal(bz, &accountSpaceMiner)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			panic(err)
		}
	}
	return accountSpaceMiner
}


func (k Keeper) LockAccountSpaceMiner(ctx sdk.Context, address string, lockedSize int64) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var accountSpaceMiner AccountSpaceMiner
	bz := store.Get(spaceMinerKey + address)
	if len(bz) > 0 {
		err := util.Json.Unmarshal(bz, &accountSpaceMiner)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return err
		}
		// := --
		leftSpace := accountSpaceMiner.SpaceTotal.Sub(accountSpaceMiner.UsedSpace).Sub(accountSpaceMiner.LockedSpace)
		lockedDecimal := decimal.NewFromInt(lockedSize)
		if leftSpace.Sub(lockedDecimal).Sign() >= 0 {
			accountSpaceMiner.LockedSpace = accountSpaceMiner.LockedSpace.Add(lockedDecimal)
			err = store.Set(spaceMinerKey+address, accountSpaceMiner)
			if err != nil {
				log.WithError(err).Error("store.Set")
				return err
			}
		}
	} else {
		return types.SpaceFirstError
	}
	return nil
}


func (k Keeper) LockAccountSpaceMinerReturn(ctx sdk.Context, address string, lockedSize int64) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var accountSpaceMiner AccountSpaceMiner
	bz := store.Get(spaceMinerKey + address)
	if len(bz) > 0 {
		err := util.Json.Unmarshal(bz, &accountSpaceMiner)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			panic(err)
		}
		lockedDecimal := decimal.NewFromInt(lockedSize)
		accountSpaceMiner.LockedSpace = accountSpaceMiner.LockedSpace.Sub(lockedDecimal)
		err = store.Set(spaceMinerKey+address, accountSpaceMiner)
		if err != nil {
			log.WithError(err).Error("store.Set")
			return err
		}
	} else {
		return types.AccountSpaceReturnError
	}
	return nil
}


func (k Keeper) UpdateDeflationMinerInforByVote(ctx sdk.Context, minerAmount decimal.Decimal) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var deflationMiner DeflationMinerInfor
	deflationMiner, _ = k.QueryDeflationMinerInfor(ctx)

	
	deflationMiner.HasMinerAmount = deflationMiner.HasMinerAmount.Add(minerAmount)
	
	remainMiner := deflationMiner.MinerTotalAmount.Sub(deflationMiner.HasMinerAmount)
	if remainMiner.Sign() <= 0 { 
		remainMiner = decimal.Zero
	}
	
	deflationMiner.RemainMinerAmount = remainMiner
	
	deflationMiner.DayMinerRemain = remainMiner.DivRound(minerAmount, 0).IntPart()

	err := store.Set(deflationMinerKey, deflationMiner)
	if err != nil {
		log.WithError(err).Error("store.Set")
		panic(err)
	}
}


func (k Keeper) UpdateDeflationMinerInfor(ctx sdk.Context, minerAmount decimal.Decimal, blockNum int64) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var deflationMiner DeflationMinerInfor
	deflationMiner, _ = k.QueryDeflationMinerInfor(ctx)

	deflationMiner.DeflationStatus = 1
	deflationMiner.MinerBlockNum = blockNum
	deflationMiner.DayMinerAmount = minerAmount
	
	deflationMiner.HasMinerAmount = deflationMiner.HasMinerAmount.Add(minerAmount)
	
	remainMiner := deflationMiner.MinerTotalAmount.Sub(deflationMiner.HasMinerAmount)
	if remainMiner.Sign() <= 0 { 
		remainMiner = decimal.Zero
	}
	
	deflationMiner.RemainMinerAmount = remainMiner
	
	deflationMiner.DayMinerRemain = remainMiner.DivRound(minerAmount, 0).IntPart()

	err := store.Set(deflationMinerKey, deflationMiner)
	if err != nil {
		log.WithError(err).Error("store.Set")
		panic(err)
	}
}

func (k Keeper) QueryDeflationMinerInfor(ctx sdk.Context) (DeflationMinerInfor, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var deflationMiner DeflationMinerInfor
	if !store.Has(deflationMinerKey) {
		deflationMiner = DeflationMinerInfor{
			DeflationStatus:   0,
			DayMinerRemain:    core.DayMinerRemain,
			DayMinerAmount:    core.SpaceMinerPerDayStandand.Add(core.ValidatorMinerPerDayStandand),
			MinerTotalAmount:  core.MinerUpperLimitStandand,
			HasMinerAmount:    core.ChuangshiFee,
			RemainMinerAmount: core.MinerUpperLimitStandand.Sub(core.ChuangshiFee),
			MinerBlockNum:     core.MinerStartHeight,
		}
		if ctx.BlockHeight() > core.MinerStartHeight {
			deflationMiner.DeflationStatus = 1
		}
		return deflationMiner, nil
	}
	err := store.GetUnmarshal(deflationMinerKey, &deflationMiner)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return DeflationMinerInfor{}, err
	}
	return deflationMiner, nil
}


func (k Keeper) QuerySpaceMinerBonusAmount(ctx sdk.Context) decimal.Decimal {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	if !store.Has(spaceMinerBonusKey) {
		return core.SpaceMinerPerDayStandand
	} else {
		var hasBonusAmount decimal.Decimal
		err := store.GetUnmarshal(spaceMinerBonusKey, &hasBonusAmount)
		if err != nil {
			log.WithError(err).Error("GetUnmarshal")
			panic(err)
		}
		if core.MinerUpperLimitStandand.LessThanOrEqual(hasBonusAmount.Add(core.SpaceMinerPerDayStandand)) {
			if core.MinerUpperLimitStandand.GreaterThan(hasBonusAmount) {
				currentBonus := core.MinerUpperLimitStandand.Sub(hasBonusAmount)
				return currentBonus
			} else {
				return decimal.Decimal{}
			}
		} else { 
			return core.SpaceMinerPerDayStandand
		}
	}
}


func (k Keeper) QueryValidatorMinerBonusAmount(ctx sdk.Context) decimal.Decimal {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	if !store.Has(spaceMinerBonusKey) {
		return core.ValidatorMinerPerDayStandand
	}
	var hasBonusAmount decimal.Decimal
	err := store.GetUnmarshal(spaceMinerBonusKey, &hasBonusAmount)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		panic(err)
	}
	
	if core.MinerUpperLimitStandand.LessThanOrEqual(hasBonusAmount.Add(core.ValidatorMinerPerDayStandand)) {
		if core.MinerUpperLimitStandand.GreaterThan(hasBonusAmount) {
			currentBonus := core.MinerUpperLimitStandand.Sub(hasBonusAmount)
			return currentBonus
		} else {
			return decimal.Decimal{}
		}
	} else { 
		return core.ValidatorMinerPerDayStandand
	}
}


func (k Keeper) QueryHasMinerBonusAmount(ctx sdk.Context) decimal.Decimal {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	if !store.Has(spaceMinerBonusKey) {
		return core.ValidatorMinerPerDayStandand
	}
	var hasBonusAmount decimal.Decimal
	err := store.GetUnmarshal(spaceMinerBonusKey, &hasBonusAmount)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		panic(err)
	}
	return hasBonusAmount
}


func (k Keeper) SetSpaceMinerBonusAmount(ctx sdk.Context, amount decimal.Decimal) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(spaceMinerBonusKey)
	var amountSpaceMiner decimal.Decimal
	if bz != nil {
		err := util.Json.Unmarshal(bz, &amountSpaceMiner)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return err
		}
		amountSpaceMiner = amountSpaceMiner.Add(amount)
	} else {
		amountSpaceMiner = amount
	}
	err := store.Set(spaceMinerBonusKey, amountSpaceMiner)
	if err != nil {
		log.WithError(err).Error("store.Set")
	}
	return err
}


func (k Keeper) ValidatorBonus(ctx sdk.Context) decimal.Decimal {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	
	bonusDecimal := k.QueryValidatorMinerBonusAmount(ctx)
	//zeroDecimal := decimal.RequireFromString("0")
	if bonusDecimal.Sign() > 0 {
		realCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, bonusDecimal.String())
		bonusLedgerCoin := types.MustRealCoin2LedgerCoin(realCoin)
		err := k.CoinKeeper.SendCoinsFromModuleToModule(ctx, core.KeyCopyrighDeflation, authtypes.FeeCollectorName, sdk.NewCoins(bonusLedgerCoin))
		if err != nil {
			log.WithError(err).Error("SendCoinsFromModuleToModule")
			panic(err)
		}
		
		err = k.SetSpaceMinerBonusAmount(ctx, bonusDecimal)
		if err != nil {
			return decimal.Zero
		}
		feeRealCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, "0")
		
		tradeInfor := types.TradeInfor{
			From:           core.ContractAddressDeflation.String(),
			To:             core.ContractAddressFee.String(),
			Amount:         realCoin,
			Fee:            feeRealCoin,
			BlockNum:       ctx.BlockHeight(),
			TradeType:      core.TradeTypeValidatorMinerBonus,
			TimeStemp:      ctx.BlockTime().Unix(),
			FromFsvBalance: k.GetBalance(ctx, core.MainToken, core.ContractAddressDeflation),
			FromTipBalance: k.GetBalance(ctx, core.InviteToken, core.ContractAddressDeflation),
			ToFsvBalance:   k.GetBalance(ctx, core.MainToken, core.ContractAddressFee),
			ToTipBalance:   k.GetBalance(ctx, core.InviteToken, core.ContractAddressFee),
		}
		err = k.AddBlockRDS(ctx, types.NewBlockRD(tradeInfor))
		if err != nil {
			log.WithError(err).Error("AddBlockRDS")
		}
		return bonusDecimal
	}
	return decimal.Zero

}


func (k Keeper) UpdateAccountSpaceUsed(ctx sdk.Context, address sdk.Address, size int64) bool {
	accountSpace := k.QueryAccountSpaceMinerInfor(ctx, address.String())
	addSpace := decimal.New(size, 0)
	if accountSpace.SpaceTotal.Sub(accountSpace.UsedSpace).LessThan(addSpace) {
		return false
	}
	accountSpace.UsedSpace = accountSpace.UsedSpace.Add(addSpace).Round(4)
	k.SetAccountSpaceMinerInfor(ctx, accountSpace)
	return true
}


func (k Keeper) setDistributeFee(ctx sdk.Context, leftBonus decimal.Decimal) {
	feePool := k.distributionKeeper.GetFeePool(ctx)
	leftRealCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, util.DecimalStringFixed(leftBonus.String(), core.CoinPlaces))
	ledgeCoin := types.MustRealCoin2LedgerCoin(leftRealCoin)
	ledgeDecCoin := sdk.NewDecCoinFromDec(ledgeCoin.Denom, ledgeCoin.Amount.ToDec())
	feePool.CommunityPool = feePool.CommunityPool.Add(ledgeDecCoin)
	k.distributionKeeper.SetFeePool(ctx, feePool)
}
func (k Keeper) setDistributeFeeV2(ctx sdk.Context, contractAddress string, ledgeCoin sdk.Coin) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	feePool := k.distributionKeeper.GetFeePool(ctx)
	ledgeDecCoin := sdk.NewDecCoinFromDec(ledgeCoin.Denom, ledgeCoin.Amount.ToDec())
	err := k.CoinKeeper.SendCoinsFromModuleToModule(ctx, contractAddress, distributionTypes.ModuleName, sdk.NewCoins(ledgeCoin))
	if err != nil {
		log.WithError(err).WithField("ledgecoin", ledgeCoin.String()).Error("setDistributeFeeV2")
		return
	}
	feePool.CommunityPool = feePool.CommunityPool.Add(ledgeDecCoin)
	k.distributionKeeper.SetFeePool(ctx, feePool)
}


/*
1GB - > 1024GB	0.01
1TB -> 1024TB	0.009
1PB -> 1024PB   0.008
1EB -> 1024EB   0.007
1ZB -> 1024ZB   0.006
1YB -> 1024ZB   0.005
1BB -> 1024YB   0.004
1NB -> 1024YB   0.003
1DB -> 1024NB   0.001
*/
var (
	spaceUnit = decimal.NewFromInt(1024)
	gB        = decimal.NewFromInt(1024 * 1024 * 1024 * 1024)
	tB        = gB.Mul(spaceUnit)
	pB        = tB.Mul(spaceUnit)
	eB        = pB.Mul(spaceUnit)
	zB        = eB.Mul(spaceUnit)
	yB        = zB.Mul(spaceUnit)
	bB        = yB.Mul(spaceUnit)
	nB        = bB.Mul(spaceUnit)
	dB        = nB.Mul(spaceUnit)
	maxDB     = dB.Mul(spaceUnit)
)

var spaceFeeMap = map[decimal.Decimal]string{
	gB:    "0.01",
	tB:    "0.009",
	pB:    "0.008",
	eB:    "0.007",
	zB:    "0.006",
	yB:    "0.005",
	bB:    "0.004",
	nB:    "0.003",
	dB:    "0.002",
	maxDB: "0.001",
}


func (k Keeper) SpaceFeeEstimate(ctx sdk.Context) string {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	totalSpaceStr := k.QueryDeflatinSpaceTotal(ctx)
	totalSpace, err := decimal.NewFromString(totalSpaceStr)
	if err != nil {
		log.WithError(err).Error("NewFromString")
		panic(err)
	}
	if totalSpace.LessThanOrEqual(gB) {
		return spaceFeeMap[gB]
	} else if totalSpace.LessThanOrEqual(tB) {
		return spaceFeeMap[tB]
	} else if totalSpace.LessThanOrEqual(pB) {
		return spaceFeeMap[pB]
	} else if totalSpace.LessThanOrEqual(eB) {
		return spaceFeeMap[eB]
	} else if totalSpace.LessThanOrEqual(zB) {
		return spaceFeeMap[zB]
	} else if totalSpace.LessThanOrEqual(yB) {
		return spaceFeeMap[yB]
	} else if totalSpace.LessThanOrEqual(bB) {
		return spaceFeeMap[bB]
	} else if totalSpace.LessThanOrEqual(nB) {
		return spaceFeeMap[nB]
	} else if totalSpace.LessThanOrEqual(dB) {
		return spaceFeeMap[dB]
	} else if totalSpace.LessThanOrEqual(maxDB) {
		return spaceFeeMap[maxDB]
	}
	return "0.001"
}
