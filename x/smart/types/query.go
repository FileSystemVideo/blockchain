package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	
	QueryGatewayInfo = "gateway_info"
	
	QueryGatewayList = "gateway_list"
	
	QueryGatewayNum = "gateway_num"
	
	QueryGatewayRedeemNum = "gateway_redeem_num"
	
	QueryValidatorByConsAddress = "validatorByConsAddress"
)


type QueryGatewayInfoParams struct {
	GatewayAddress  string `json:"gateway_address"`
	GatewayNumIndex string `json:"gateway_num_index"`
}


type QueryValidatorByConsAddrParams struct {
	ValidatorConsAddress sdk.ConsAddress
}
