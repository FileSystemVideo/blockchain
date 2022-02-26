<!--
order: 1
-->

# State

## LastTotalPower

LastTotalPower tracks the total amounts of bonded tokens recorded during the previous end block.

- LastTotalPower: `0x12 -> ProtocolBuffer(sdk.Int)`

## Params

Params is a module-wide configuration structure that stores system parameters
and defines overall functioning of the staking module.

- Params: `Paramsspace("staking") -> legacy_amino(params)`

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.1/proto/cosmos/staking/v1beta1/staking.proto#L230-L241

## Validator

Validators can have one of three statuses

- `Unbonded`:验证器不在活动集中。他们不能签署区块，也不能赚取奖励。他们可以接受委托。.
- `Bonded`": 一旦验证器接收到足够的绑定代币，它们将在期间自动加入活动集，并将其状态更新为`Bonded`。
  他们签了区块并得到了奖励。
  他们可以接受更多的委托。
  他们可能因为行为不端而被削减。
  此验证程序的委托者如果解除了其委托的绑定，则必须等待解除绑定时间的持续时间，这是一个特定于链的参数。
  在此期间，如果源验证器的违规行为是在代币被绑定的时间段内发生的，那么它们仍然可以被削减。
- `Unbonding`: 当一个验证器离开活动集时，不管是选择还是由于削减或逻辑删除(永久监禁)，它们的所有委托都开始解除绑定。
  然后，所有委托必须等待解除绑定时间，然后才能从`BondedPool`将其代币转移到其帐户。


验证器对象应该主要由`OperatorAddr`存储和访问，OperatorAddr是验证器操作员的SDK验证器地址。
每个验证器对象维护两个额外的索引，以便完成削减和验证器集更新所需的查找。
第三个特殊索引（`LastValidatorPower`）也保持不变，但是它在每个块中保持不变，不像前两个索引在块中镜像验证器记录。

- Validators: `0x21 | OperatorAddr -> ProtocolBuffer(validator)`
- ValidatorsByConsAddr: `0x22 | ConsAddr -> OperatorAddr`
- ValidatorsByPower: `0x23 | BigEndian(ConsensusPower) | OperatorAddr -> OperatorAddr`
- LastValidatorsPower: `0x11 OperatorAddr -> ProtocolBuffer(ConsensusPower)`

`Validators` is the primary index - it ensures that each operator can have only one
associated validator, where the public key of that validator can change in the
future. Delegators can refer to the immutable operator of the validator, without
concern for the changing public key.

`ValidatorByConsAddr` 是一个额外的索引，用于查找削减。
当Tendermint报告证据时，它提供了验证器地址，因此需要这个映射来找到操作员。
注意，`ConsAddr`对应于可以从验证器的`ConsPubKey`派生的地址。

`ValidatorsByPower` is an additional index that provides a sorted list o
potential validators to quickly determine the current active set. Here
ConsensusPower is validator.Tokens/10^6. Note that all validators where
`Jailed` is true are not stored within this index.

`LastValidatorsPower` is a special index that provides a historical list of the
last-block's bonded validators. This index remains constant during a block but
is updated during the validator set update process which takes place in [`EndBlock`](./05_end_block.md).

Each validator's state is stored in a `Validator` struct:

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/staking/v1beta1/staking.proto#L65-L99

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/staking/v1beta1/staking.proto#L24-L63

## Delegation

Delegations are identified by combining `DelegatorAddr` (the address of the delegator)
with the `ValidatorAddr` Delegators are indexed in the store as follows:

- Delegation: `0x31 | DelegatorAddr | ValidatorAddr -> ProtocolBuffer(delegation)`

Stake holders may delegate coins to validators; under this circumstance their
funds are held in a `Delegation` data structure. It is owned by one
delegator, and is associated with the shares for one validator. The sender of
the transaction is the owner of the bond.

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/staking/v1beta1/staking.proto#L159-L170

### Delegator Shares

When one Delegates tokens to a Validator they are issued a number of delegator shares based on a
dynamic exchange rate, calculated as follows from the total number of tokens delegated to the
validator and the number of shares issued so far:

`Shares per Token = validator.TotalShares() / validator.Tokens()`

Only the number of shares received is stored on the DelegationEntry. When a delegator then
Undelegates, the token amount they receive is calculated from the number of shares they currently
hold and the inverse exchange rate:

`Tokens per Share = validator.Tokens() / validatorShares()`

These `Shares` are simply an accounting mechanism. They are not a fungible asset. The reason for
this mechanism is to simplify the accounting around slashing. Rather than iteratively slashing the
tokens of every delegation entry, instead the Validators total bonded tokens can be slashed,
effectively reducing the value of each issued delegator share.

## UnbondingDelegation

Shares in a `Delegation` can be unbonded, but they must for some time exist as
an `UnbondingDelegation`, where shares can be reduced if Byzantine behavior is
detected.

`UnbondingDelegation` are indexed in the store as:

- UnbondingDelegation: `0x32 | DelegatorAddr | ValidatorAddr -> ProtocolBuffer(unbondingDelegation)`
- UnbondingDelegationsFromValidator: `0x33 | ValidatorAddr | DelegatorAddr -> nil`

The first map here is used in queries, to lookup all unbonding delegations for
a given delegator, while the second map is used in slashing, to lookup all
unbonding delegations associated with a given validator that need to be
slashed.

A UnbondingDelegation object is created every time an unbonding is initiated.

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/staking/v1beta1/staking.proto#L172-L198

## Redelegation

The bonded tokens worth of a `Delegation` may be instantly redelegated from a
source validator to a different validator (destination validator). However when
this occurs they must be tracked in a `Redelegation` object, whereby their
shares can be slashed if their tokens have contributed to a Byzantine fault
committed by the source validator.

`Redelegation` are indexed in the store as:

- Redelegations: `0x34 | DelegatorAddr | ValidatorSrcAddr | ValidatorDstAddr -> ProtocolBuffer(redelegation)`
- RedelegationsBySrc: `0x35 | ValidatorSrcAddr | ValidatorDstAddr | DelegatorAddr -> nil`
- RedelegationsByDst: `0x36 | ValidatorDstAddr | ValidatorSrcAddr | DelegatorAddr -> nil`

The first map here is used for queries, to lookup all redelegations for a given
delegator. The second map is used for slashing based on the `ValidatorSrcAddr`,
while the third map is for slashing based on the `ValidatorDstAddr`.

A redelegation object is created every time a redelegation occurs. To prevent
"redelegation hopping" redelegations may not occur under the situation that:

- the (re)delegator already has another immature redelegation in progress
  with a destination to a validator (let's call it `Validator X`)
- and, the (re)delegator is attempting to create a _new_ redelegation
  where the source validator for this new redelegation is `Validator-X`.

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/staking/v1beta1/staking.proto#L200-L228

## Queues

All queues objects are sorted by timestamp. The time used within any queue is
first rounded to the nearest nanosecond then sorted. The sortable time format
used is a slight modification of the RFC3339Nano and uses the the format string
`"2006-01-02T15:04:05.000000000"`. Notably this format:

- right pads all zeros
- drops the time zone info (uses UTC)

In all cases, the stored timestamp represents the maturation time of the queue
element.

### UnbondingDelegationQueue

For the purpose of tracking progress of unbonding delegations the unbonding
delegations queue is kept.

- UnbondingDelegation: `0x41 | format(time) -> []DVPair`

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/staking/v1beta1/staking.proto#L123-L133

### RedelegationQueue

For the purpose of tracking progress of redelegations the redelegation queue is
kept.

- UnbondingDelegation: `0x42 | format(time) -> []DVVTriplet`

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/staking/v1beta1/staking.proto#L140-L152

### ValidatorQueue

For the purpose of tracking progress of unbonding validators the validator
queue is kept.

- ValidatorQueueTime: `0x43 | format(time) -> []sdk.ValAddress`

The stored object as each key is an array of validator operator addresses from
which the validator object can be accessed. Typically it is expected that only
a single validator record will be associated with a given timestamp however it is possible
that multiple validators exist in the queue at the same location.

## HistoricalInfo

HistoricalInfo objects are stored and pruned at each block such that the staking keeper persists
the `n` most recent historical info defined by staking module parameter: `HistoricalEntries`.

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/staking/v1beta1/staking.proto#L15-L22

At each BeginBlock, the staking keeper will persist the current Header and the Validators that committed
the current block in a `HistoricalInfo` object. The Validators are sorted on their address to ensure that
they are in a determisnistic order.
The oldest HistoricalEntries will be pruned to ensure that there only exist the parameter-defined number of
historical entries.
