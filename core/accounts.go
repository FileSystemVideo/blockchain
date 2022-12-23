package core

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

/**
: fsv17xpfvakm2amg962yls6f84z3kell8c5lmrfnas
:  fsv1gwqac243g2z3vryqsev6acq965f9ttwhjpnmg9
:  fsv1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8v9w0lj
: fsv1fl48vsnmsdzcv85q5d2q4z5ajdha8yu37pru06
: fsv1tygms3xhhs3yv487phx3dw4a95jn7t7l2pldew
*/


var ContractAddressFee = authtypes.NewModuleAddress(authtypes.FeeCollectorName)


var ContractAddressBank = authtypes.NewModuleAddress(bankTypes.ModuleName)


var ContractAddressDistribution = authtypes.NewModuleAddress(distrtypes.ModuleName)

//staking 
var ContractAddressStakingBonded = authtypes.NewModuleAddress(stakingtypes.BondedPoolName)

//staking 
var ContractAddressStakingNotBonded = authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName)

//gov
var ContractAddressGov = authtypes.NewModuleAddress(govtypes.ModuleName)


var ContractAddressBonus = authtypes.NewModuleAddress(KeyCopyrightBonus)


var ContractAddressMortgage = authtypes.NewModuleAddress(KeyCopyrightMortgage)


var ContractAddressDeflation = authtypes.NewModuleAddress(KeyCopyrighDeflation)


var ContractAddressDestory = sdk.AccAddress([]byte(sdk.BlackHoleAddress))

//IBC
var ContractAddressIbcTransfer = authtypes.NewModuleAddress(ibcTransferTypes.ModuleName)


const CommunityRewardAccount string = "fsv1gna7wuqcqp3l8nrw3q59fsc64hqep3hsulk234"


const CrossChainInManageAccount string = "fsv16knzs948zchx8dlaxl5tey5hs0hxr2rzj2uvfa"

//fsv1upm8ejcj4hlnqd6f5gsegwgrpmq490s4qqwd5z

const CrossChainAccount string = "fsv16knzs948zchx8dlaxl5tey5hs0hxr2rzj2uvfa"

//fsv1q48gzqz0lznegex8ymfcvpsp22nlxwrmapqqaa

const CrossChainFeeAccount string = "fsv17wf6w29k9g5s5qccqz2kcalqq8dchlskx297hc"


const CrossChainAutoDump string = "fsv1yvwstgw8q398cjtvq6vxu2n6z3ec5tc6wsal9t"
