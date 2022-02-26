<!--
order: 0
title: Distribution Overview
parent:
  title: "distribution"
-->

# `distribution`

## Overview


这个简单的分配机制描述了一种在验证者和委托者之间被动分配奖励的功能性方法。
请注意，这一机制并不像积极的奖励分配机制那样精确地分配资金，因此将在未来升级。

该机构的工作原理如下。收集到的奖励在全球范围内汇集，并被动地分配给验证者和授权者。
每位验证人有机会就代表授权人收取的奖励向授权人收取佣金。
费用直接收集到全球奖励池和验证提议者奖励池中。
由于被动会计的性质，每当影响报酬分配率的参数发生变化时，也必须收回报酬。

- 每次提取时，必须提取他们有权获得的最大金额，池中不留任何东西。
- 当绑定、解除绑定或将代币重新委托给现有帐户时，必须全额提取奖励（因为懒散会计规则发生了变化）。
- 每当验证人选择更改奖励佣金时，必须同时撤销所有累积的佣金奖励。

以上场景在 `hooks.md`.

本文概述的分发机制用于在验证器和相关委托人之间延迟分发以下奖励：

- 多代币费用社会化分配
- 申请人奖励池
- 虚增的代币准备金
- 委托人获得的所有奖励的验证人佣金

费用汇集在一个全球池中，以及验证特定的提议者奖励池中。
所使用的机制允许验证者和授权者独立地、懒洋洋地撤回他们的奖励。

## Shortcomings

作为延迟计算的一部分，每个委托人持有一个特定于每个验证器的累积项， 用于估计他们在全球费用池中持有的代币的大致公平部分欠他们多少。

```
entitlement = delegator-accumulation / all-delegators-accumulation
```

在每个区块有恒定且相等的传入奖励令牌流的情况下，该分配机制将等于主动分配（每个区块分别分配给所有委托人）。
然而，这是不现实的，因此，根据传入奖励代币的波动以及其他委托人的奖励撤销时间，会出现偏离主动分配的情况。

如果你碰巧知道即将到来的奖励将显著增加，你会受到激励，在该事件发生之前不要退出，从而增加你现有的累积价值。
 参考 [#2764](https://github.com/cosmos/cosmos-sdk/issues/2764) 更多细节。

## Effect on Staking

对Token条款收取佣金，同时允许Token条款自动绑定（直接分发给验证者绑定的股权）在BPoS中是有问题的。
从根本上说，这两种机制是相互排斥的。
如果委托和自动绑定机制同时应用于锁紧令牌，则任何验证器及其委托者之间的锁紧令牌分布将随每个块而改变。
这就需要对每个块的每个委托记录进行计算——这在计算上是昂贵的。

总之，我们只能有token佣金和未绑定的token条款，或者没有token委托的token条款，我们选择实施前者。
希望重新发行其条款的利益相关者可以选择建立一个脚本，定期撤回和重新发行奖励。

## Contents

1. **[Concepts](01_concepts.md)**
    - [Reference Counting in F1 Fee Distribution](01_concepts.md#reference-counting-in-f1-fee-distribution)
2. **[State](02_state.md)**
3. **[End Block](03_end_block.md)**
4. **[Messages](04_messages.md)**
    - [MsgSetWithdrawAddress](04_messages.md#msgsetwithdrawaddress)
    - [MsgWithdrawDelegatorReward](04_messages.md#msgwithdrawdelegatorreward)
        - [Withdraw Validator Rewards All](04_messages.md#withdraw-validator-rewards-all)
    - [Common calculations ](04_messages.md#common-calculations-)
5. **[Hooks](05_hooks.md)**
    - [Create or modify delegation distribution](05_hooks.md#create-or-modify-delegation-distribution)
    - [Commission rate change](05_hooks.md#commission-rate-change)
    - [Change in Validator State](05_hooks.md#change-in-validator-state)
6. **[Events](06_events.md)**
    - [BeginBlocker](06_events.md#beginblocker)
    - [Handlers](06_events.md#handlers)
7. **[Parameters](07_params.md)**
