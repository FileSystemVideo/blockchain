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
	InvalidFromAddress  = NewErr("^invalid from address: (.*)$", "发起账号不合法")
	ParseError          = NewErr("failed to parse request", "格式化请求数据失败")
	ChainIdError        = NewErr("chain-id required but not specified", "链编号必填")
	FeeGasError         = NewErr("cannot provide both fees and gas prices", "手续费和价格不能同时为零")
	FeeGasInvalid       = NewErr("invalid fees or gas prices provided", "手续费或价格不能同时为空")
	PassGroupError      = NewErr("password is invalid", "密码组不合法")
	ChargeRateEmpty     = NewErr("chargeRate is empty", "佣金比例必填")
	ChargeRateTooSmall  = NewErr("chargeRate is too small", "佣金比例太小")
	ChargeRateTooHigh   = NewErr("chargeRate is too high", "佣金比例太高")
	Published           = NewErr("current datahash has published", "当前版权资源哈希已发布")
	QueryBindError      = NewErr("query bind relationship error", "查询账号版权上链资料失败")
	HasNoBind           = NewErr("current account has not bind relationship", "当前账号还未申请版权上链资料")
	InvalidAddress      = NewErr("invalid address", "账号地址不合法")
	InsufficientError   = NewErr("insufficient funds", "余额不足")
	MarshalError        = NewErr("failed to marshal JSON bytes", "序列化失败")
	UnMarshalError      = NewErr("failed to unmarshal JSON bytes", "序列化失败")
	DataHashEmpty       = NewErr("datahash is empty", "资源哈希为空")
	OriginDataHashEmpty = NewErr("origindatahash is empty", "原始资源哈希为空")
	DataHashNotEmpty    = NewErr("datahash cannot be empty", "资源哈希不能为空")
	GasAdjustmentError  = NewErr("invalid gas adjustment", "价格调整不合法")
	AccountNoExist      = NewErr("account not exist", "账号不存在")
	AccountError        = NewErr("decoding Bech32 address failed: must provide an address", "解析账号地址失败")
	CoinError           = NewErr("coin can not be empty", "币不能为空")
	IllegalError        = NewErr("Illegal account proportion", "分享账号占比不合法")
	DeleteError         = NewErr("delete fee amount is not invalid", "删除版权金额不合法")
	HasDeleted          = NewErr("current datahash has deleted", "当前资源哈希已删除")
	BalanceError        = NewErr("Insufficient account balance", "账号余额不足")
	AccountEmpty        = NewErr("account can not be empty", "账号不能为空")
	PageInvalid         = NewErr("page is invalid", "页码格式错误")
	PageSizeINvalid     = NewErr("pageSize is invalid", "页面数据格式错误")
	TokenEmpty          = NewErr("token is empty", "币种为空")
	QueryParamError     = NewErr("query param empty", "入参为空")
	NumError            = NewErr("query param num invalid", "数量参数为空")
	QueryExtError       = NewErr("query extype invalid", "")
	IdsError            = NewErr("query param idsString is empty", "")
	TxhashEmpty         = NewErr("query param txhash empty", "交易哈希为空")
	TxhashNotExist      = NewErr("txhash not exist", "交易哈希不存在")
	DeleAddressError    = NewErr("must use own delegator address", "POS抵押账号不合法")
	QueryDeletorError   = NewErr("query delegators amount error", "查询POS抵押余额错误")
	UnbondInfuffiError  = NewErr("unbond amount is not enough", "解绑余额不足")
	DeleAddressEmpty    = NewErr("delegatorAddr can not be empty", "POS抵押账号不能为空")
	_                   = NewErr("param author is empty", "作者不能为空")
	_                   = NewErr("author not empty", "著作人为空")
	_                   = NewErr("param productor is empty", "制片人不能为空")
	_                   = NewErr("param legalNumber is empty", "版权序列号不能为空")
	_                   = NewErr("param legalTime is empty", "登记时间不能为空")
	_                   = NewErr("param complainInfor is empty", "申诉描述不能为空")
	_                   = NewErr("param complainAccount is empty", "申诉账号不能为空")
	_                   = NewErr("param accuseAccount is empty", "被告账号不能为空")
	_                   = NewErr("do not appeal to yourself", "不能申诉自己")
	_                   = NewErr("current datahash is appealing", "当前资源哈希在申诉中")
	_                   = NewErr("complainType can not be empty", "申诉类型不能为空")
	_                   = NewErr("complainType is invalid", "申诉类型不合法")
	_                   = NewErr("aram accuseInfor is empty", "申诉回复描述不能为空")
	_                   = NewErr("param complainId is empty", "版权申诉编号不能为空")
	_                   = NewErr("complain not exist", "当期版权申诉不存在")
	_                   = NewErr("current complain has response", "版权申诉已回复")
	_                   = NewErr("current account  has no right to response", "当前账号无权回复当前版权申诉")
	_                   = NewErr("param voteStatus is empty", "投票状态不能为空")
	_                   = NewErr("param voteStatus is invalid", "投票状态不合法")
	_                   = NewErr("current complain status invalid", "当前版权申诉状态不合法")
	_                   = NewErr("complain vote can not repeat", "版权申诉投票不能重复")
	_                   = NewErr("current account has no vote right", "当前账号无投票权")
	_                   = NewErr("current account has no right to get complain result", "当前账号无权获取版权申诉结果")
	_                   = NewErr("current complain has finished", "当前版权申诉已结束")
	_                   = NewErr("complain vote has not reached end time", "版权申诉还未到结束时间")
	_                   = NewErr("calculate vote amount error", "计算投票数量错误")
	_                   = NewErr("query copyright complain error", "查看版权申诉失败")
	_                   = NewErr("complain account is invalid", "版权申诉账号不合法")
	_                   = NewErr("complain status does not allow query for ip", "版权申诉状态不允许查看网络地址")
	_                   = NewErr("format gas error", "价格格式化失败")
	_                   = NewErr("invalid gas adjustment", "价格调整参数不合法")
	_                   = NewErr("query copyright bonus status error", "查阅版权分红状态失败")
	_                   = NewErr("datahash does not exist", "资源哈希不存在")
	_                   = NewErr("datahash has exist", "资源哈希已经存在")
	_                   = NewErr("down copyright price error", "下载扣费价格错误")
	_                   = NewErr("datahash has download", "已经支付过下载费用")
	_                   = NewErr("query data error", "查询数据失败")
	_                   = NewErr("datahash creator error", "资源哈希创建者失败")
	_                   = NewErr("copyright bonus has demand", "版权下载已分红")
	_                   = NewErr("account bind has exist", "账号绑定信息已存在")
	_                   = NewErr("current accout has no right to delete", "当前账号无权删除")
	_                   = NewErr("copyright has deleted", "当前账号已删除")
	_                   = NewErr("copyright files is empty", "资源信息文件列表为空")
	_                   = NewErr("copyright complain account is invalid", "版权申诉账号不合法")
	_                   = NewErr("complainId does not exist", "版权申诉编号不存在")
	_                   = NewErr("complain status is invalid", "版权申诉状态不合法")
	_                   = NewErr("complain vote status is invalid", "版权申诉投票状态不合法")
	_                   = NewErr("current account exist block chain request ,please wait a minute", "当前账号已存在上链请求，请稍后再试")
	_                   = NewErr("sign error", "钱包密码验证错误")
	_                   = NewErr("get accountManage error", "获取账号管理器失败")
	_                   = NewErr("Entropy length must be \\[128, 256\\] and a multiple of 32", "生成助记词位数错误")
	_                   = NewErr("account only exist", "账号已存在")
	_                   = NewErr("encoding bech32 failed", "编码类型失败")
	_                   = NewErr("password verification failed", "账号密码错误")
	_                   = NewErr("account not exist", "账号不存在")
	_                   = NewErr("account key not exist", "账号秘钥未找到")

	_ = NewErr("failed to decrypt private key", "私钥解密失败")
	_ = NewErr("invalid mnemonic", "助记词格式错误")
	_ = NewErr("height must be equal or greater than zero", "块高度不能小于零")
	_ = NewErr("empty delegator address", "委托账号为空")
	_ = NewErr("empty validator address", "POS矿工器地址为空")
	_ = NewErr("invalid delegation amount", "委托金额不合法")
	_ = NewErr("invalid shares amount", "解绑金额不合法")
	_ = NewErr("validator does not exist", "POS矿工地址不存在")
	_ = NewErr("invalid coin denomination", "委托币种错误")
	_ = NewErr("delegate progress error", "股权质押失败")
	_ = NewErr("no validator distribution info", "没有验证程序分发信息")
	_ = NewErr("no delegation distribution info", "无委派分发信息")
	_ = NewErr("module account (.*)$ does not exist", "模块账号不存在")
	_ = NewErr("no validator commission to withdraw", "无验证器佣金可提取")
	_ = NewErr("signature verification failed; verify correct account sequence and chain-id", "签名验证失败,请检查账号序列号和链ID")
	_ = NewErr("database error", "数据库错误")
	_ = NewErr("parse account error", "账号格式化错误")
	_ = NewErr("parse coin error", "价格错误")
	_ = NewErr("parse string to number error", "类型转化失败")
	_ = NewErr("query chain infor errors", "查询链上数据失败")
	_ = NewErr("valid chain request error", "请求参数验证失败")
	_ = NewErr("parse json error", "格式化失败")
	_ = NewErr("parse byte to struct error", "类型转化错误")
	_ = NewErr("format tx struct error", "交易数转化失败")
	_ = NewErr("broadcast error", "广播失败")
	_ = NewErr("format string to int error", "类型转化失败")
	_ = NewErr("datahash format account", "资源哈希转化账号失败")
	_ = NewErr("parse valitor error", "验证地址格式化失败")
	_ = NewErr("parse time error", "时间格式化失败")
	_ = NewErr("current account has payed for datahash", "已经支付过该资源")
	_ = NewErr("current account has vote for this complain", "当前账号已经投票")
	_ = NewErr("datahash not exist", "资源不存在")
	_ = NewErr("sign error", "钱包密码验证错误")
	_ = NewErr("operator address exist", "当前钱包地址已经申请过POS矿工了")
	_ = NewErr("pub key for validator exist", "当前公钥已生成验证器")
	_ = NewErr("current validator not jail", "POS矿工非监禁状态")
	_ = NewErr("has no right to oprate validator", "无权解禁")
	_ = NewErr("validator description length error", "矿工信息过长")
	_ = NewErr("validator min mortgage amount", "创建POS矿工最小抵押数50")
	_ = NewErr("current account does not have right to delete datahash ", "当前账号无权删除")
	_ = NewErr("too many unbonding delegation entries for \\(delegator, validator\\) tuple", "当前申请OS赎回的记录已达到上限,需等待到账后再操作.")
	_ = NewErr("no delegation for \\(address, validator\\) tuple", "钱包地址下没有POS抵押金额")
	_ = NewErr("parse validator address error", "解析POS矿工地址出错")
	_ = NewErr("There is no reward to receive", "没有奖励可以领取")
	_ = NewErr("delegation does not exist", "没有抵押数据")
	_ = NewErr("query sensitive words error", "查询敏感词错误")
	_ = NewErr("save sensitive words error", "保存敏感词错误")
	_ = NewErr("sensitive status illegal", "敏感词状态不合法")
	_ = NewErr("origindatahash is empty", "原始hash不能为空")
	_ = NewErr("origindatahash has exist", "原始hash已存在")
	_ = NewErr("origindatahash not exist", "原始hash不存在")
	_ = NewErr("contain sensitive words", "包含敏感词")
	_ = NewErr("fee can not be zero", "手续费不能为零")
	_ = NewErr("fee is too less", "手续费太小")
	_ = NewErr("fee can not empty", "手续费不能为空")
	_ = NewErr("delegation coin less then min", "抵押金额最多支持到小数点后六位")
	_ = NewErr("unbonding delegation shares less then min", "赎回投票权最多支持到小数点后六位")
	_ = NewErr("delegation reward coin less then min", "POS奖励金额单项超过0.000001才可以领取")
	_ = NewErr("not enough delegation shares", "投票权不足")
	_ = NewErr("must use own validator address", "必须使用验证器所属的地址发起请求")
	_ = NewErr("private key is empty", "私钥为空")
	_ = NewErr("^verification fail$", "验签失败")
	_ = NewErr("pubkey not exist", "公钥不存在")
	_ = NewErr("verification error", "验签失败")
	_ = NewErr("con address is invalid", "共识地址不合法")
	_ = NewErr("con address can not empty", "共识地址不能为空")
	_ = NewErr("classify id is invalid", "分类id不合法")
	_ = NewErr("dir name has exist", "目录已存在")
	_ = NewErr("delegator does not contain delegation", "未找到POS抵押信息")
	_ = NewErr("The account to be unlocked must have a valid POS mortgage", "申请解禁的账号必须存在有效的POS抵押")
	_ = NewErr("The mortgage amount is less than the self-mortgage amount", "申请解禁的账号抵押值小于最小自抵押值")
	_ = NewErr("dir type illegal", "目录类型不合法")
	_ = NewErr("dir not exist", "目录不存在")
	_ = NewErr("query classify list error", "查询分类列表错误")
	_ = NewErr("query classify copyright list error", "查询分类下的版权列表错误")
	_ = NewErr("query classify list error", "查询分类列表错误")
	_ = NewErr("copyright not exist", "版权不存在")
	_ = NewErr("Illegal share account", "账号分享比例错误")
	_ = NewErr("current account does not have right to editor datahash ", "当前账号无权修改版权")
	_ = NewErr("current account does not have right to delete datahash ", "当前账号无权删除版权")
	_ = NewErr("lower min coin", "超出最小精度")
	_ = NewErr("cannot move to subdirectory", "不能移动自己的子目录")
	_ = NewErr("cannot move to self directory", "不能移动到自身目录")
	_ = NewErr("mortg miner has finish", "抵押挖矿已结束")
	_ = NewErr("SendCoins error", "转账失败")
	_ = NewErr("current classify has exist the same name", "当前目录下已存在相同名称的分类")
	_ = NewErr("classify not dir", "当前分类不是目录")
	_ = NewErr("bindinfor has exist", "绑定关系已存在")
	_ = NewErr("vote index not empty", "投票id不能为空")
	_ = NewErr("has not reach mortgage height", "还未达到抵押高度")
	_ = NewErr("miner space not enough", "硬盘空间不足")
	_ = NewErr("param inviteCode is empty", "邀请码为空")
	_ = NewErr("token id not empty", "nft tokenId 不能为空")
	_ = NewErr("has vote deflation", "当期已经存在投票信息")
	_ = NewErr("empty description: invalid request", "请填写矿工信息")
	_ = NewErr("signature verification failed, invalid chainid or account number", "验签失败,无效的chainId或者账号number")
	_ = NewErr("account serial number expired, the reason may be: node block behind or repeatedly sent messages", "账号序列号过期,原因可能是:节点区块落后或者重复发送了消息.")
	_ = NewErr("copyright ID has been used, please register again", "版权id已经使用，请重新注册")
	_ = NewErr("The copyright ID is empty", "版权id为空")
	_ = NewErr("Verifier information can only be changed once in 24 hours", "验证器信息24小时内仅可更改一次")
	_ = NewErr("Please use the main currency", "请使用主币")
	_ = NewErr("chain error", "上链失败")
	_ = NewErr("min self delegation cannot be zero", "最小自委托不能为零")

	_ = NewErr("commission must be positive", "佣金必须是整数")
	_ = NewErr("commission cannot be more than 100%", "佣金不能超过100%")
	_ = NewErr("commission cannot be more than the max rate", "佣金不能超过最高费率")
	_ = NewErr("commission cannot be changed more than once in 24h", "佣金修改24小时只能执行一次")
	_ = NewErr("commission change rate must be positive", "更改的佣金必须是整数")
	_ = NewErr("commission change rate cannot be more than the max rate", "更改的佣金不能超过最高费率")
	_ = NewErr("commission cannot be changed more than max change rate", "佣金变化不能超过最大变化率")
	_ = NewErr("validator's self delegation must be greater than their minimum self delegation", "DPOS矿工的抵押金额必须大于最小自抵押")
	_ = NewErr("minimum self delegation must be a positive integer", "最小自抵押必须是整数")
	_ = NewErr("minimum self delegation cannot be decrease", "不能减少最小自抵押")
	_ = NewErr("copyright vote power not enough", "版权审核投票票数不足")
	_=  NewErr("tip balance must greate than one","tip 余额大于1")
	_= NewErr("current datahash has not approve","当前资源还未审核")
	_= NewErr("The resource is not in the approval period","资源不在审核期")

)
