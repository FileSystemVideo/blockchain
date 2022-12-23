package cmd

import (
	"context"
	"fs.video/blockchain/app"
	"fs.video/blockchain/cmd/vcd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"os"
)

func ChainRun(ctx context.Context, arg ...string) error {
	os.Args = arg
	rootCmd, _ := cmd.NewRootCmd()
	//rootCmd.SetOut(os.Stdout)
	//rootCmd.SetErr(os.Stderr)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome, ctx); err != nil {
		//os.Exit panic
		//os.Exit(1)
		return err
	}
	return nil
}
