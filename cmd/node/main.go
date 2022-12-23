package main

import (
	"context"
	ncmd "fs.video/blockchain/cmd/node/cmd"
	"fs.video/blockchain/commands"
	"fs.video/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "chain",
		Short: "fsv chain node",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(commands.VersionCmd())
	rootCmd.AddCommand(commands.DposCmd())
	rootCmd.AddCommand(commands.ParamCmd())
	rootCmd.AddCommand(commands.StatusCmd())
	rootCmd.Execute()

	/*runtimePath, _ := os.Getwd()
	logFile := filepath.Join(runtimePath, "log.txt")
	rootCmd, _ := cmd.NewRootCmd()
	// log
	log.EnableLogStorage(logFile, time.Hour*24*7, time.Hour*24) 
	log.InitLogger(logrus.DebugLevel)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome, context.Background()); err != nil {
		os.Exit(1)
	}*/
}

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start chain",
		Run: func(cmd *cobra.Command, args []string) {

			programPath, _ := filepath.Abs(os.Args[0])

			//.exe 
			runtimePath, _ := filepath.Split(programPath)

			logPath := filepath.Join(runtimePath, "log")

			if !PathExists(logPath) {
				err := os.Mkdir(logPath, 0644)
				if err != nil {
					panic(err)
				}
			}

			daemonLogPath := filepath.Join(logPath, "chain.log")

			log.EnableLogStorage(daemonLogPath, time.Hour*24*7, time.Hour*24) 

			var cosmosRepoPath string
			home := cmd.Flag("home").Value.String()
			if home == "" {
				pwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				cosmosRepoPath = filepath.Join(pwd, ".data")
			} else {
				cosmosRepoPath = home
			}
			ncmd.Start(cosmosRepoPath, context.Background())
		},
	}
	cmd.Flags().String("home", "", "chain repo path")
	cmd.Flags().Bool("log", true, "enabled log storage")
	return cmd
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
