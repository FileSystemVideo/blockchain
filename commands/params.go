package commands

import (
	"fmt"
	"fs.video/blockchain/client"
	"fs.video/blockchain/core"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func ParamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "lookup params",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("chain-ID:", core.ChainID)
			fmt.Println("MainToken:", core.MainToken)
			fmt.Println("InviteToken:", core.InviteToken)
			fmt.Println("CoinPlaces:", core.CoinPlaces)
			fmt.Println("CommitTime:", core.CommitTime, "S")
			fmt.Println("-------------Fee-----------------------")
			fmt.Println("CopyrightFee:", core.ChainDefaultFee)
			fmt.Println("CopyrightInviteFee:", core.CopyrightInviteFee)
			fmt.Println("TransferRate:", core.TransferRate)
			fmt.Println("---------------SpaceMiner---------------------")
			fmt.Println("SpaceMinerBonusBlockNum:", core.SpaceMinerBonusBlockNum)
			fmt.Println("MinerUpperLimitStandand:", core.MinerUpperLimitStandand.String())
			fmt.Println("SpaceMinerPerDayStandand:", core.SpaceMinerPerDayStandand.String())
			fmt.Println("ValidatorMinerPerDayStandand:", core.ValidatorMinerPerDayStandand.String())
			fmt.Println("ChargeRateLow:", types.RemoveDecLastZero(core.ChargeRateLow), " ChargeRateHigh:", types.RemoveDecLastZero(core.ChargeRateHigh))
			fmt.Println("---------------Vote---------------------")
			fmt.Println("DeflationVoteDealBlockNum:", core.DeflationVoteDealBlockNum)
			fmt.Println("VoteResultTimePerioad:", core.VoteResultTimePerioad.Hours(), " hours")
			fmt.Println("CopyrightVoteTimePerioad:", core.CopyrightVoteTimePerioad.Hours(), " hours")
			fmt.Println("CopyrightVoteRedeemTimePerioad:", core.CopyrightVoteRedeemTimePerioad.Hours(), " hours")
			fmt.Println("CopyrightVoteAwardRate:", core.CopyrightVoteAwardRateV2.String())
			txClient := client.NewTxClient()

			param, err := txClient.GetParams()
			if err != nil {
				return
			}
			fmt.Println("---------------Staking params---------------------")
			fmt.Println("UnbondingTime:", param.StakingParam.UnbondingTime.Hours(), " hours")
			fmt.Println("MaxValidators:", param.StakingParam.MaxValidators)
			fmt.Println("HistoricalEntries:", param.StakingParam.HistoricalEntries)
			fmt.Println("MaxEntries:", param.StakingParam.MaxEntries)
			fmt.Println("---------------Distribution params---------------------")
			fmt.Println("BaseProposerReward:", types.RemoveDecLastZero(param.DistributionParam.BaseProposerReward))
			fmt.Println("BonusProposerReward:", types.RemoveDecLastZero(param.DistributionParam.BonusProposerReward))
			fmt.Println("CommunityTax:", types.RemoveDecLastZero(param.DistributionParam.CommunityTax))
			fmt.Println("---------------Slashing params---------------------")
			fmt.Println("DowntimeJailDuration:", param.SlashingParam.DowntimeJailDuration.Hours(), " hours")

			signedBlocksWindow := sdk.NewDec(param.SlashingParam.SignedBlocksWindow)
			fmt.Println("SignedBlocksWindow:", signedBlocksWindow.TruncateInt64())
			fmt.Println("MinSignedBlockPerWindow:", types.RemoveDecLastZero(param.SlashingParam.MinSignedPerWindow), " ≈ ", signedBlocksWindow.Mul(param.SlashingParam.MinSignedPerWindow).TruncateInt64(), " block")

			fmt.Println("SlashFractionDoubleSign:", types.RemoveDecLastZero(param.SlashingParam.SlashFractionDoubleSign))
			fmt.Println("SlashFractionDowntime:", types.RemoveDecLastZero(param.SlashingParam.SlashFractionDowntime))
			fmt.Println("---------------Hardware mining params---------------------")
		},
	}
	return cmd
}
