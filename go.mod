module fs.video/blockchain

go 1.15

require (
	fs.video/log v0.0.0
	fs.video/trerr v0.0.0
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cosmos/cosmos-sdk v0.42.1
	github.com/ethereum/go-ethereum v1.10.3
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/ipfs/go-blockservice v0.1.3 // indirect
	github.com/ipfs/go-cid v0.0.5
	github.com/ipfs/go-ipld-cbor v0.0.4 // indirect
	github.com/ipfs/go-ipld-format v0.2.0
	github.com/ipfs/go-log v1.0.4 // indirect
	github.com/ipfs/go-merkledag v0.3.2
	github.com/ipfs/go-path v0.0.7 // indirect
	github.com/ipfs/go-unixfs v0.2.4
	github.com/ipfs/interface-go-ipfs-core v0.2.7
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/kr/text v0.2.0 // indirect
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/multiformats/go-multibase v0.0.2 // indirect
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	github.com/tyler-smith/go-bip39 v1.1.0
	google.golang.org/genproto v0.0.0-20210207032614-bba0dbe2a9ea
	google.golang.org/grpc v1.35.0
)

replace (
	fs.video/log v0.0.0 => ./gomod/log
	fs.video/trerr v0.0.0 => ./gomod/trerr
	github.com/ipfs/go-merkledag v0.3.2 => ./../ipfs/vendor/github.com/ipfs/go-merkledag
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
	github.com/tendermint/tendermint v0.34.8 => ./gomod/tendermint@v0.34.8
	github.com/cosmos/cosmos-sdk v0.42.1 => ./gomod/cosmos-sdk@v0.42.1
)
