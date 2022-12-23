package commands

import (
	"encoding/base64"
	"fmt"
	"fs.video/blockchain/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"time"
)

func StatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "lookup status",
		Run: func(cmd *cobra.Command, args []string) {
			txClient := client.NewTxClient()
			dposClient := client.NewDposClient(&txClient)
			nodeClient := client.NewNodeClient(dposClient)
			statusInfo, err := nodeClient.StatusInfo()
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error StatusInfo | ", err.Error())
				fmt.Println("----------------------------------------------------------------------")
				return
			}
			fmt.Println("----------------------------------------------------------------------")
			fmt.Println("Sync Info")
			fmt.Println("LatestBlockHeight:", statusInfo.SyncInfo.LatestBlockHeight)
			fmt.Println("LatestBlockTime:", statusInfo.SyncInfo.LatestBlockTime)
			fmt.Println("LatestAppHash:", statusInfo.SyncInfo.LatestAppHash)
			fmt.Println("LatestBlockHash:", statusInfo.SyncInfo.LatestBlockHash)
			fmt.Println("EarliestBlockHeight:", statusInfo.SyncInfo.EarliestBlockHeight)
			fmt.Println("EarliestBlockTime:", statusInfo.SyncInfo.EarliestBlockTime)
			fmt.Println("EarliestAppHash:", statusInfo.SyncInfo.EarliestAppHash)
			fmt.Println("EarliestBlockHash:", statusInfo.SyncInfo.EarliestBlockHash)
			fmt.Println("----------------------------------------------------------------------")
			fmt.Println("Node Info")
			fmt.Println("ListenAddr:", statusInfo.NodeInfo.ListenAddr)
			fmt.Println("Network:", statusInfo.NodeInfo.Network)
			fmt.Println("NodeID:", statusInfo.NodeInfo.DefaultNodeID)
			fmt.Println("----------------------------------------------------------------------")
			fmt.Println("Validator Info")
			fmt.Println("Address:", statusInfo.ValidatorInfo.Address)
			fmt.Println("VotingPower:", statusInfo.ValidatorInfo.VotingPower)
			fmt.Println("PubKey:", base64.StdEncoding.EncodeToString(statusInfo.ValidatorInfo.PubKey.Bytes()))

			
			validatorConsAddress, err := sdk.ConsAddressFromHex(statusInfo.ValidatorInfo.Address.String())
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error ConsAddressFromHex | ", err.Error())
				fmt.Println("----------------------------------------------------------------------")
				return
			}
			fmt.Println("ConsAddress:", validatorConsAddress.String())

			validatorInfo, noFound, err := dposClient.FindValidatorByConsAddress(validatorConsAddress.String())
			if noFound {
				
				return
			}
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error FindValidatorByConsAddress | ", err.Error())
				fmt.Println("----------------------------------------------------------------------")
				return
			}

			//dpos
			signInfo, err := dposClient.FindSigningInfo(validatorConsAddress.String())
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error FindSigningInfo | ", err.Error())
				fmt.Println("----------------------------------------------------------------------")
				return
			}
			fmt.Println("----------------------------------------------------------------------")
			fmt.Println("Sign Info")
			if validatorInfo.Jailed {
				shanghaiLocation, _ := time.LoadLocation("Asia/Shanghai")
				fmt.Println("Jailed Status: jailed")
				fmt.Println("Jailed Height:", validatorInfo.UnbondingHeight)
				fmt.Println("Jailed End Time:", GetTimeStringFormatTime(signInfo.JailedUntil.In(shanghaiLocation)))
				//fmt.Println("Unbonding Time:",utilTime.GetTimeStringFormatTime(validatorInfo.UnbondingTime.In(shanghaiLocation)))
			} else {
				fmt.Println("Jailed Status: nothing")
			}

			fmt.Println("Tombstoned:", signInfo.Tombstoned)
			fmt.Println("Sign Start Height:", signInfo.StartHeight)
			fmt.Println("Signed blocks:", signInfo.IndexOffset)
			fmt.Println("Unsigned blocks:", signInfo.MissedBlocksCounter)
		},
	}
	return cmd
}

var YMRDHS_Format = "2006-01-02 15:04:05"


func GetTimeStringFormatTime(t time.Time) string {
	return t.Format(YMRDHS_Format)
}
