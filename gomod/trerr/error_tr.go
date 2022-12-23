package trerr

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var Language = "EN"

type TrErr struct {
	RawMsg string
	TrMsg  string
	ErrMsg string
}

// 
func (tr *TrErr) Error() string {
	return tr.ErrMsg
}

//，
func (tr *TrErr) GetError() error {
	if Language == "CHC" {
		tr.ErrMsg = tr.TrMsg
		return tr
	} else {
		tr.ErrMsg = tr.RawMsg
		return tr
	}
}

// 
func (tr *TrErr) DumpError(errMsg string) {
	if errMsg != "" {
		r, _ := regexp.Compile(tr.RawMsg)
		params := r.FindStringSubmatch(errMsg)
		if len(params) <= 1 || !strings.Contains(tr.TrMsg, "%s") {
			return
		} else {
			tmp := make([]interface{}, 0)
			for _, param := range params[1:] {
				tmp = append(tmp, param)
			}
			tr.ErrMsg = fmt.Sprintf(tr.TrMsg, tmp...)
			return
		}
	}
	return
}

var trMap = make(map[string]TrErr)

// 
func NewErr(unTrMsg, trMsg string) TrErr {
	te := TrErr{
		RawMsg: unTrMsg,
		TrMsg:  trMsg,
		ErrMsg: trMsg,
	}
	trMap[unTrMsg] = te
	return te
}

func TransError(unTrMsg string) error {
	for regC, te := range trMap {
		r, _ := regexp.Compile(regC)
		if r.Match([]byte(unTrMsg)) && Language == "CHC" {
			te.DumpError(unTrMsg)
			e := te
			return &e
		}
	}
	return errors.New(unTrMsg)
}

var (
	InvalidFromAddress  = NewErr("^invalid from address: (.*)$", "")
	ParseError          = NewErr("failed to parse request", "")
	ChainIdError        = NewErr("chain-id required but not specified", "")
	FeeGasError         = NewErr("cannot provide both fees and gas prices", "")
	FeeGasInvalid       = NewErr("invalid fees or gas prices provided", "")
	PassGroupError      = NewErr("password is invalid", "")
	ChargeRateEmpty     = NewErr("chargeRate is empty", "")
	ChargeRateTooSmall  = NewErr("chargeRate is too small", "")
	ChargeRateTooHigh   = NewErr("chargeRate is too high", "")
	Published           = NewErr("current datahash has published", "")
	QueryBindError      = NewErr("query bind relationship error", "")
	HasNoBind           = NewErr("current account has not bind relationship", "")
	InvalidAddress      = NewErr("invalid address", "")
	InsufficientError   = NewErr("insufficient funds", "")
	MarshalError        = NewErr("failed to marshal JSON bytes", "")
	UnMarshalError      = NewErr("failed to unmarshal JSON bytes", "")
	DataHashEmpty       = NewErr("datahash is empty", "")
	OriginDataHashEmpty = NewErr("origindatahash is empty", "")
	DataHashNotEmpty    = NewErr("datahash cannot be empty", "")
	GasAdjustmentError  = NewErr("invalid gas adjustment", "")
	AccountNoExist      = NewErr("account not exist", "")
	AccountError        = NewErr("decoding Bech32 address failed: must provide an address", "")
	CoinError           = NewErr("coin can not be empty", "")
	IllegalError        = NewErr("Illegal account proportion", "")
	DeleteError         = NewErr("delete fee amount is not invalid", "")
	HasDeleted          = NewErr("current datahash has deleted", "")
	BalanceError        = NewErr("Insufficient account balance", "")
	AccountEmpty        = NewErr("account can not be empty", "")
	PageInvalid         = NewErr("page is invalid", "")
	PageSizeINvalid     = NewErr("pageSize is invalid", "")
	TokenEmpty          = NewErr("token is empty", "")
	QueryParamError     = NewErr("query param empty", "")
	NumError            = NewErr("query param num invalid", "")
	QueryExtError       = NewErr("query extype invalid", "")
	IdsError            = NewErr("query param idsString is empty", "")
	TxhashEmpty         = NewErr("query param txhash empty", "")
	TxhashNotExist      = NewErr("txhash not exist", "")
	DeleAddressError    = NewErr("must use own delegator address", "POS")
	QueryDeletorError   = NewErr("query delegators amount error", "POS")
	UnbondInfuffiError  = NewErr("unbond amount is not enough", "")
	DeleAddressEmpty    = NewErr("delegatorAddr can not be empty", "POS")
	_                   = NewErr("param author is empty", "")
	_                   = NewErr("author not empty", "")
	_                   = NewErr("param productor is empty", "")
	_                   = NewErr("param legalNumber is empty", "")
	_                   = NewErr("param legalTime is empty", "")
	_                   = NewErr("param complainInfor is empty", "")
	_                   = NewErr("param complainAccount is empty", "")
	_                   = NewErr("param accuseAccount is empty", "")
	_                   = NewErr("do not appeal to yourself", "")
	_                   = NewErr("current datahash is appealing", "")
	_                   = NewErr("complainType can not be empty", "")
	_                   = NewErr("complainType is invalid", "")
	_                   = NewErr("aram accuseInfor is empty", "")
	_                   = NewErr("param complainId is empty", "")
	_                   = NewErr("complain not exist", "")
	_                   = NewErr("current complain has response", "")
	_                   = NewErr("current account  has no right to response", "")
	_                   = NewErr("param voteStatus is empty", "")
	_                   = NewErr("param voteStatus is invalid", "")
	_                   = NewErr("current complain status invalid", "")
	_                   = NewErr("current complain response status invalid", "")
	_                   = NewErr("complain vote can not repeat", "")
	_                   = NewErr("current account has no vote right", "")
	_                   = NewErr("current account has no right to get complain result", "")
	_                   = NewErr("current complain has finished", "")
	_                   = NewErr("complain vote has not reached end time", "")
	_                   = NewErr("calculate vote amount error", "")
	_                   = NewErr("query copyright complain error", "")
	_                   = NewErr("complain account is invalid", "")
	_                   = NewErr("complain status does not allow query for ip", "")
	_                   = NewErr("format gas error", "")
	_                   = NewErr("invalid gas adjustment", "")
	_                   = NewErr("query copyright bonus status error", "")
	_                   = NewErr("datahash does not exist", "")
	_                   = NewErr("datahash has exist", "")
	_                   = NewErr("down copyright price error", "")
	_                   = NewErr("datahash has download", "")
	_                   = NewErr("query data error", "")
	_                   = NewErr("datahash creator error", "")
	_                   = NewErr("copyright bonus has demand", "")
	_                   = NewErr("account bind has exist", "")
	_                   = NewErr("current accout has no right to delete", "")
	_                   = NewErr("copyright has deleted", "")
	_                   = NewErr("copyright files is empty", "")
	_                   = NewErr("copyright complain account is invalid", "")
	_                   = NewErr("complainId does not exist", "")
	_                   = NewErr("complain status is invalid", "")
	_                   = NewErr("complain vote status is invalid", "")
	_                   = NewErr("current account exist block chain request ,please wait a minute", "，")
	_                   = NewErr("sign error", "")
	_                   = NewErr("get accountManage error", "")
	_                   = NewErr("Entropy length must be \\[128, 256\\] and a multiple of 32", "")
	_                   = NewErr("account only exist", "")
	_                   = NewErr("encoding bech32 failed", "")
	_                   = NewErr("password verification failed", "")
	_                   = NewErr("account not exist", "")
	_                   = NewErr("account key not exist", "")

	_ = NewErr("failed to decrypt private key", "")
	_ = NewErr("invalid mnemonic", "")
	_ = NewErr("height must be equal or greater than zero", "")
	_ = NewErr("empty delegator address", "")
	_ = NewErr("empty validator address", "POS")
	_ = NewErr("invalid delegation amount", "")
	_ = NewErr("invalid shares amount", "")
	_ = NewErr("validator does not exist", "POS")
	_ = NewErr("invalid coin denomination", "")
	_ = NewErr("delegate progress error", "")
	_ = NewErr("no validator distribution info", "")
	_ = NewErr("no delegation distribution info", "")
	_ = NewErr("module account (.*)$ does not exist", "")
	_ = NewErr("no validator commission to withdraw", "")
	_ = NewErr("signature verification failed; verify correct account sequence and chain-id", ",ID")
	_ = NewErr("database error", "")
	_ = NewErr("parse account error", "")
	_ = NewErr("parse coin error", "")
	_ = NewErr("parse string to number error", "")
	_ = NewErr("query chain infor errors", "")
	_ = NewErr("valid chain request error", "")
	_ = NewErr("parse json error", "")
	_ = NewErr("parse byte to struct error", "")
	_ = NewErr("format tx struct error", "")
	_ = NewErr("broadcast error", "")
	_ = NewErr("format string to int error", "")
	_ = NewErr("datahash format account", "")
	_ = NewErr("parse valitor error", "")
	_ = NewErr("parse time error", "")
	_ = NewErr("current account has payed for datahash", "")
	_ = NewErr("current account has vote for this complain", "")
	_ = NewErr("datahash not exist", "")
	_ = NewErr("sign error", "")
	_ = NewErr("operator address exist", "POS")
	_ = NewErr("pub key for validator exist", "")
	_ = NewErr("current validator not jail", "POS")
	_ = NewErr("has no right to oprate validator", "")
	_ = NewErr("validator description length error", "")
	_ = NewErr("validator min mortgage amount", "POS50")
	_ = NewErr("current account does not have right to delete datahash ", "")
	_ = NewErr("too many unbonding delegation entries for \\(delegator, validator\\) tuple", "OS,.")
	_ = NewErr("no delegation for \\(address, validator\\) tuple", "POS")
	_ = NewErr("parse validator address error", "POS")
	_ = NewErr("There is no reward to receive", "")
	_ = NewErr("delegation does not exist", "")
	_ = NewErr("query sensitive words error", "")
	_ = NewErr("save sensitive words error", "")
	_ = NewErr("sensitive status illegal", "")
	_ = NewErr("origindatahash is empty", "hash")
	_ = NewErr("origindatahash has exist", "hash")
	_ = NewErr("origindatahash not exist", "hash")
	_ = NewErr("contain sensitive words", "")
	_ = NewErr("fee can not be zero", "")
	_ = NewErr("fee is too less", "gas")
	_ = NewErr("fee can not empty", "")
	_ = NewErr("delegation coin less then min", "")
	_ = NewErr("unbonding delegation shares less then min", "")
	_ = NewErr("delegation reward coin less then min", "POS0.000001")
	_ = NewErr("not enough delegation shares", "")
	_ = NewErr("must use own validator address", "")
	_ = NewErr("private key is empty", "")
	_ = NewErr("^verification fail$", "")
	_ = NewErr("pubkey not exist", "")
	_ = NewErr("verification error", "")
	_ = NewErr("con address is invalid", "")
	_ = NewErr("con address can not empty", "")
	_ = NewErr("classify id is invalid", "id")
	_ = NewErr("dir name has exist", "")
	_ = NewErr("delegator does not contain delegation", "POS")
	_ = NewErr("The account to be unlocked must have a valid POS mortgage", "POS")
	_ = NewErr("The mortgage amount is less than the self-mortgage amount", "")
	_ = NewErr("dir type illegal", "")
	_ = NewErr("dir not exist", "")
	_ = NewErr("query classify list error", "")
	_ = NewErr("query classify copyright list error", "")
	_ = NewErr("query classify list error", "")
	_ = NewErr("copyright not exist", "")
	_ = NewErr("Illegal share account", "")
	_ = NewErr("current account does not have right to editor datahash ", "")
	_ = NewErr("current account does not have right to delete datahash ", "")
	_ = NewErr("lower min coin", "")
	_ = NewErr("cannot move to subdirectory", "")
	_ = NewErr("cannot move to self directory", "")
	_ = NewErr("mortg miner has finish", "")
	_ = NewErr("SendCoins error", "")
	_ = NewErr("current classify has exist the same name", "")
	_ = NewErr("classify not dir", "")
	_ = NewErr("bindinfor has exist", "")
	_ = NewErr("bindinfor not exist", "")
	_ = NewErr("vote index not empty", "id")
	_ = NewErr("has not reach mortgage height", "")
	_ = NewErr("miner space not enough", "")
	_ = NewErr("param inviteCode is empty", "")
	_ = NewErr("token id not empty", "nft tokenId ")
	_ = NewErr("has vote deflation", "")
	_ = NewErr("empty description: invalid request", "")
	_ = NewErr("signature verification failed, invalid chainid or account number", ",chainIdnumber")
	_ = NewErr("account serial number expired, the reason may be: node block behind or repeatedly sent messages", ",:.")
	_ = NewErr("copyright ID has been used, please register again", "id，")
	_ = NewErr("The copyright ID is empty", "id")
	_ = NewErr("Verifier information can only be changed once in 24 hours", "24")
	_ = NewErr("Please use the main currency", "")
	_ = NewErr("chain error", "")
	_ = NewErr("min self delegation cannot be zero", "")

	_                    = NewErr("commission must be positive", "")
	_                    = NewErr("commission cannot be more than 100%", "100%")
	_                    = NewErr("commission cannot be more than the max rate", "")
	_                    = NewErr("commission cannot be changed more than once in 24h", "24")
	_                    = NewErr("commission change rate must be positive", "")
	_                    = NewErr("commission change rate cannot be more than the max rate", "")
	_                    = NewErr("commission cannot be changed more than max change rate", "")
	_                    = NewErr("validator's self delegation must be greater than their minimum self delegation", "DPOS")
	_                    = NewErr("minimum self delegation must be a positive integer", "")
	_                    = NewErr("minimum self delegation cannot be decrease", "")
	_                    = NewErr("copyright vote power not enough", "")
	_                    = NewErr("Copyright audit voting has ended", "")
	_                    = NewErr("Copyright voting is invalid", "")
	_                    = NewErr("Minimum amount 50fsv", "50fsv")
	_                    = NewErr("invalid amount", "")
	_                    = NewErr("tip balance must greate than one", "tip 1")
	_                    = NewErr("current datahash has not approve", "")
	_                    = NewErr("The resource is not in the approval period", "")
	_                    = NewErr("is not allowed to receive funds", "")
	_                    = NewErr("Fee coin only supports fsv", "fsv")
	_                    = NewErr("There can only be one handling fee currency", "")
	_                    = NewErr("miner space not enough", "")
	_                    = NewErr("fee amount is invalid", "")
	_                    = NewErr("Please redeem the space first", "")
	_                    = NewErr("Inconsistent ownership of copyright", "")
	_                    = NewErr("Unauthorized account number", "")
	_                    = NewErr("The copyright has not expired, and the purchase is invalid", "，")
	_                    = NewErr("Reward space settlement failed", "")
	_                    = NewErr("There is no mining income to be claimed", "")
	_                    = NewErr("Storage space return failed", "")
	_                    = NewErr("Invalid download contribution value", "")
	_                    = NewErr("NFT data formatting failed", "NFT")
	_                    = NewErr("You do not have permission to operate on the current NFT", "NFT")
	_                    = NewErr("NFT does not exist", "NFT")
	_                    = NewErr("Already voted", "")
	_                    = NewErr("Illegal deflation voting option", "")
	_                    = NewErr("RDS information not found", "RDS")
	_                    = NewErr("Copyright information not found", "")
	_                    = NewErr("insufficient fee", "gas")
	HistoryBlockNotQuery = NewErr("Historical block cannot be queried currently", "")
	_                    = NewErr("there are currently frozen ballots that cannot be redeemed", "")
	_                    = NewErr("amount quantity illegal", "")
	_                    = NewErr("copyright price illegal", "")
	_                    = NewErr("copyright charge rate illegal", "")
	_                    = NewErr("decimal point cannot exceed 18 digits", "6")
)
