package keeper

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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



type AccountSpaceMiner struct {
	Account         string               `json:"account"`
	SpaceTotal      decimal.Decimal      `json:"space_total"`
	UsedSpace       decimal.Decimal      `json:"used_space"`
	DeflationAmount types.RealCoin       `json:"deflation_amount"`
	BuySpace        decimal.Decimal      `json:"buy_space"`
	RewardSpace     decimal.Decimal      `json:"reward_space"`
	Settlement      map[int64]Settlement `json:"settlement"`
	SettlementEnd   Settlement           `json:"settlement_end"`
	LockedSpace     decimal.Decimal      `json:"locked_space"`
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
func (k Keeper) QueryAccountMinerBonusAmount(ctx sdk.Context, account string) (string, AccountSpaceMiner, error) {
	accountSpaceMiner := k.QueryAccountSpaceMinerInfor(ctx, account)
	if (accountSpaceMiner.Settlement == nil || len(accountSpaceMiner.Settlement) == 0) && accountSpaceMiner.SettlementEnd.Index == 0 {
		return "0", accountSpaceMiner, types.NoSpaceMinerRewardErr
	}
	height := ctx.BlockHeight()
	if height <= config.MinerStartHeight {
		return "0", accountSpaceMiner, types.NoSpaceMinerRewardErr
	}
	index := (height-config.MinerStartHeight)/config.SpaceMinerBonusBlockNum + 1
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
	var totalSpace decimal.Decimal
	var accountBonus decimal.Decimal
	var accountBonusString string
	bonusDecimal := k.QuerySpaceMinerBonusAmount(ctx)
	bonusDec := sdk.NewDec(bonusDecimal.IntPart())
	if !bonusDec.IsPositive() {
		return "0", accountSpaceMiner, nil
	}

	totalSpace = accountSpaceMiner.SettlementEnd.IndexSpace
	for _, index := range indexArray {
		spaceTotal := k.QuerySpaceTotalIndex(ctx, index)
		if spaceTotal.Sign() <= 0 {
			continue
		}
		perBonus := bonusDecimal.Div(spaceTotal)
		logs.Info("per bonus " + perBonus.String())
		logs.Info("times" + strconv.FormatInt(index, 10))
		if _, ok := accountSpaceMiner.Settlement[index]; ok {
			totalSpace = accountSpaceMiner.Settlement[index].IndexSpace
		}
		accountBonus = accountBonus.Add(totalSpace.Mul(perBonus))
	}
	accountBonusString = util.DecimalStringFixed(accountBonus.String(), config.CoinPlaces)
	accountSpaceMiner.SettlementEnd.Index = index
	accountSpaceMiner.SettlementEnd.IndexSpace = accountSpaceMiner.SpaceTotal
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
	store := k.KVHelper(ctx)
	bz := store.Get(spaceTotalIndexKey + strconv.Itoa(int(index)))
	if bz != nil {
		logs.Error("index", string(bz))
		return decimal.RequireFromString(string(bz))
	}
	return decimal.Decimal{}
}

func (k Keeper) SetSpaceTotalIndex(ctx sdk.Context) {
	store := k.KVHelper(ctx)
	spaceTotal := k.QueryDeflatinSpaceTotal(ctx)
	height := ctx.BlockHeight()
	index := (height - config.MinerStartHeight) / config.SpaceMinerBonusBlockNum
	store.Set(spaceTotalIndexKey+strconv.Itoa(int(index)), spaceTotal)
}

func (k Keeper) SpaceMinerRewardSettlement(ctx sdk.Context, account string) error {
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
		logrus.Error("format account error", err)
		return err
	}
	err = k.CoinKeeper.SendCoinsFromModuleToAccount(ctx, config.KeyCopyrighDeflation, accAddress, bonusCoins)
	if err != nil {
		return err
	}
	err = k.SetSpaceMinerBonusAmount(ctx, decimal.RequireFromString(accountBonusString))
	if err != nil {
		return err
	}
	feeRealCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, strconv.FormatFloat(config.CopyrightFee, 'f', 3, 64))
	var tradeRd []types.BlockRD
	tradeInfor := types.TradeInfor{
		From:           config.ContractAddressDeflation.String(),
		To:             account,
		Amount:         realCoinBonus,
		Fee:            feeRealCoin,
		BlockNum:       ctx.BlockHeight(),
		TradeType:      types.TradeTypeSpaceMinerBonus,
		TimeStemp:      ctx.BlockTime().Unix(),
		Txhash:         formatTxHash(ctx.TxBytes()),
		FromFsvBalance: k.GetBalance(ctx, config.MainToken, config.ContractAddressDeflation),
		FromTipBalance: k.GetBalance(ctx, config.InviteToken, config.ContractAddressDeflation),
		ToFsvBalance:   k.GetBalance(ctx, config.MainToken, accAddress),
		ToTipBalance:   k.GetBalance(ctx, config.InviteToken, accAddress),
	}
	tradeRd = append(tradeRd, types.NewBlockRD(tradeInfor))
	err = k.AddBlockRDS(ctx, tradeRd...)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) AddSpaceMiner(ctx sdk.Context, spaceMiner types.SpaceMinerData) error {
	amount := spaceMiner.DeflationAmount.AmountDec()
	return k.InviteReward(ctx, amount.Round(4), spaceMiner.AwardAccount, 1)
}


func (k Keeper) SetAccountSpaceMinerInfor(ctx sdk.Context, accountSpace AccountSpaceMiner) {
	store := k.KVHelper(ctx)
	key := spaceMinerKey + accountSpace.Account
	store.Set(key, accountSpace)
}


func (k Keeper) QueryDeflatinSpaceTotal(ctx sdk.Context) string {
	store := k.KVHelper(ctx)
	bz := store.Get(deflationSpaceTotalKey)
	logs.Info("size", string(bz))
	if bz != nil {
		return string(bz)
	} else {
		return "0"
	}
}


func (k Keeper) SetDeflationSpaceTotal(ctx sdk.Context, space decimal.Decimal) {
	store := k.KVHelper(ctx)
	var amountSpaceMiner decimal.Decimal
	if !store.Has(deflationSpaceTotalKey) {
		amountSpaceMiner = space
	} else {
		amountSpaceMiner = decimal.RequireFromString(string(store.Get(deflationSpaceTotalKey)))
		amountSpaceMiner = amountSpaceMiner.Add(space)
	}
	store.Set(deflationSpaceTotalKey, amountSpaceMiner.StringFixed(4))
}


func (k Keeper) QuerySpaceMinerAmount(ctx sdk.Context) decimal.Decimal {
	store := k.KVHelper(ctx)

	var amountSpaceMiner decimal.Decimal
	if !store.Has(spaceMinerAmountKey) {
		return amountSpaceMiner
	}
	err := store.GetUnmarshal(spaceMinerAmountKey, &amountSpaceMiner)
	if err != nil {
		logs.Error("unmarshal error", err)
	}
	return amountSpaceMiner
}


func (k Keeper) SetSpaceMinerAmount(ctx sdk.Context, amount decimal.Decimal) error {
	store := k.KVHelper(ctx)
	var amountSpaceMiner decimal.Decimal

	if store.Has(spaceMinerAmountKey) {
		err := store.GetUnmarshal(spaceMinerAmountKey, &amountSpaceMiner)
		if err != nil {
			logs.Error("save sapce miner error", err)
			return err
		}
		amountSpaceMiner = amountSpaceMiner.Add(amount)
	} else {
		amountSpaceMiner = amount
	}
	return store.Set(spaceMinerAmountKey, amountSpaceMiner)
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
	store := k.KVHelper(ctx)
	if store.Has(spaceMinerAccountKey) {
		return
	}
	accountString := account
	store.Set(spaceMinerAccountKey, accountString)
}



func (k Keeper) QueryAccountSpaceMinerInfor(ctx sdk.Context, address string) AccountSpaceMiner {
	store := k.KVHelper(ctx)
	var accountSpaceMiner AccountSpaceMiner
	bz := store.Get(spaceMinerKey + address)
	if len(bz) > 0 {
		err := util.Json.Unmarshal(bz, &accountSpaceMiner)
		if err != nil {
			panic(err)
		}
	}
	return accountSpaceMiner
}


func (k Keeper) LockAccountSpaceMiner(ctx sdk.Context, address string, lockedSize int64) error {
	store := k.KVHelper(ctx)
	var accountSpaceMiner AccountSpaceMiner
	bz := store.Get(spaceMinerKey + address)
	if len(bz) > 0 {
		err := util.Json.Unmarshal(bz, &accountSpaceMiner)
		if err != nil {
			panic(err)
		}

		leftSpace := accountSpaceMiner.SpaceTotal.Sub(accountSpaceMiner.UsedSpace).Sub(accountSpaceMiner.LockedSpace)
		lockedDecimal := decimal.NewFromInt(lockedSize)
		if leftSpace.Sub(lockedDecimal).Sign() >= 0 {
			accountSpaceMiner.LockedSpace = accountSpaceMiner.LockedSpace.Add(lockedDecimal)
			err = store.Set(spaceMinerKey+address, accountSpaceMiner)
			if err != nil {
				return err
			}
		}
	} else {
		return types.SpaceFirstError
	}
	return nil
}


func (k Keeper) LockAccountSpaceMinerReturn(ctx sdk.Context, address string, lockedSize int64) error {
	store := k.KVHelper(ctx)
	var accountSpaceMiner AccountSpaceMiner
	bz := store.Get(spaceMinerKey + address)
	if len(bz) > 0 {
		err := util.Json.Unmarshal(bz, &accountSpaceMiner)
		if err != nil {
			panic(err)
		}
		lockedDecimal := decimal.NewFromInt(lockedSize)
		accountSpaceMiner.LockedSpace = accountSpaceMiner.LockedSpace.Sub(lockedDecimal)
		err = store.Set(spaceMinerKey+address, accountSpaceMiner)
		if err != nil {
			return err
		}
	} else {
		return types.AccountSpaceReturnError
	}
	return nil
}




func (k Keeper) SetDeflationRateVoteIndex(ctx sdk.Context, index string) {
	store := k.KVHelper(ctx)
	store.Set(deflation_rate_index_key, index)
}


func (k Keeper) UpdateDeflationMinerInforByVote(ctx sdk.Context, minerAmount decimal.Decimal) {
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
		logs.Error("save deflation miner error", err)
	}
}


func (k Keeper) UpdateDeflationMinerInfor(ctx sdk.Context, minerAmount decimal.Decimal, blockNum int64) {
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
		logs.Error("保存通缩挖矿信息失败", err)
	}
}

func (k Keeper) QueryDeflationMinerInfor(ctx sdk.Context) (DeflationMinerInfor, error) {
	store := k.KVHelper(ctx)
	var deflationMiner DeflationMinerInfor
	if !store.Has(deflationMinerKey) {
		deflationMiner = DeflationMinerInfor{
			DeflationStatus:   0,
			DayMinerRemain:    config.DayMinerRemain,
			DayMinerAmount:    config.SpaceMinerPerDayStandand.Add(config.ValidatorMinerPerDayStandand),
			MinerTotalAmount:  config.MinerUpperLimitStandand,
			HasMinerAmount:    config.ChuangshiFee,
			RemainMinerAmount: config.MinerUpperLimitStandand.Sub(config.ChuangshiFee),
			MinerBlockNum:     config.MinerStartHeight,
		}
		if ctx.BlockHeight() > config.MinerStartHeight {
			deflationMiner.DeflationStatus = 1
		}
		return deflationMiner, nil
	}
	err := store.GetUnmarshal(deflationMinerKey, &deflationMiner)
	if err != nil {
		logs.Error("unmarshal error", err)
		return DeflationMinerInfor{}, err
	}
	return deflationMiner, nil
}


func (k Keeper) QuerySpaceMinerBonusAmount(ctx sdk.Context) decimal.Decimal {
	store := k.KVHelper(ctx)
	if !store.Has(spaceMinerBonusKey) {
		return config.SpaceMinerPerDayStandand
	} else {
		var hasBonusAmount decimal.Decimal
		err := store.GetUnmarshal(spaceMinerBonusKey, &hasBonusAmount)
		if err != nil {
			panic(err)
		}
		if config.MinerUpperLimitStandand.LessThanOrEqual(hasBonusAmount.Add(config.SpaceMinerPerDayStandand)) {
			if config.MinerUpperLimitStandand.GreaterThan(hasBonusAmount) {
				currentBonus := config.MinerUpperLimitStandand.Sub(hasBonusAmount)
				return currentBonus
			} else {
				return decimal.Decimal{}
			}
		} else {
			return config.SpaceMinerPerDayStandand
		}
	}
}


func (k Keeper) QueryValidatorMinerBonusAmount(ctx sdk.Context) decimal.Decimal {
	store := k.KVHelper(ctx)
	if !store.Has(spaceMinerBonusKey) {
		return config.ValidatorMinerPerDayStandand
	}
	var hasBonusAmount decimal.Decimal
	err := store.GetUnmarshal(spaceMinerBonusKey, &hasBonusAmount)
	if err != nil {
		panic(err)
	}

	if config.MinerUpperLimitStandand.LessThanOrEqual(hasBonusAmount.Add(config.ValidatorMinerPerDayStandand)) {
		if config.MinerUpperLimitStandand.GreaterThan(hasBonusAmount) {
			currentBonus := config.MinerUpperLimitStandand.Sub(hasBonusAmount)
			return currentBonus
		} else {
			return decimal.Decimal{}
		}
	} else {
		return config.ValidatorMinerPerDayStandand
	}
}


func (k Keeper) SetSpaceMinerBonusAmount(ctx sdk.Context, amount decimal.Decimal) error {
	store := k.KVHelper(ctx)
	bz := store.Get(spaceMinerBonusKey)
	var amountSpaceMiner decimal.Decimal
	if bz != nil {
		err := util.Json.Unmarshal(bz, &amountSpaceMiner)
		if err != nil {
			return err
		}
		amountSpaceMiner = amountSpaceMiner.Add(amount)
	} else {
		amountSpaceMiner = amount
	}
	store.Set(spaceMinerBonusKey, amountSpaceMiner)
	return nil
}




func (k Keeper) ValidatorBonus(ctx sdk.Context) decimal.Decimal {

	bonusDecimal := k.QueryValidatorMinerBonusAmount(ctx)

	if bonusDecimal.Sign() > 0 {
		realCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, bonusDecimal.String())
		bonusLedgerCoin := types.MustRealCoin2LedgerCoin(realCoin)
		err := k.CoinKeeper.SendCoinsFromModuleToModule(ctx, config.KeyCopyrighDeflation, authtypes.FeeCollectorName, sdk.NewCoins(bonusLedgerCoin))
		if err != nil {
			logs.Error("SendCoinsFromModuleToModule error", err)
			panic(err)
		}

		err = k.SetSpaceMinerBonusAmount(ctx, bonusDecimal)
		if err != nil {
			return decimal.Zero
		}
		feeRealCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, "0")

		tradeInfor := types.TradeInfor{
			From:           config.ContractAddressDeflation.String(),
			To:             config.ContractAddressFee.String(),
			Amount:         realCoin,
			Fee:            feeRealCoin,
			BlockNum:       ctx.BlockHeight(),
			TradeType:      types.TradeTypeValidatorMinerBonus,
			TimeStemp:      ctx.BlockTime().Unix(),
			FromFsvBalance: k.GetBalance(ctx, config.MainToken, config.ContractAddressDeflation),
			FromTipBalance: k.GetBalance(ctx, config.InviteToken, config.ContractAddressDeflation),
			ToFsvBalance:   k.GetBalance(ctx, config.MainToken, config.ContractAddressFee),
			ToTipBalance:   k.GetBalance(ctx, config.InviteToken, config.ContractAddressFee),
		}
		k.AddBlockRDS(ctx, types.NewBlockRD(tradeInfor))
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
	leftRealCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, util.DecimalStringFixed(leftBonus.String(), config.CoinPlaces))
	ledgeCoin := types.MustRealCoin2LedgerCoin(leftRealCoin)
	ledgeDecCoin := sdk.NewDecCoinFromDec(ledgeCoin.Denom, ledgeCoin.Amount.ToDec())
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
	totalSpaceStr := k.QueryDeflatinSpaceTotal(ctx)
	totalSpace, err := decimal.NewFromString(totalSpaceStr)
	if err != nil {
		logs.Error("SpaceFeeEstimate error | ", err.Error())
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
