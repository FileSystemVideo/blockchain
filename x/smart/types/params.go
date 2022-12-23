package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// NewParams creates a new Params object
func NewParams() Params {
	return Params{}
}

func DefaultParams() Params {

	return Params{}
}

func (p Params) Validate() error {
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable()
}
