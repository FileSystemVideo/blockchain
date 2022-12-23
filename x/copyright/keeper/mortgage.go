package keeper

import (
	"encoding/hex"
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	tenType "github.com/tendermint/tendermint/types"
	"sort"
	"strconv"
	"strings"
)

const (
	MortgageKey       = "mortgage-for-copyright" //key()
	MiningTotalAmount = "mortgage-total-amount"  //key
	MortgageTimeInfor = "mortgage-time-infor-"   
)


type MortgageAccountAmount struct {
	MortageAccount sdk.AccAddress `json:"mortage_account"` 
	MortgageAmount types.RealCoin `json:"mortgage_amount"` 
	CopyrightPrice types.RealCoin `json:"copyright_price"` 
	MortgTime      int64          `json:"mortg_time"`      
	TxHashAccount  sdk.AccAddress `json:"txhash_account"`  //txhash
	Txhash         string         `json:"txhash"`          //txhash 
}


type TimeMortgage struct {
	MortgageAccount []MortgageAccountAmount `json:"mortgage_account"` 
	NextTime        int64                   `json:"next_time"`        
}


type MortgageCopyrightInfor struct {
	//MortgageAccountArray []MortgageAccountAmount `json:"mortage_account_list"`  
	//MortgageTotalAmount  sdk.Coin                `json:"mortgage_total_amount"` 
	TriggerTime int64 `json:"trigger_time"` 
	LastTime    int64 `json:"last_time"`    
}

func formatTxHash(tx tenType.Tx) string {
	return strings.ToUpper(hex.EncodeToString(tx.Hash()))
}


type MortgMinerInfor struct {
	Status          int    `json:"status"`           
	ContractAddress string `json:"contract_address"` 
	MinerTotal      string `json:"miner_total"`      
	HasMiner        string `json:"has_miner"`        
	RedeemTime      int64  `json:"redeem_time"`      
	ValidTime       int64  `json:"valid_time"`       
}


func (k Keeper) Mortgage(ctx sdk.Context, mortgage types.MortgageData) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	txhash := formatTxHash(ctx.TxBytes())
	blockHeight := ctx.BlockHeight()
	
	copyrightByte := k.GetCopyrightBonusInfo(ctx, mortgage.DataHash, mortgage.MortgageAccount.String())
	resourceRelation := CopyrightExtrainfor{}
	if copyrightByte != nil {
		err := util.Json.Unmarshal(copyrightByte, &resourceRelation)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return errors.New("format mortgage error")
		}
		
		if resourceRelation.Species == "buy" || resourceRelation.Height > blockHeight {
			return types.CopyrightMortgageErr
		}
	}
	if blockHeight < core.MortgageStartHeight { 
		return errors.New("has not reach mortgage height")
	}
	txhashMortgAccountString, err := GetTxHashAccount(txhash)
	if err != nil {
		log.WithError(err).WithField("txhash", txhash).Error("GetTxHashAccount")
		return errors.New("format mortgage account error")
	}
	//txhash,
	txhashMortgAccount, err := sdk.AccAddressFromBech32(txhashMortgAccountString)
	if err != nil {
		log.WithError(err).WithField("txhashMortgAccountString", txhashMortgAccountString).Error("AccAddressFromBech32")
		return err
	}
	ledgeCoin := types.MustRealCoin2LedgerCoin(mortgage.MortgageAmount)
	mortgageCoins := sdk.NewCoins(ledgeCoin)
	
	err = k.CoinKeeper.SendCoins(ctx, mortgage.MortgageAccount, txhashMortgAccount, mortgageCoins)
	//err = k.CoinKeeper.SendCoinsFromAccountToModule(ctx, mortgage.MortgageAccount, config.KeyCopyrightMortgage, mortgageCoins)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"fromAddr": mortgage.MortgageAccount.String(),
			"toAddr":   txhashMortgAccount.String(),
			"amt":      mortgageCoins.String(),
		}).Error("SendCoins")
		return err
	}
	mortgageCopyrightInfor := k.GetMortgageCopyrightInfor(ctx)
	
	err = appendMortgageAccount(ctx, k, mortgage.CopyrightPrice, mortgage.MortgageAmount, mortgage.MortgageAccount, txhashMortgAccount, txhash, mortgage.CreateTime, mortgageCopyrightInfor)
	if err != nil {
		log.WithError(err).Error("appendMortgageAccount")
		return err
	}
	mortgageAmountRealCoin := mortgage.MortgageAmount
	feeCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, "0")
	fee := "0"
	if mortgage.Fee.Amount != nil && mortgage.Fee.Amount.IsValid() {
		log.WithField("Current handling fee amount", mortgage.Fee.Amount[0].Amount.String()).Debug("Mortgage")
		fee = types.MustParseLedgerInt(mortgage.Fee.Amount[0].Amount)
		coin := sdk.NewCoin(sdk.DefaultBondDenom, mortgage.Fee.Amount[0].Amount)
		feeCoin = types.MustLedgerCoin2RealCoin(coin)
	}

	
	mortgTradeInfor := types.TradeInfor{
		From:           mortgage.MortgageAccount.String(),
		To:             txhashMortgAccount.String(),
		Txhash:         txhash,
		TradeType:      core.TradeTypeCopyrightBuyMortgage,
		Amount:         mortgageAmountRealCoin,
		Fee:            feeCoin,
		BlockNum:       blockHeight,
		TimeStemp:      ctx.BlockTime().Unix(),
		FromFsvBalance: k.GetBalance(ctx, core.MainToken, mortgage.MortgageAccount),
		FromTipBalance: k.GetBalance(ctx, core.InviteToken, mortgage.MortgageAccount),
		ToFsvBalance:   k.GetBalance(ctx, core.MainToken, txhashMortgAccount),
		ToTipBalance:   k.GetBalance(ctx, core.InviteToken, txhashMortgAccount),
	}
	
	buySourceInfor := types.BuySourceInfor{
		Txhash:    txhash,
		DataHash:  mortgage.DataHash,
		Creator:   mortgage.Creator.String(),
		Purchaser: mortgage.MortgageAccount.String(),
		Price:     mortgage.CopyrightPrice,
		Fee:       feeCoin,
		Remark:    "mortgage",
		TimeStemp: ctx.BlockTime().Unix(),
	}
	
	flag := k.AddMiningAmount(ctx, mortgage.CopyrightPrice)
	if !flag {
		return errors.New(types.MortgMinerHasFinish.Error())
	}
	
	//k.AddSupply(ctx, msg.CopyrightPrice[0])
	
	delayTime := k.judgeMortgRedeemTimeDelay(ctx)
	
	height := int64(delayTime/core.CommitTime) + blockHeight + core.MortgageHeight
	k.SetCopyrightDownRelation(ctx, mortgage.DataHash, mortgage.CopyrightPrice, mortgage.MortgageAccount, "mortgage", height)
	/*mortgageCoin := types.MustRealCoin2LedgerCoin(mortgage.MortgageAmount)
	bonusCoins := sdk.NewCoins(mortgageCoin)
	err = k.CoinKeeper.SendCoins(ctx, mortgage.MortgageAccount, mortgage.DataHashAccount, bonusCoins)*/

	
	copyrightInfoByte, err := k.GetCopyright(ctx, mortgage.DataHash)
	if err != nil {
		return err
	}
	copyrightData := &types.CopyrightData{}
	err = util.Json.Unmarshal(copyrightInfoByte, copyrightData)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	
	if mortgage.BonusType == core.BonusTypeFront {
		err = dealBonusAuthrizeLogic(ctx, k, mortgage.OfferAccountShare, txhash, fee, core.KeyCopyrightMortgage, blockHeight, *copyrightData)
		if err != nil {
			log.WithError(err).Error("dealBonusAuthrizeLogic")
			return err
		}
	} else { 
		ledgeCoin := types.MustRealCoin2LedgerCoin(copyrightData.Price)
		bonusCoins := sdk.NewCoins(ledgeCoin)
		err = k.CoinKeeper.SendCoinsFromModuleToModule(ctx, core.KeyCopyrightMortgage, core.KeyCopyrightBonus, bonusCoins)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{
				"fromAddr": core.KeyCopyrightMortgage,
				"toAddr":   core.KeyCopyrightBonus,
				"amt":      bonusCoins.String(),
			}).Error("SendCoinsFromModuleToModule")
			return err
		}
	}
	err = k.AddBlockRDS(ctx, types.NewBlockRD(mortgTradeInfor))
	if err != nil {
		log.WithError(err).Error("AddBlockRDS-mortgTradeInfor")
		return err
	}
	err = k.AddBlockRDS(ctx, types.NewBlockRD(buySourceInfor))
	if err != nil {
		log.WithError(err).Error("AddBlockRDS2-buySourceInfor")
		return err
	}
	return err

}


func (k Keeper) AddMiningAmount(ctx sdk.Context, mingAmount types.RealCoin) bool {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(MiningTotalAmount)
	var mortgageAmount types.RealCoin
	if len(bz) == 0 {
		mortgageAmount = mingAmount
	} else {
		err := util.Json.Unmarshal(bz, &mortgageAmount)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return false
		}
		mingAmountDecimal := mingAmount.AmountDec().Add(mortgageAmount.AmountDec())
		if mingAmountDecimal.GreaterThan(core.MinerUpperLimitStandand) {
			return false
		}
		mortgageAmount = types.NewRealCoinFromStr(sdk.DefaultBondDenom, mingAmountDecimal.String())
	}
	err := store.Set(MiningTotalAmount, mortgageAmount)
	if err != nil {
		log.WithError(err).Error("store.Set")
		return false
	}
	return true
}


func (k Keeper) GetMiningAmount(ctx sdk.Context) types.RealCoin {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(MiningTotalAmount)
	var miningAmount types.RealCoin
	if bz != nil {
		err := util.Json.Unmarshal(bz, &miningAmount)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return miningAmount
		}
	} else {
		miningAmount = types.NewRealCoinFromStr(sdk.DefaultBondDenom, "0")
	}
	return miningAmount
}


func (k Keeper) QueryMortgMinerInfor(ctx sdk.Context) MortgMinerInfor {
	hasMiner := k.GetMiningAmount(ctx)
	mortgMinerInfor := MortgMinerInfor{
		HasMiner: hasMiner.Amount,
	}
	if ctx.BlockHeader().Height > core.MortgageStartHeight {
		mortgMinerInfor.Status = 1
	}
	addTime := k.judgeMortgRedeemTimeDelay(ctx)
	mortgMinerInfor.ContractAddress = core.ContractAddressMortgage.String()
	mortgMinerInfor.MinerTotal = core.MinerUpperLimitStandand.String()
	mortgMinerInfor.RedeemTime = int64(core.CopyrightMortgRedeemTimePerioad + addTime)
	mortgMinerInfor.ValidTime = core.MortgageHeight
	return mortgMinerInfor
}


func (k Keeper) GetMortgageCopyrightInfor(ctx sdk.Context) MortgageCopyrightInfor {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(MortgageKey)
	var mortgageAccountCopyright = MortgageCopyrightInfor{}
	if bz != nil {
		err := util.Json.Unmarshal(bz, &mortgageAccountCopyright)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			panic(err)
		}
	}
	return mortgageAccountCopyright
}


func (k Keeper) SetMortgageCopyrightInfor(ctx sdk.Context, mortgageCopyright MortgageCopyrightInfor) error {
	store := k.KVHelper(ctx)
	return store.Set(MortgageKey, mortgageCopyright)
}


func (k Keeper) DeleteMortgageCopyrightInfor(ctx sdk.Context) {
	store := k.KVHelper(ctx)
	store.Delete(MortgageKey)
}


func (k Keeper) SetMortgageTimeInfor(ctx sdk.Context, timeMortgage TimeMortgage, mortgageTime int64) error {
	store := k.KVHelper(ctx)
	return store.Set(MortgageTimeInfor+strconv.FormatInt(mortgageTime, 10), timeMortgage)
}


func (k Keeper) GetMortgageTimeInfor(ctx sdk.Context, mortgageTime int64) TimeMortgage {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(MortgageTimeInfor + strconv.FormatInt(mortgageTime, 10))
	var timeMortgage = TimeMortgage{}
	if bz != nil {
		err := util.Json.Unmarshal(bz, &timeMortgage)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			panic(err)
		}
	}
	return timeMortgage
}


func (k Keeper) DeleteMortgageTimeInfor(ctx sdk.Context, mortgageTime int64) {
	store := k.KVHelper(ctx)
	store.Delete(MortgageTimeInfor + strconv.FormatInt(mortgageTime, 10))
}


func appendMortgageAccount(ctx sdk.Context, k Keeper, copyrightPrice, mortgageAmount types.RealCoin, address, txhashMortgAccount sdk.AccAddress, txhash string, timStempstring int64, mortgageCopyrightInfor MortgageCopyrightInfor) error {
	mortAccountAmount := MortgageAccountAmount{}
	mortAccountAmount.CopyrightPrice = copyrightPrice
	mortAccountAmount.MortageAccount = address
	mortAccountAmount.MortgageAmount = mortgageAmount
	mortAccountAmount.TxHashAccount = txhashMortgAccount
	mortAccountAmount.Txhash = txhash
	mortAccountAmount.MortgTime = timStempstring
	
	if mortgageCopyrightInfor.TriggerTime == 0 {
		mortgageCopyrightInfor.TriggerTime = timStempstring
	}
	if mortgageCopyrightInfor.LastTime == timStempstring { 
		lastTimeMortgage := k.GetMortgageTimeInfor(ctx, mortgageCopyrightInfor.LastTime)
		lastTimeMortgage.MortgageAccount = append(lastTimeMortgage.MortgageAccount, mortAccountAmount)
		
		return k.SetMortgageTimeInfor(ctx, lastTimeMortgage, mortgageCopyrightInfor.LastTime)
	} else {
		
		timeMortgage := TimeMortgage{}
		timeMortgage.MortgageAccount = append(timeMortgage.MortgageAccount, mortAccountAmount)
		err := k.SetMortgageTimeInfor(ctx, timeMortgage, timStempstring)
		if err != nil {
			return err
		}
		
		lastTimeMortgage := k.GetMortgageTimeInfor(ctx, mortgageCopyrightInfor.LastTime)
		if lastTimeMortgage.MortgageAccount == nil { 
			mortgageCopyrightInfor.TriggerTime = timStempstring
		} else { 
			lastTimeMortgage.NextTime = timStempstring
			err = k.SetMortgageTimeInfor(ctx, lastTimeMortgage, mortgageCopyrightInfor.LastTime)
			if err != nil {
				return err
			}
		}

		
		mortgageCopyrightInfor.LastTime = timStempstring
		return k.SetMortgageCopyrightInfor(ctx, mortgageCopyrightInfor)
	}
}


func buildBonusRelationAndSize(ctx sdk.Context, keeper Keeper, bonusMap map[string]string, sizeMap map[string]int64, txhash, fee, contractAddressKey, creator, dataHash string, blockHeight int64) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	copntractAddress := authtypes.NewModuleAddress(contractAddressKey)
	var bonusKeySlice []string
	//map key
	for key, _ := range bonusMap {
		bonusKeySlice = append(bonusKeySlice, key)
	}
	sort.Strings(bonusKeySlice)
	//key
	for _, account := range bonusKeySlice {
		//mapkey value
		if amount, ok := bonusMap[account]; ok {
			tradeType := core.TradeTypeCopyrightSharesReward
			toAccount := ""
			if account == "creator" { 
				tradeType = core.TradeTypeCopyrightSell
				toAccount = creator
			} else {
				toAccount = account
			}
			realCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, amount)
			ledgerCoin := types.MustRealCoin2LedgerCoin(realCoin)
			//datahash
			toAddress, err := sdk.AccAddressFromBech32(toAccount)
			if err != nil { 
				log.WithError(err).WithField("toAccount", toAccount).Error("AccAddressFromBech32")

				keeper.setDistributeFeeV2(ctx, contractAddressKey, ledgerCoin)

				continue
			}
			err = keeper.CoinKeeper.SendCoinsFromModuleToAccount(ctx, contractAddressKey, toAddress, sdk.NewCoins(ledgerCoin))
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{
					"fromAddr": contractAddressKey,
					"toAddr":   toAddress.String(),
					"amt":      ledgerCoin.String(),
				}).Error("SendCoinsFromModuleToAccount")
				continue
			}
			
			feeCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, "0")
			mortgTradeInfor := types.TradeInfor{
				From:           copntractAddress.String(),
				To:             toAccount,
				Txhash:         txhash,
				TradeType:      tradeType,
				Amount:         realCoin,
				Fee:            feeCoin,
				BlockNum:       blockHeight,
				TimeStemp:      ctx.BlockTime().Unix(),
				FromFsvBalance: keeper.GetBalance(ctx, core.MainToken, copntractAddress),
				FromTipBalance: keeper.GetBalance(ctx, core.InviteToken, copntractAddress),
				ToFsvBalance:   keeper.GetBalance(ctx, core.MainToken, toAddress),
				ToTipBalance:   keeper.GetBalance(ctx, core.InviteToken, toAddress),
			}
			err = keeper.AddBlockRDS(ctx, types.NewBlockRD(mortgTradeInfor))
			if err != nil {
				log.WithError(err).Error("AddBlockRDS")
			}
		}
	}
	var sizeKeySlice []string
	for key, _ := range sizeMap {
		sizeKeySlice = append(sizeKeySlice, key)
	}
	sort.Strings(sizeKeySlice)
	for _, account := range sizeKeySlice {
		if size, ok := sizeMap[account]; ok {
			sourceSizeInfor := types.SourceSizeInfor{
				Account:  account,
				Txhash:   txhash,
				DataHash: dataHash,
				BlockNum: blockHeight,
				Size:     size,
			}
			err := keeper.AddBlockRDS(ctx, types.NewBlockRD(sourceSizeInfor))
			if err != nil {
				log.WithError(err).Error("AddBlockRDS")
			}
		}
	}
}




func dealBonusAuthrizeLogic(ctx sdk.Context, keeper Keeper, offerAccountShare, txhash, fee, contractAddressKey string, blockHeight int64, copyrightInfo types.CopyrightData) error {
	
	
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	realPriceString := copyrightInfo.Price.Amount
	realPriceDeciaml, err := decimal.NewFromString(realPriceString)
	if err != nil {
		log.WithError(err).Error("NewFromString")
		return err
	}
	reduceFlag := keeper.JudgeCopyrightBonusRate(ctx)
	if reduceFlag { 
		realPriceDeciaml = realPriceDeciaml.Div(decimal.RequireFromString("2"))
	}
	
	chargeRateDecimal := decimal.RequireFromString(copyrightInfo.ChargeRate)
	createBonus := realPriceDeciaml.Mul(decimal.RequireFromString("1").Sub(chargeRateDecimal))
	creatorAward := types.NewLedgerInt(createBonus)
	
	//fmt.Println(creatorAward.GT(sdk.ZeroInt()), "---------")
	offerAwardMap := map[string]string{}
	if creatorAward.GT(sdk.ZeroInt()) {
		
		//offerAwardMap[copyrightInfo.Creator.String()] = createBonus
		//key,
		offerAwardMap["creator"] = createBonus.StringFixedBank(core.CoinPlaces)

		
		
		chargeBonus := realPriceDeciaml.Mul(chargeRateDecimal)

		offerAccountsBonus := chargeBonus
		//offerAccountsBonusFloat64, _ := strconv.ParseFloat(offerAccountsBonus.StringFixedBank(core.CoinPlaces), 64)
		offerAccountsAward := types.NewLedgerInt(offerAccountsBonus)
		if offerAccountsAward.GT(sdk.ZeroInt()) {
			normalAccount, authorizeAccount := anylseAccount(ctx, keeper, offerAccountShare)
			if len(normalAccount.AccoutShares) > 0 && len(authorizeAccount.AccoutShares) > 0 {
				authrizeAward := offerAccountsBonus.Mul(decimal.RequireFromString(core.AuthorizeRate))
				normalAward := offerAccountsBonus.Sub(authrizeAward)
				offerAwardMap = caculateAccountAward(authorizeAccount, authrizeAward, offerAwardMap)
				offerAwardMap = caculateAccountAward(normalAccount, normalAward, offerAwardMap)
			} else if len(normalAccount.AccoutShares) > 0 {
				offerAwardMap = caculateAccountAward(normalAccount, offerAccountsBonus, offerAwardMap)
			} else if len(authorizeAccount.AccoutShares) > 0 {
				offerAwardMap = caculateAccountAward(authorizeAccount, offerAccountsBonus, offerAwardMap)
			}
		}
	}
	
	offerSizeMap, err := caculateCopyrightSize(copyrightInfo.Size, offerAccountShare)
	if err != nil {
		return err
	}
	buildBonusRelationAndSize(ctx, keeper, offerAwardMap, offerSizeMap, txhash, fee, contractAddressKey, copyrightInfo.Creator.String(), copyrightInfo.DataHash, blockHeight)
	return nil
}

// true   false 
func (k Keeper) JudgeCopyrightBonusRate(ctx sdk.Context) bool {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(MiningTotalAmount)
	var mortgageAmount types.RealCoin
	if len(bz) == 0 {
		return false
	} else {
		err := util.Json.Unmarshal(bz, &mortgageAmount)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return false
		}
		mortgageMinerAmount := mortgageAmount.AmountDec()
		if mortgageMinerAmount.LessThan(core.MinerUpperLimitStandand) {
			return false
		}
	}
	return true
}


func caculateAccountAward(accountShareInfo AccountShareInfor, award decimal.Decimal, awardArray map[string]string) map[string]string {
	for i := 0; i < len(accountShareInfo.AccoutShares); i++ {
		accountShare := accountShareInfo.AccoutShares[i]
		accountAward := accountShare.ShareNum.Div(accountShareInfo.TotalNum).Mul(award)
		awardArray[accountShare.Account] = accountAward.StringFixedBank(core.CoinPlaces)
	}
	return awardArray
}


func caculateCopyrightSize(copyrightSize int64, offerAccountShare string) (map[string]int64, error) {
	accountShares := strings.Split(offerAccountShare, "-") 
	copyrightSizeDecimal := decimal.New(copyrightSize, 0)
	zeroDecimal := decimal.NewFromInt(0)
	perShareDecimal := copyrightSizeDecimal.Div(decimal.RequireFromString("1000")) 
	var accounShareDecimalMap = map[string]int64{}
	for i := 0; i < len(accountShares); i++ {
		accountShare := strings.Split(accountShares[i], ":")       
		shareDecimal := decimal.RequireFromString(accountShare[1]) 
		if shareDecimal.LessThanOrEqual(zeroDecimal) {             
			return accounShareDecimalMap, types.ErrInvalidContributionValue
		}
		accountShareDecimal := perShareDecimal.Mul(shareDecimal)
		accounShareDecimalMap[accountShare[0]] = accountShareDecimal.IntPart()
	}
	return accounShareDecimalMap, nil
}


type AccountShareInfor struct {
	TotalNum     decimal.Decimal 
	AccoutShares []AccountShare
}

type AccountShare struct {
	Account  string          
	ShareNum decimal.Decimal 
}


func anylseAccount(ctx sdk.Context, keeper Keeper, offerAccountShare string) (normalAccount AccountShareInfor, authorizeAccount AccountShareInfor) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	accountShares := strings.Split(offerAccountShare, "-")
	normalAccount = AccountShareInfor{}    
	authorizeAccount = AccountShareInfor{} 
	size := len(accountShares)
	for i := 0; i < size; i++ {
		log.WithField("Number of executions", i).Debug("anylseAccount")
		accountShare := strings.Split(accountShares[i], ":")
		shareDecimal := decimal.RequireFromString(accountShare[1])
		normalAccount.TotalNum = normalAccount.TotalNum.Add(shareDecimal)
		accountShareObj := AccountShare{}
		accountShareObj.Account = accountShare[0]
		accountShareObj.ShareNum = shareDecimal
		normalAccount.AccoutShares = append(normalAccount.AccoutShares, accountShareObj)
	}
	return
}

func (k Keeper) judgeMortgRedeemTimeDelay(ctx sdk.Context) float64 {
	realCoin := k.GetMiningAmount(ctx) 
	mortgage := realCoin.AmountDec()
	
	mortgageTimes := mortgage.Div(core.MortgageMinerStand)
	delayTime := mortgageTimes.IntPart() * core.MortgageMinerAddTimeStand
	return float64(delayTime)
}


func (k Keeper) CopyrightMortgRedeemNew(ctx sdk.Context) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	blockHeight := ctx.BlockHeight()
	mortgCopyrightInfor := k.GetMortgageCopyrightInfor(ctx)
	if mortgCopyrightInfor.TriggerTime == 0 {
		if mortgCopyrightInfor.LastTime != 0 {
			k.DeleteMortgageCopyrightInfor(ctx)
		}
		return
	}
	delayTimes := k.judgeMortgRedeemTimeDelay(ctx)
	flag := true        
	editorFlag := false 
	for {
		lastRedeemTime := util.TimeStampToTime(mortgCopyrightInfor.TriggerTime)
		secondDiff := ctx.BlockHeader().Time.Sub(lastRedeemTime).Seconds()
		log.WithField("time difference", secondDiff).Debug("CopyrightMortgRedeemNew")
		if mortgCopyrightInfor.TriggerTime == 0 {
			break
		}
		
		if secondDiff >= (core.CopyrightMortgRedeemTimePerioad + delayTimes) { 
			timeMortgage := k.GetMortgageTimeInfor(ctx, mortgCopyrightInfor.TriggerTime)

			for i := 0; i < len(timeMortgage.MortgageAccount); i++ {
				mortgAccount := timeMortgage.MortgageAccount[i]
				
				ledgeCoin := types.MustRealCoin2LedgerCoin(mortgAccount.MortgageAmount)
				coins := sdk.Coins{ledgeCoin}
				err := k.CoinKeeper.SendCoins(ctx, mortgAccount.TxHashAccount, mortgAccount.MortageAccount, coins)
				if err != nil {
					log.WithError(err).WithFields(logrus.Fields{
						"fromAddr": mortgAccount.TxHashAccount.String(),
						"toAddr":   mortgAccount.MortageAccount.String(),
						"amt":      coins.String(),
					}).Error("SendCoins")
					panic(err)
				}
				
				mortgRedeemTradeInfor := types.TradeInfor{
					From:           mortgAccount.TxHashAccount.String(),
					To:             mortgAccount.MortageAccount.String(),
					TradeType:      core.TradeTypeCopyrightBuyRedeem,
					Amount:         mortgAccount.MortgageAmount,
					BlockNum:       blockHeight,
					TimeStemp:      ctx.BlockTime().Unix(),
					FromFsvBalance: k.GetBalance(ctx, core.MainToken, mortgAccount.TxHashAccount),
					FromTipBalance: k.GetBalance(ctx, core.InviteToken, mortgAccount.TxHashAccount),
					ToFsvBalance:   k.GetBalance(ctx, core.MainToken, mortgAccount.MortageAccount),
					ToTipBalance:   k.GetBalance(ctx, core.InviteToken, mortgAccount.MortageAccount),
				}
				err = k.AddBlockRDS(ctx, types.NewBlockRD(mortgRedeemTradeInfor))
				if err != nil {
					log.WithError(err).Error("AddBlockRDS-mortgRedeemTradeInfor")
				}
				mortgRedeemInfor := types.MortgageRedeem{
					MortageAccount: mortgAccount.TxHashAccount.String(),
					MortgageAmount: mortgAccount.MortgageAmount,
					RedeemAccount:  mortgAccount.MortageAccount.String(),
					RedeemAmount:   mortgAccount.MortgageAmount,
					BlockNum:       blockHeight,
					TimeStemp:      ctx.BlockTime().Unix(),
					TxHash:         mortgAccount.Txhash,
				}
				err = k.AddBlockRDS(ctx, types.NewBlockRD(mortgRedeemInfor))
				if err != nil {
					log.WithError(err).Error("AddBlockRDS-mortgRedeemInfor")
				}
				editorFlag = true
			}
			
			k.DeleteMortgageTimeInfor(ctx, mortgCopyrightInfor.TriggerTime)
			
			mortgCopyrightInfor.TriggerTime = timeMortgage.NextTime
		} else {
			flag = false
			break
		}
		if !flag {
			break
		}
	}
	
	if editorFlag {
		err := k.SetMortgageCopyrightInfor(ctx, mortgCopyrightInfor)
		if err != nil {
			log.WithError(err).Error("SetMortgageCopyrightInfor")
			return
		}
	} else { 
		if mortgCopyrightInfor.LastTime == 0 {
			k.DeleteMortgageTimeInfor(ctx, mortgCopyrightInfor.TriggerTime)
		}
	}

	/*for i := 0; i < len(mortgAccountArray); i++ {
		mortgAccount := mortgAccountArray[i]
		lastRedeemTime := util.TimeStampToTime(mortgAccount.MortgTime)
		//hourDiff := ctx.BlockHeader().Time.Sub(lastRedeemTime).Hours()
		secondDiff := ctx.BlockHeader().Time.Sub(lastRedeemTime).Seconds()
		log.Info("%s", secondDiff)
		
		if secondDiff >= (config.CopyrightMortgRedeemTimePerioad + delayTimes) { 
			if mortgAccount.TxHashAccount.String() == "" {
				continue
			}
			
			ledgeCoin := types.MustRealCoin2LedgerCoin(mortgAccount.MortgageAmount)
			coins := sdk.Coins{ledgeCoin}
			err := k.CoinKeeper.SendCoins(ctx, mortgAccount.TxHashAccount, mortgAccount.MortageAccount, coins)
			if err != nil {
				log.Error("", err, mortgAccount)
				continue
			}
			
			mortgRedeemTradeInfor := types.TradeInfor{
				From:           mortgAccount.TxHashAccount.String(),
				To:             mortgAccount.MortageAccount.String(),
				TradeType:      types.TradeTypeCopyrightBuyRedeem,
				Amount:         mortgAccount.MortgageAmount,
				BlockNum:       blockHeight,
				TimeStemp:      ctx.BlockTime().Unix(),
				FromFsvBalance: k.GetBalance(ctx, config.MainToken, mortgAccount.TxHashAccount),
				FromTipBalance: k.GetBalance(ctx, config.InviteToken, mortgAccount.TxHashAccount),
				ToFsvBalance:   k.GetBalance(ctx, config.MainToken, mortgAccount.MortageAccount),
				ToTipBalance:   k.GetBalance(ctx, config.InviteToken, mortgAccount.MortageAccount),
			}
			k.AddBlockRDS(ctx, types.NewBlockRD(mortgRedeemTradeInfor))
			mortgRedeemInfor := types.MortgageRedeem{
				MortageAccount: mortgAccount.TxHashAccount.String(),
				MortgageAmount: mortgAccount.MortgageAmount,
				RedeemAccount:  mortgAccount.MortageAccount.String(),
				RedeemAmount:   mortgAccount.MortgageAmount,
				BlockNum:       blockHeight,
				TimeStemp:      ctx.BlockTime().Unix(),
				TxHash:         mortgAccount.Txhash,
			}
			k.AddBlockRDS(ctx, types.NewBlockRD(mortgRedeemInfor))
			
			deleteAccount = append(deleteAccount, mortgAccount.Txhash)
		}
	}
	
	if len(deleteAccount) > 0 {
		log.Error("***************************", deleteAccount)
		log.Error("***************************", mortgAccountArray)
		for i := 0; i < len(deleteAccount); i++ {
			txhash := deleteAccount[i]
			for j := 0; j < len(mortgAccountArray); j++ {
				if mortgAccountArray[j].Txhash == txhash {
					mortgAccountArray = append(mortgAccountArray[0:j], mortgAccountArray[j+1:]...)
					break
				}
			}
		}
		log.Error("***************************", mortgAccountArray)
		mortgAccountInfor.MortgageAccountArray = mortgAccountArray
		k.SetMortgageCopyrightInfor(ctx, mortgAccountInfor)
	}*/
}
