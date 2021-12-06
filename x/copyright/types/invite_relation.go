package types

import (
	"github.com/shopspring/decimal"
)

type InviteRewardStatistics struct {
	InviteRewardSpace     decimal.Decimal `json:"invite_reward_space"`
	InviteRewardCounts    int             `json:"invite_reward_counts"`
	ExpansionRewardSpace  decimal.Decimal `json:"expansion_reward_space"`
	ExpansionRewardCounts int             `json:"expansion_reward_counts"`
}

type InviteRecording struct {
	Address    string          `json:"address"`
	InviteTime int64           `json:"invite_time"`
	Space      decimal.Decimal `json:"space"`
}


type Settlement struct {
	ExpansionRewardSpace decimal.Decimal `json:"expansion_reward_space"`
	InviteRewardSpace    decimal.Decimal `json:"invite_reward_space"`
}
