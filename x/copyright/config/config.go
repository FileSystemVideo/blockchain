package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"time"
)

const (

	ChainID = "fsv20211021"

	RpcUrl = "tcp://127.0.0.1:26657"


	MainToken = "fsv"


	InviteToken = "tip"


	KeyCopyrightBonus = "copyright_bonus"


	KeyCopyrightMortgage = "copyright_mortgage"


	KeyCopyrighDeflation = "copyright_deflation"

	CopyrightFee float64 = 0.001

	MinFsvTransfer float64 = 0.001

	CopyrightInviteFee float64 = 0.0001

	TransferRate string = "0.002"

	MortgageFee = "0.002"

	MortgageHeight = 14400



	MinRealAmountFloat64 float64 = 0.000001

	CoinPlaces = 6

	DiskSpaceMortgRate = 1


	MinLedgerAmountInt64 int64 = 1


	RealToLedgerRate float64 = float64(RealToLedgerRateInt64)


	RealToLedgerRateInt64 int64 = 1000000


	LedgerToRealRate = "0.000001"


	InvestorBonusRate = "0.1"


	AuthorizeRate = "0.9"


	ValidatorRegisterMinAmount = 50


	SpaceMinerBonusBlockNum = 14400



	DeflationVoteDealBlockNum = 250


	MinerStartHeight int64 = 100800



	MortgageStartHeight int64 = 5256000

	MortgageRate int64 = 100

	CopyrightMortgRedeemTimePerioad float64 = 60 * 60 * 48


	InitPublisherId = 100000


	CommitTime = 6


	BonusTypeFront = "front"

	BonusTypeRear = "rear"

	MinimumGasPrices = "0.00005"
)

var (


	MinerUpperLimitStandand = decimal.RequireFromString("2079000000")

	ChuangshiFee = decimal.RequireFromString("3.912")


	ByteToMb = decimal.RequireFromString("1048576")


	ChargeRateLow, ChargeRateHigh = sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("0.8")

	SpaceMinerPerDayStandand = decimal.RequireFromString("30000")

	ValidatorMinerPerDayStandand = decimal.RequireFromString("6000")

	DayMinerRemain = MinerUpperLimitStandand.Div(SpaceMinerPerDayStandand.Add(ValidatorMinerPerDayStandand)).IntPart()

	VoteResultTimePerioad = time.Duration(24) * time.Hour


	MortgageMinerStand = decimal.RequireFromString("100000")

	MortgageMinerAddTimeStand = int64(6 * 15)

	CopyrightVoteTimePerioad = time.Duration(24) * time.Hour


	CopyrightVoteRedeemTimePerioad = time.Duration(360) * time.Hour


	CopyrightVoteAwardRate = decimal.RequireFromString("0.01")


	CrossChainOutMinAmount = "50000000"

	CrossChainOutFeeRatio = "0.01"
)
