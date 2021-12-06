package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var ContractAddressFee = authtypes.NewModuleAddress(authtypes.FeeCollectorName)

var BankAddress = authtypes.NewModuleAddress(bankTypes.ModuleName)

var ContractAddressDistribution = authtypes.NewModuleAddress(distrtypes.ModuleName)

var ContractAddressStakingBonded = authtypes.NewModuleAddress(stakingtypes.BondedPoolName)

var ContractAddressStakingNotBonded = authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName)

var ContractAddressBonus = authtypes.NewModuleAddress(KeyCopyrightBonus)

var ContractAddressMortgage = authtypes.NewModuleAddress(KeyCopyrightMortgage)

var ContractAddressDeflation = authtypes.NewModuleAddress(KeyCopyrighDeflation)

var ContractAddressDestory = sdk.AccAddress([]byte(sdk.BlackHoleAddress))

const InvestorAddress = "fsv1ycmwsuvq9x74pe0xfuk22r9zz8phl0hz437plj"

const CommunityRewardAccount string = "fsv1gna7wuqcqp3l8nrw3q59fsc64hqep3hsulk234"

const CrossChainInManageAccount string = "fsv16knzs948zchx8dlaxl5tey5hs0hxr2rzj2uvfa"

const CrossChainAccount string = "fsv16knzs948zchx8dlaxl5tey5hs0hxr2rzj2uvfa"

const CrossChainFeeAccount string = "fsv17wf6w29k9g5s5qccqz2kcalqq8dchlskx297hc"

const CrossChainAutoDump string = "fsv1yvwstgw8q398cjtvq6vxu2n6z3ec5tc6wsal9t"
