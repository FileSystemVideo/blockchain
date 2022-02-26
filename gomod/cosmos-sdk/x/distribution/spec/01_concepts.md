<!--
order: 1
-->

# Concepts

## Reference Counting in F1 Fee Distribution

在F1费用分配中，为了计算委托人撤回其委托时应获得的奖励，我们必须阅读从委托结束的期间到最后一个期间（撤回时创建）的奖励总和除以代币的条款。
此外，由于削减改变了委托人将拥有的代币数量（但我们只是在委托人取消委托的情况下才延迟计算），因此我们必须在任何削减之前/之后计算奖励，这些削减发生在委托人委托和他们撤回奖励之间。
因此，削减和代表团一样，指的是斜杠事件结束的时间段。

因此，任何委派或任何削减不再引用的期间的所有存储的历史奖励记录都可以安全地删除，因为它们永远不会被读取（将来的委派和将来的削减将始终引用将来的期间）。
这是通过跟踪“ReferenceCount”以及每个历史奖励存储条目来实现的。

每次创建可能需要引用历史记录的新对象（委派或斜杠）时，引用计数都会递增。
每次删除一个以前需要引用历史记录的对象时，引用计数就会减少。
如果引用计数为零，则删除历史记录。