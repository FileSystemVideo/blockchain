package keeper

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

const (
	deflation_rate_key       = "deflation_rate_"
	deflation_rate           = "1"
	deflation_rate_index_key = "deflation_rate_index_key"
	deflation_rate_index     = "1"
	deflation_vote_key       = "deflation_vote_key_"
	deflation_vote_title_key = "deflation_vote_title_key_"
	deflation_rate_min           = "0.000001"
)

var rateVoteInfor = map[string]string{
	"30": "30",
	"28": "28",
	"26": "26",
	"24": "24",
	"22": "22",
	"20": "20",
	"18": "18",
	"16": "16",
	"14": "14",
	"12": "12",
	"10": "10",
}

type DeflationVoteAccount struct {
	VoteAccount sdk.AccAddress `json:"vote_account"`
	VoteOption  string         `json:"vote_option"`
}

type DeflationVoteAccountArray struct {
	DeflationVoteAccountArray []DeflationVoteAccount `json:"deflation_vote_account_list"`
}

func (k Keeper) DeflationVote(ctx sdk.Context, deflationVote types.DeflationVoteData) error {
	flag := k.CheckVoteOPtion(deflationVote.Option)
	if !flag {
		return types.DeflationVoteOptionErr
	}
	_, deflationVoteInfor := k.QueryDeflationVoteInfor(ctx)
	//判断是否已经投票
	if len(deflationVoteInfor.DeflationVoteAccountArray) > 0 {
		for i := 0; i < len(deflationVoteInfor.DeflationVoteAccountArray); i++ {
			deflationVoteAccount := deflationVoteInfor.DeflationVoteAccountArray[i]
			if deflationVoteAccount.VoteAccount.Equals(deflationVote.Creator) {
				return types.DeflationVoted
			}
		}
	}
	deflationVoteAccount := DeflationVoteAccount{}
	deflationVoteAccount.VoteAccount = deflationVote.Creator
	deflationVoteAccount.VoteOption = deflationVote.Option
	k.SaveDeflationVoteInfor(ctx, deflationVoteAccount)
	txhash := formatTxHash(ctx.TxBytes())
	deflationIndex := k.QueryDeflationRateVoteIndex(ctx)
	txDeflationVote := types.DeflationVote{
		Txhash:    txhash,
		VoteIndex: deflationIndex,
	}
	k.AddBlockRDS(ctx, types.NewBlockRD(txDeflationVote))
	return nil
}

func queryDeflationVoteInfor(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	_, voteArray := keeper.QueryDeflationVoteInfor(ctx)
	return util.Json.Marshal(voteArray)
}

func (k Keeper) QueryDeflationVoteInfor(ctx sdk.Context) (string, DeflationVoteAccountArray) {
	deflationIndex := k.QueryDeflationRateVoteIndex(ctx)
	store := ctx.KVStore(k.storeKey)
	deflation_vote_day_key := deflation_vote_key + deflationIndex
	bz := store.Get([]byte(deflation_vote_day_key))
	var deflationVoteArray DeflationVoteAccountArray
	if len(bz) > 0 {
		err := util.Json.Unmarshal(bz, &deflationVoteArray)
		if err != nil {
			logs.Error("unmarshal error", err)
		}
	}
	return deflation_vote_day_key, deflationVoteArray
}

func (k Keeper) SaveDeflationVoteInfor(ctx sdk.Context, voteAccount DeflationVoteAccount) {
	deflation_vote_day_key, deflationVoteArray := k.QueryDeflationVoteInfor(ctx)
	store := ctx.KVStore(k.storeKey)
	deflationVoteArray.DeflationVoteAccountArray = append(deflationVoteArray.DeflationVoteAccountArray, voteAccount)
	deflationVoteBytes, err := util.Json.Marshal(deflationVoteArray)
	if err != nil {
		ctx.Logger().Error("unmarshal error", err)
	} else {
		store.Set([]byte(deflation_vote_day_key), deflationVoteBytes)
	}
}

func (k Keeper) CheckVoteOPtion(option string) bool {
	if _, ok := rateVoteInfor[option]; ok {
		return true
	} else {
		return false
	}
}


func (k Keeper) QueryRateVoteTitle(ctx sdk.Context) map[string]string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(deflation_vote_title_key))
	if bz == nil {
		return rateVoteInfor
	} else {
		var voteTitleInfor map[string]string
		err := util.Json.Unmarshal(bz, &voteTitleInfor)
		if err != nil {
			return nil
		}
		return voteTitleInfor
	}
}


func (k Keeper) QueryDeflationRate(ctx sdk.Context) string {
	/*store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(deflation_rate_key))
	if bz == nil {
		return deflation_rate
	} else {
		return string(bz)
	}*/
	spaceTotal := k.QueryDeflatinSpaceTotal(ctx)
	//if ctx.BlockHeader().Height > config.MinerStartHeight && spaceTotal != "0" {
	if spaceTotal != "0" {
		return k.calculatDeflationRate(ctx)
	} else {
		return deflation_rate
	}
}

func (k Keeper) calculatDeflationRate(ctx sdk.Context) string {
	bonusDecimal := k.QuerySpaceMinerBonusAmount(ctx)
	spaceTotal := k.QueryDeflatinSpaceTotal(ctx)
	spaceTotalDecimal := decimal.RequireFromString(spaceTotal)
	totalSpaceDecimalM := spaceTotalDecimal.Div(ByteToMb)
	preMbonus := bonusDecimal.Div(totalSpaceDecimalM)
	currentRate := preMbonus.Mul(returnDays)
	if currentRate.GreaterThan(decimal.RequireFromString(deflation_rate_min)){
		return util.DecimalStringFixed(currentRate.String(), config.CoinPlaces)
	}else{
		return deflation_rate_min
	}

}


func (k Keeper) SetDeflationRate(ctx sdk.Context, rate string) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(deflation_rate_key), []byte(rate))
}


func (k Keeper) QueryDeflationRateVoteIndex(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(deflation_rate_index_key))
	if bz == nil {
		return deflation_rate_index
	} else {
		return string(bz)
	}
}

type DeflationVotePer struct {
	VoteShare        string `json:"vote_share"`
	VoteSharePercent string `json:"vote_share_percent"`
}

func (k Keeper) DeflationVoteDeal(ctx sdk.Context) {

	deflationIndex, deflationVoteInfor := k.QueryDeflationVoteInfor(ctx)
	deflationIndexArray := strings.Split(deflationIndex, "_")
	deflationRate := k.QueryDeflationRate(ctx)
	deflationRateVoteInfor := types.DeflationRateVoteInfor{}
	deflationRateVoteInfor.DeflationPre = deflationRate
	voteHeader := k.QueryRateVoteTitle(ctx)
	voteHeaderByte, _ := util.Json.Marshal(voteHeader)
	deflationRateVoteInfor.VoteTitle = string(voteHeaderByte)
	if len(deflationVoteInfor.DeflationVoteAccountArray) > 0 {
		//totalVoteShare := k.stakingKeeper.GetTotalDelegate(ctx)
		//totalVoteShare := sdk.NewDec(100000)
		realTotalShares := k.GetAllDelegatorShares(ctx)
		realTotalSharesDecimal := decimal.RequireFromString(realTotalShares)
		accountVoteMap := make(map[string]decimal.Decimal, 0)
		accountVoteStringMap := make(map[string]DeflationVotePer, 0)
		for i := 0; i < len(deflationVoteInfor.DeflationVoteAccountArray); i++ {
			voteAccount := deflationVoteInfor.DeflationVoteAccountArray[i]
			//accountShare := k.stakingKeeper.GetTotalDelegateFromAddr(ctx, voteAccount.VoteAccount)
			//accountShare := sdk.NewDec(52000)
			accountShare, _ := k.GetAccountDelegatorShares(ctx, voteAccount.VoteAccount)
			accountShareDecimal := decimal.RequireFromString(accountShare)
			if _, ok := accountVoteMap[voteAccount.VoteOption]; ok {
				accountVoteMap[voteAccount.VoteOption] = accountVoteMap[voteAccount.VoteOption].Add(accountShareDecimal)
				//accountVoteStringMap[voteAccount.VoteOption] = accountVoteMap[voteAccount.VoteOption].Add(accountShare).String()
				deflationVoteper := buildDeflationVotePer(accountVoteMap[voteAccount.VoteOption].Add(accountShareDecimal), realTotalSharesDecimal)
				accountVoteStringMap[voteAccount.VoteOption] = deflationVoteper
				deflationAccountVote := types.DeflationAccountVote{
					Account:   voteAccount.VoteAccount.String(),
					Power:     accountShare,
					VoteIndex: deflationIndexArray[len(deflationIndexArray)-1],
				}
				k.AddBlockRDS(ctx, types.NewBlockRD(deflationAccountVote))
				//blockchain.UpdatePower(deflationVoteper.VoteShare, voteAccount.VoteAccount.String(), deflationIndexArray[len(deflationIndexArray)-1])
			} else {
				if !accountShareDecimal.IsZero() {
					accountVoteMap[voteAccount.VoteOption] = accountShareDecimal
					deflationVoteper := buildDeflationVotePer(accountShareDecimal, realTotalSharesDecimal)
					accountVoteStringMap[voteAccount.VoteOption] = deflationVoteper
				}
				deflationAccountVote := types.DeflationAccountVote{
					Account:   voteAccount.VoteAccount.String(),
					Power:     accountShare,
					VoteIndex: deflationIndexArray[len(deflationIndexArray)-1],
				}
				k.AddBlockRDS(ctx, types.NewBlockRD(deflationAccountVote))
				//blockchain.UpdatePower(accountShare.FormatStandDecToString6(), voteAccount.VoteAccount.String(), deflationIndexArray[len(deflationIndexArray)-1])
			}
		}
		dd := realTotalSharesDecimal.Mul(decimal.RequireFromString("0.51"))
		for option, voteNum := range accountVoteMap {
			if voteNum.GreaterThanOrEqual(dd) {
				//rateDecimal := decimal.RequireFromString(deflationRate)
				//optionDecimal := decimal.RequireFromString(option)
				//rateDecimal = rateDecimal.Add(rateDecimal.Mul(optionDecimal))
				k.SetDeflationRate(ctx, option)
				newVoteTitle := calculationOption(option)
				k.SetDeflationVoteTitle(ctx, newVoteTitle)
				deflationRateVoteInfor.Deflation = option
				deflationRateVoteInfor.Status = 1
				break
			}
		}
		accountVoteByte, _ := util.Json.Marshal(accountVoteStringMap)
		deflationRateVoteInfor.VoteOption = string(accountVoteByte)
	} else {
		deflationRateVoteInfor.Deflation = deflationRate
	}


	indexInt, err := strconv.Atoi(deflationIndexArray[len(deflationIndexArray)-1])
	if err != nil {
		panic(err)
	}
	deflationRateVoteInfor.VoteIndex = deflationIndexArray[len(deflationIndexArray)-1]
	deflationRateVoteInfor.EndNum = ctx.BlockHeight()
	deflationRateVoteInfor.BeginNum = ctx.BlockHeight() - config.DeflationVoteDealBlockNum
	indexInt += 1
	k.SetDeflationRateVoteIndex(ctx, strconv.Itoa(indexInt))
	k.AddBlockRDS(ctx, types.NewBlockRD(deflationRateVoteInfor))
}

var percentDecimal = decimal.RequireFromString("100")

func buildDeflationVotePer(accountShare, totalVoteShare decimal.Decimal) DeflationVotePer {
	var deflationVoteper = DeflationVotePer{}
	deflationVoteper.VoteShare = accountShare.String()
	deflationVoteper.VoteSharePercent = accountShare.Div(totalVoteShare).Mul(percentDecimal).StringFixed(2)
	return deflationVoteper
}

func calculationOption(options string) map[string]string {
	option := decimal.RequireFromString(options)
	min := option.Sub(option.Mul(decimal.NewFromFloat(0.5)))
	max := option.Mul(decimal.NewFromFloat(0.5)).Add(option)
	cha := max.Sub(min)
	increment := cha.Div(decimal.NewFromInt(10))
	optionMap := make(map[string]string)
	optionMap[min.String()] = min.String()
	op := min
	for i := 0; i < 10; i++ {
		op = op.Add(increment)
		optionMap[op.String()] = op.String()
	}
	return optionMap
}


func (k Keeper) SetDeflationVoteTitle(ctx sdk.Context, voteTitle map[string]string) {
	store := ctx.KVStore(k.storeKey)
	voteTitleBytes, _ := util.Json.Marshal(voteTitle)
	store.Set([]byte(deflation_vote_title_key), voteTitleBytes)
}
