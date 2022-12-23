package core

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

const (
	
	AppName = "vcd"

	CommandName = "videochain"

	
	Version = "22.12.10"

	EvmRpcURL = "http://localhost:8545"

	RpcURL = "tcp://127.0.0.1:26657"

	ServerURL = "http://127.0.0.1:1317"

	
	MainToken = "fsv"

	
	InviteToken = "tip"

	BaseDenomUnit = 18

	//key
	KeyCopyrightBonus = "copyright_bonus"

	//key
	KeyCopyrightMortgage = "copyright_mortgage"

	//key
	KeyCopyrighDeflation = "copyright_deflation"

	//ChainDefaultFeeStr string = "0.001" 

	//ChainDefaultFee float64 = 0.001 //   fsv  ,  

	//MinFsvTransfer float64 = 0.001 
	
	//CopyrightInviteFee float64 = 0.0001 
	
	//TransferRate string = "0.002" 
	
	//MortgageFee = "0.002" 

	MortgageHeight = 14400 

	CoinPlaces = 18 

	DiskSpaceMortgRate = 1 

	
	MinLedgerAmountInt64 int64 = 1

	
	AuthorizeRate = "0.9"

	
	ValidatorRegisterMinAmount = 50

	
	SpaceMinerBonusBlockNum = 14400
	//SpaceMinerBonusBlockNum = 100

	
	DeflationVoteDealBlockNum = 250

	
	MinerStartHeight int64 = 100800 //30 * 24 * 60 * 10

	
	MortgageStartHeight int64 = 5256000 //365 * 24 * 60 * 10

	MortgageRate int64 = 100 

	CopyrightMortgRedeemTimePerioad float64 = 60 * 10 

	
	InitPublisherId = 100000

	
	CommitTime = 6

	
	BonusTypeFront = "front"
	
	BonusTypeRear = "rear"

	
	InitHeight = 5341742

	
	DefaultBlockHash = "0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	
	MinRealAmountDec, _ = sdk.NewDecFromStr("0.000000000000000001") //dec

	
	RealToLedgerRate, _ = sdk.NewDecFromStr("1000000000000000000")

	
	LedgerToRealRate, _ = sdk.NewDecFromStr("0.000000000000000001")

	//   fsv  ,  
	ChainDefaultFee = decimal.RequireFromString("0.001")

	CoryrightChargeRateDecimalInputlLength = 3 

	CopyrightPriceDecimalInputlLength = 6 

	TransferAmountDecimalInputLength = 18 

	MinFsvTransfer = decimal.RequireFromString("0.001") 

	CopyrightInviteFee = decimal.RequireFromString("0.0001") 

	TransferRate = decimal.RequireFromString("0.002") 

	MortgageFee = decimal.RequireFromString("0.002") 

	MinimumGasPrices = decimal.RequireFromString("0.00005") //gas

	ChainID = "vc_8889-1" // vc_8888-1   vc_8889-1

	PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(CoinPlaces), nil))
	
	//MinerUpperLimit = sdk.NewDec(int64(1000000000000000))
	MinerUpperLimitStandand = decimal.RequireFromString("2079000000")
	//0.02 1956*0.002
	ChuangshiFee = decimal.RequireFromString("3.912")
	//MinerUpperLimitStandandOrigin = decimal.RequireFromString("2079000000")

	ByteToMb = decimal.RequireFromString("1048576")

	
	ChargeRateLow, ChargeRateHigh = sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("0.8")
	
	SpaceMinerPerDayStandand = decimal.RequireFromString("30000")
	
	ValidatorMinerPerDayStandand = decimal.RequireFromString("6000")
	
	DayMinerRemain = MinerUpperLimitStandand.Div(SpaceMinerPerDayStandand.Add(ValidatorMinerPerDayStandand)).IntPart()

	VoteResultTimePerioad = time.Duration(24) * time.Hour 
	//VoteResultTimePerioad = time.Duration(10) * time.Minute 

	MortgageMinerStand = decimal.RequireFromString("100000") 

	MortgageMinerAddTimeStand = int64(6 * 15) 

	CopyrightVoteTimePerioad = time.Duration(24) * time.Hour 
	//CopyrightVoteTimePerioad = time.Duration(5) * time.Minute 

	CopyrightVoteRedeemTimePerioad = time.Duration(360) * time.Hour 
	//CopyrightVoteRedeemTimePerioad = time.Duration(15) * time.Minute 

	CopyrightVoteAwardRate   = decimal.RequireFromString("0.01")  
	CopyrightVoteAwardRateV2 = decimal.RequireFromString("0.005") 

	
	CrossChainOutMinAmount = "50000000"

	CrossChainOutFeeRatio = "0.01"
)
