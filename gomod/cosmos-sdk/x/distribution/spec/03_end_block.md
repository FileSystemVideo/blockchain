<!--
order: 3
-->

# End Block

在每个 `EndBlock` 中，收到的费用都会转移到 `ModuleAccount`分发帐户中，因为它是一个跟踪硬币进出模块的帐户。
费用也分配给提议者、社区和全球基金。

当验证者是这一轮的提议者时，验证者（及其委托人）获得1%至5%的费用奖励，然后收取储备社区税，然后剩余部分通过投票权按比例分配给所有绑定的验证者，与他们是否投票无关（社会分配）。
注：除提议者奖励外，社会分配也适用于提议者验证器。

提议者奖励的金额是从预提交Tendermint消息中计算出来的，以便激励验证器等待并在块中包含额外的预提交。

爆块奖励 = 当前累计dpos资金池 * (0.01(爆块基础奖励) + 0.04(爆块额外奖励) * ( 参与投票的股权数 / 全网累计股权数 ) ))

所有供应奖励都添加到验证器单独持有的供应奖励池（`ValidatorDistribution.ProvisionRewardPool`）。

```go
func AllocateTokens(feesCollected sdk.Coins, feePool FeePool, proposer ValidatorDistribution, 
              sumPowerPrecommitValidators, totalBondedTokens, communityTax, 
              proposerCommissionRate sdk.Dec)

     SendCoins(FeeCollectorAddr, DistributionModuleAccAddr, feesCollected)
     feesCollectedDec = MakeDecCoins(feesCollected)
     
     proposerReward = feesCollectedDec * (0.01 + 0.04 
                       * sumPowerPrecommitValidators / totalBondedTokens)

     commission = proposerReward * proposerCommissionRate
     proposer.PoolCommission += commission
     proposer.Pool += proposerReward - commission

     communityFunding = feesCollectedDec * communityTax
     feePool.CommunityFund += communityFunding

     poolReceived = feesCollectedDec - proposerReward - communityFunding
     feePool.Pool += poolReceived

     SetValidatorDistribution(proposer)
     SetFeePool(feePool)
```
