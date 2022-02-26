---
order: 13
---

# 轻客户端

对于大多数应用程序而言，轻客户端是完整区块链系统的重要组成部分。Tendermint为轻型客户端应用程序提供独特的速度和安全属性。.

See our [light
package](https://pkg.go.dev/github.com/tendermint/tendermint/light?tab=doc).

## Overview

轻型客户端协议的目标是为最近的块散列获取提交，其中提交包括来自最后一个已知验证器集的大多数签名。
从那里，所有的应用程序状态都可以用[merkle证明](https://github.com/tendermint/spec/blob/953523c3cb99fdb8c8f7a2d21e3a99094279e9de/spec/blockchain/encoding.md#iavl-tree) 进行验证 .

## Properties

- 您可以获得Tendermint的全部担保安全利益；无需等待确认。
- 您可以享受Tendermint的全速优势；事务立即提交。
- 您可以以非交互方式获取应用程序状态的最新版本（无需向区块链提交任何内容）。
例如，这意味着您可以从名称注册表获取名称的最新值，而无需担心fork审查攻击，也无需发布提交和等待确认。
它快速、安全、免费！

## Where to obtain trusted height & hash

[Trust Options](https://pkg.go.dev/github.com/tendermint/tendermint/light?tab=doc#TrustOptions)

获取半可信哈希和高度的一种方法是查询多个完整节点并比较它们的哈希：

```bash
$ curl -s https://233.123.0.140:26657:26657/commit | jq "{height: .result.signed_header.header.height, hash: .result.signed_header.commit.block_id.hash}"
{
  "height": "273",
  "hash": "188F4F36CBCD2C91B57509BBF231C777E79B52EE3E0D90D06B1A25EB16E6E23D"
}
```

## 将轻型客户端作为HTTP代理服务器运行

Tendermint comes with a built-in `tendermint light` command, which can be used
to run a light client proxy server, verifying Tendermint RPC. All calls that
can be tracked back to a block header by a proof will be verified before
passing them back to the caller. Other than that, it will present the same
interface as a full Tendermint node.

You can start the light client proxy server by running `tendermint light <chainID>`,
with a variety of flags to specify the primary node,  the witness nodes (which cross-check
the information provided by the primary), the hash and height of the trusted header,
and more.

For example:

```bash
$ tendermint light supernova -p tcp://233.123.0.140:26657 \
  -w tcp://179.63.29.15:26657,tcp://144.165.223.135:26657 \
  --height=10 --hash=37E9A6DD3FA25E83B22C18835401E8E56088D0D7ABC6FD99FCDC920DD76C1C57
```

For additional options, run `tendermint light --help`.
