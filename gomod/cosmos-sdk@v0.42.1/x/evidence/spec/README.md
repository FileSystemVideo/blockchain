<!--
order: 0
title: Evidence Overview
parent:
  title: "evidence"
-->

# `x/evidence`

## Table of Contents

<!-- TOC -->

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Messages](03_messages.md)**
4. **[Events](04_events.md)**
5. **[Params](05_params.md)**
6. **[BeginBlock](06_begin_block.md)**

## Abstract

`x/evidence` is an implementation of a Cosmos SDK module, per [ADR 009](./../../../docs/architecture/adr-009-evidence-module.md),
that allows for the submission and handling of arbitrary evidence of misbehavior such
as equivocation and counterfactual signing.


证据模块不同于标准的证据处理，标准的证据处理通常期望潜在的共识引擎（如Tendermint）在发现证据时自动提交证据，允许客户和国外连锁店直接提交更复杂的证据。


所有具体证据类型都必须实现“证据”接口契约。
提交的“证据”首先通过证据模块的“路由器”路由，在该路由器中，它尝试为特定的“证据”类型找到相应的已注册“处理程序”。
每个“证据”类型都必须在证据模块的keeper中注册一个“Handler”，才能成功地路由和执行它。


每个相应的处理程序还必须满足“handler”接口约定。

给定“证据”类型的“Handler”可以执行任何任意状态转换，如斜杠、监禁和墓碑。
