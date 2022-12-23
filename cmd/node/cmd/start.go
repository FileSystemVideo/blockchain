package cmd

import (
	"context"
	vcmd "fs.video/blockchain/cmd/vcd/cmd"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/config"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"os"
	"path/filepath"
	"time"
)

//fsv
func Start(cosmosRepoPath string, ctx context.Context) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	var err error
	config.BlackList = types.DefaultParams().BlackList
	config.WhiteList = types.DefaultParams().WhiteList

	log.Info("check config dir")
	//config,
	if e, _ := CheckSourceExist(filepath.Join(cosmosRepoPath, "config")); !e {
		log.Info("init chain")

		
		err = ChainRun(context.Background(), core.CommandName, "init", "node1", "--chain-id", core.ChainID, "--home", cosmosRepoPath)
		if err != nil {
			panic(err)
		}
		
		err = replaceConfig(cosmosRepoPath)
		if err != nil {
			panic(err)
		}
	}
	log.Info("check genesis.json")
	//genesis.json,,
	err = checkGenesisFile(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	log.Info("check config.toml")
	//config.toml,,
	err = checkConfigFile(cosmosRepoPath)
	if err != nil {
		panic(err)
	}
	log.Info("check client.toml")
	//client.toml,,
	err = checkClientToml(cosmosRepoPath)
	if err != nil {
		panic(err)
	}
	log.Info("check app.toml")
	//app.toml
	err = checkAppToml(cosmosRepoPath)
	if err != nil {
		panic(err)
	}
	log.Info("check priv_validator_state.json")
	//priv_validator_state.json
	err = checkValidatorStateJson(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	log.WithField("path", cosmosRepoPath).Info("chain repo")

	//logPath := filepath.Join(cosmosRepoPath, "chain.log") //cosmos

	//  debug | info | error
	logLevel := "error"
	logLevelSet, ok := os.LookupEnv("CHAIN_LOGGING") //chain
	if ok {
		logLevel = logLevelSet 
	}
	
	go func() {
		for {
			//cosmos
			isOk := util.TestConnectivity("http://127.0.0.1:1317", time.Second*1)
			if isOk {
				break
			}
			<-time.After(time.Second) 
		}
	}()
	
	log.Info("start chain")
	os.Args = []string{core.CommandName, "start", "--log_format", "json", "--log_level", logLevel, "--home", cosmosRepoPath}
	//os.Args = []string{core.CommandName, "start", "--log_format", "json", "--log_level", logLevel, "--log-file", logPath, "--home", cosmosRepoPath}
	rootCmd, _ := vcmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, cosmosRepoPath, ctx); err != nil {
		//os.Exit panic
		os.Exit(1)
	}
}
