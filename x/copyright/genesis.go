package copyright

import (
	"fs.video/blockchain/x/copyright/keeper"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	//genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export

	//types.NewGenesisState()
	return types.NewGenesisState(k.ExportAccountSpace(ctx), k.ExportAccountInvite(ctx), k.ExportDeflationMinerInfor(ctx), k.ExportAccountRecord(ctx))
}
