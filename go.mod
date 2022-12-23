module fs.video/blockchain

go 1.15

require (
	fs.video/log v0.0.0
	fs.video/trerr v0.0.0
	github.com/StackExchange/wmi v1.2.1
	github.com/btcsuite/btcd v0.22.1
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.45.4
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/ibc-go/v3 v3.0.0
	github.com/ethereum/go-ethereum v1.10.16
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
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
	github.com/json-iterator/go v1.1.12
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/multiformats/go-multibase v0.0.2 // indirect
	github.com/otiai10/copy v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.1
	github.com/tendermint/tendermint v0.34.19
	github.com/tendermint/tm-db v0.6.7
	github.com/tharsis/ethermint v0.15.0
	github.com/tharsis/evmos/v4 v4.0.1
	github.com/tyler-smith/go-bip39 v1.1.0
	go.opencensus.io v0.23.0
	google.golang.org/genproto v0.0.0-20220429170224-98d788798c3e
	google.golang.org/grpc v1.45.0
)

replace (
	fs.video/log v0.0.0 => ./gomod/log
	fs.video/trerr v0.0.0 => ./gomod/trerr
	github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76
	github.com/cosmos/cosmos-sdk v0.45.4 => ./gomod/cosmos-sdk@v0.45.4
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	//github.com/ipfs/go-merkledag v0.3.2 => ./../ipfs-libs/github.com/ipfs/go-ipfs@0.5.0/vendor/github.com/ipfs/go-merkledag
	github.com/tendermint/tendermint v0.34.19 => ./gomod/tendermint@v0.34.19
	github.com/tharsis/ethermint v0.15.0 => ./gomod/ethermint@v0.15.0
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
