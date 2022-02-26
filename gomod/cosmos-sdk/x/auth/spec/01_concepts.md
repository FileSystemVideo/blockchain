<!--
order: 1
-->

# Concepts

## Gas & Fees

对于网络节点来说，手续费有两个目的。

手续费(Fees)限制了每个完整节点所存储的状态的增长，并允许对几乎没有经济价值的交易进行通用审查。
手续费(Fees)最适合作为一种反垃圾邮件机制，其中验证程序对网络的使用和用户的身份不感兴趣。

手续费由交易提供的gasPrices和gasLimit决定，其中` fees = ceil(gasLimit * gasPrices)`.

tx会产生所有状态读/写、签名验证以及与tx大小成比例的 gas成本。

网络节点应在启动节点时设定最低gasPrices.

他们必须在希望支持的每个token中设定每个gas要扣除的成本:

`simd start ... --minimum-gas-prices=0.00001stake;0.05photinos`

在向mempool添加事务或gossip事务时， 验证者检查由提供的fee决定,交易gas价格是否符合验证者的任何最低gas价格。

验证者检查事务的gas价格，由提供的fees决定，是否符合验证者的最低gas价格

换言之，交易必须提供至少一种面额的fee，与验证者设定的最低gas价格相匹配。

Tendermint does not currently provide fee based mempool prioritization, and fee
based mempool filtering is local to node and not part of consensus. But with
minimum gas prices set, such a mechanism could be implemented by node operators.

Because the market value for tokens will fluctuate, validators are expected to
dynamically adjust their minimum gas prices to a level that would encourage the
use of the network.
