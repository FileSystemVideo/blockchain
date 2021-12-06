package cli

import (
	"fmt"
	//"github.com/cosmos/cosmos-sdk/client/flags"
	//"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"fs.video/blockchain/x/copyright/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
	//NewMsgSendTxCmd(),
	)
	return cmd
}

// NewSendTxCmd returns a CLI command handler for creating a MsgSend transaction.
//func NewMsgSendTxCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "send [from_key_or_address] [contentHash]",
//		Short: `Send funds from one account to another. Note, the'--from' flag is ignored as it is implied from [from_key_or_address].`,
//		Args:  cobra.ExactArgs(3),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			cmd.Flags().Set(flags.FlagFrom, args[0])
//			clientCtx, err := client.GetClientTxContext(cmd)
//			if err != nil {
//				return err
//			}
//			contentHash := args[2]
//
//			msg := types.NewMsgChat(clientCtx.GetFromAddress(), contentHash)
//			if err := msg.ValidateBasic(); err != nil {
//				return err
//			}
//			fmt.Println("0----------------------------------")
//			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
//		},
//	}
//
//	flags.AddTxFlagsToCmd(cmd)
//
//	return cmd
//}
