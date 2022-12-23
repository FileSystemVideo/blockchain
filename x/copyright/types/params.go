package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyBlackList = []byte("BlackList")
	KeyWhiteList = []byte("WhiteList")
)

// NewParams creates a new Params object
func NewParams() Params {
	return Params{}
}

func DefaultParams() Params {
	blackList := []string{
		"fsv18d9xkh77evlw4fraerermq2hxnl3l6excrc3j4", 
		"fsv1wruv90vmum0wjta70hp45hsfe87p5wvsjtkd5n",
		"fsv12fxa4ny4thmuvz0clpcmfqrcm3fnwu6vyz8qk9",
		"fsv1khydjyzetyp9wy8exc34agvt7597cc9pakzvde",
		"fsv1u2qx4hvggt88e5uw79fz8xzqv3l5n75ndjwn90",
		"fsv1m6wgmcjlxvhs3c7zh6jx8dnul38lldre0g4s9e",
		"fsv1hg6yvpvwyvj0wycllh9lctyjj7s23remcsxtk9",
		"fsv1fmm620e7ytnjngyn3pwnc4jxp8ll8qnevq9kwa",
		"fsv1vj6mngvwhtfxk45hwfh0h0rk3g2js7wx4c95nx",
		"fsv1zdj0tqz3960k5ykedk3g4uf8eyrxlvz320jnwa",
		"fsv136mm4y5a4teg2av84fn2pgcdwzanadu38dnppe",
		"fsv19asjq7y3ujnzw5x4qlx0q7af0vme9jrnr7ewyg",
		"fsv1upx2ra40xc3xfsg2a9juw4ehhcrrtlk26sn6ts",
		"fsv1j4j6yyps2pap0yv6jq9qvqzwphnc22evjnq678",
		"fsv10hkp86w32t66nxqzjgr72pq444thlstxn5y88a",
		"fsv18gsqcjmnn6phn6y88g6h9t4d8v8rlnu9se7ktk",
		"fsv1lha3y9vctfrgh265cml2qvkkg8ryfmyu68sqn8",
		"fsv15mmz5x2hsrc5pm074jfxgndgamwm7et9d7u0h0",
		"fsv1gdg4z45emtj0dxff994dsvqyrjgm27yytxy70j",
		"fsv1g3zy2rl3wdf68nstw3pl4p7fds7fe7w3wy0d7t",
		"fsv10gtezag75w0cx8vpmdwpgkxhpa7hzvm6kfs4zv",
		"fsv1xxju0vr55k5dqwxrlt2ksx58ylfz4ymyfcjuwk",
	}
	whiteList := []string{
		"fsv1jw0cv47rvmw40jpgze03xq6wnt5jp58w6plnvx", 
		"fsv1jn0fvan693d8cqnzgmasq37uxnc0hht9pw32wr",
		"fsv16u46w9t380nmflsy9ptreg9njueqg828fjzjxy",
		"fsv18t4445dhugqkln5yjmndjgp8wufglgm28e3cfv",
		"fsv1ftlrw9fnn56j27327d2d8peq4d7f7swmc62jly",
		"fsv1hj8uf2ukqzet87k5r6ylxmt4we5puzc3h5v7nt",
		"fsv1wr5ehvfknrgqalj6ekwzasjx5dv7l97yg7q59l",
		"fsv1utflwxsegyh37adgvw292le8mdu2pk9ps9tfsn",
		"fsv14u44xqgus82x28pq5yyqr046wg9qnhdrszcjvy",
		"fsv1q9dtgjr96hv2wx56z2pe70n8yr7agnlr7jfk2q",
		"fsv1htnu2v2jxdyffu9xhvhcha5agqk05n8qm2578l",
		"fsv1pkc66v7apj532az7j2y3uedrwzpkmapd595c3l",
		"fsv1sd90hqdxju2kqetcr2lyazrv8da40sx69h56eu",
		"fsv1a2833h6hrujqxzhzcyhttnmzs380d6hw7fn86d",
		"fsv1w0amlrdqptw4q8qkqyad3k3fktdh54z4d998qp",
		"fsv1kv2w2yhy08agjan8hf8uss4wk0dy5p023pu8wt",
		"fsv1lyj8lxjzleve2f8fc7m85n2umres8jpfp09n06",
		"fsv1nd52manetndja20xkh45pqymjt35dc994xtehy",
		"fsv16u4teajf9w6eyn9snw7lxtlauqrszvq0v4czhn",
		"fsv1d2uxyth9m3yyxxnj0av4vuzkmdgejc30k43pxm",
		"fsv1wqap233cc94c6jqkqp7c8rwafjen2ny38a7lr0"}

	return Params{
		BlackList: blackList,
		WhiteList: whiteList,
	}
}

func (p Params) Validate() error {
	if err := validateBlackList(p.BlackList); err != nil {
		return err
	}
	if err := validateWhiteList(p.WhiteList); err != nil {
		return err
	}
	return nil
}

//paramtypes.NewParamSetPair(KeyMaxMemoCharacters, &p.MaxMemoCharacters, validateMaxMemoCharacters),
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBlackList, &p.BlackList, validateBlackList),
		paramtypes.NewParamSetPair(KeyWhiteList, &p.WhiteList, validateWhiteList),
	}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func validateBlackList(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, val := range v {
		_, err := sdk.AccAddressFromBech32(val)
		if err != nil {
			return fmt.Errorf("Illegal address:", val)
		}
	}
	return nil
}

func validateWhiteList(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, val := range v {
		_, err := sdk.AccAddressFromBech32(val)
		if err != nil {
			return fmt.Errorf("Illegal address:", val)
		}
	}
	return nil
}
