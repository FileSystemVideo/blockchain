package cli

import (
	"errors"
	"fmt"
	"fs.video/blockchain/x/smart/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// NewTxCmd returns a root CLI command handler for erc20 transaction commands
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "smart subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewCreateValidatorCmd(),
	)
	return txCmd
}

func NewCreateValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "create new validator initialized with a self-delegation to it",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildCreateValidatorMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(cli.FlagSetPublicKey())
	cmd.Flags().AddFlagSet(cli.FlagSetAmount())
	//cmd.Flags().AddFlagSet(cli.FlagSetDescriptionCreate())
	cmd.Flags().AddFlagSet(cli.FlagSetCommissionCreate())
	cmd.Flags().AddFlagSet(cli.FlagSetMinSelfDelegation())

	cmd.Flags().String(cli.FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(cli.FlagNodeID, "", "The node's ID")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(cli.FlagAmount)
	_ = cmd.MarkFlagRequired(cli.FlagPubKey)
	_ = cmd.MarkFlagRequired(cli.FlagMoniker)

	return cmd
}

func newBuildCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgCreateSmartValidator, error) {
	fAmount, _ := fs.GetString(cli.FlagAmount)
	amount, err := sdk.ParseCoinNormalized(fAmount)
	if err != nil {
		return txf, nil, err
	}

	valAddr := clientCtx.GetFromAddress()
	pkStr, err := fs.GetString(cli.FlagPubKey)
	if err != nil {
		return txf, nil, err
	}

	moniker, _ := fs.GetString(cli.FlagMoniker)
	identity, _ := fs.GetString(cli.FlagIdentity)
	website, _ := fs.GetString(cli.FlagWebsite)
	security, _ := fs.GetString(cli.FlagSecurityContact)
	details, _ := fs.GetString(cli.FlagDetails)
	description := stakingTypes.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)

	// get the initial validator commission parameters
	rateStr, _ := fs.GetString(cli.FlagCommissionRate)
	maxRateStr, _ := fs.GetString(cli.FlagCommissionMaxRate)
	maxChangeRateStr, _ := fs.GetString(cli.FlagCommissionMaxChangeRate)

	commissionRates, err := buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr)
	if err != nil {
		return txf, nil, err
	}

	// get the initial validator min self delegation
	msbStr, _ := fs.GetString(cli.FlagMinSelfDelegation)

	minSelfDelegation, ok := sdk.NewIntFromString(msbStr)
	if !ok {
		return txf, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	msg, err := types.NewMsgCreateSmartValidator(
		sdk.ValAddress(valAddr), pkStr, amount, description, commissionRates, minSelfDelegation,
	)
	if err != nil {
		return txf, nil, err
	}
	if err := msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	genOnly, _ := fs.GetBool(flags.FlagGenerateOnly)
	if genOnly {
		ip, _ := fs.GetString(cli.FlagIP)
		nodeID, _ := fs.GetString(cli.FlagNodeID)

		if nodeID != "" && ip != "" {
			txf = txf.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}

	return txf, msg, nil
}

func buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr string) (commission stakingTypes.CommissionRates, err error) {
	if rateStr == "" || maxRateStr == "" || maxChangeRateStr == "" {
		return commission, errors.New("must specify all validator commission parameters")
	}

	rate, err := sdk.NewDecFromStr(rateStr)
	if err != nil {
		return commission, err
	}

	maxRate, err := sdk.NewDecFromStr(maxRateStr)
	if err != nil {
		return commission, err
	}

	maxChangeRate, err := sdk.NewDecFromStr(maxChangeRateStr)
	if err != nil {
		return commission, err
	}

	commission = stakingTypes.NewCommissionRates(rate, maxRate, maxChangeRate)

	return commission, nil
}
