<!--
order: 2
-->

# Messages

在本节中，我们将描述危机消息的处理以及相应的状态更新。

## MsgVerifyInvariant


可以使用 `MsgVerifyInvariant` 消息检查区块链不变量。

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/crisis/v1beta1/tx.proto#L14-L22

如果出现以下情况，则此消息将失败：
- 寄件人没有足够的token支付固定费用
- 未注册不变路由

这个消息检查提供的不变量，如果不变量被破坏，它就会panic，停止区块链。
如果不变量被破坏，那么固定费用就永远不会被扣除，因为交易永远不会被提交给一个区块（相当于被退款）。
但是，如果不变量没有被破坏，固定费用将不会退还。