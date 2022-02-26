<!--
order: 5
-->

# Hooks

## 创建或修改委派分发
 
 - triggered-by: `staking.MsgDelegate`, `staking.MsgBeginRedelegate`, `staking.MsgUndelegate`

The pool of a new delegator bond will be 0 for the height at which the bond was
added, or the withdrawal has taken place. This is achieved by setting
`DelegationDistInfo.WithdrawalHeight` to the height of the triggering transaction. 

## 佣金率变化
 
 - triggered-by: `staking.MsgEditValidator`

If a validator changes its commission rate, all commission on fees must be
simultaneously withdrawn using the transaction `TxWithdrawValidator`.
Additionally the change and associated height must be recorded in a
`ValidatorUpdate` state record.

## 验证程序状态更改
 
 - triggered-by: `staking.Slash`, `staking.UpdateValidator`

Whenever a validator is slashed or enters/leaves the validator group all of the
validator entitled reward tokens must be simultaneously withdrawn from
`Global.Pool` and added to `ValidatorDistInfo.Pool`. 
