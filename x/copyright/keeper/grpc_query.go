package keeper

import (
	"fs.video/blockchain/x/copyright/types"
)

var _ types.QueryServer = Keeper{}
