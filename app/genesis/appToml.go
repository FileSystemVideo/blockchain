package genesis

const AppToml = `
minimum-gas-prices = "0.00005fsv"
pruning = "default"
pruning-keep-recent = "0"
pruning-keep-every = "0"
pruning-interval = "0"
halt-height = 0
halt-time = 0
min-retain-blocks = 0
inter-block-cache = true
index-events = []
[telemetry]
service-name = ""
enabled = false
enable-hostname = false
enable-hostname-label = false
enable-service-label = false
prometheus-retention-time = 0
global-labels = [
]

[api]
enable = true
swagger = false
address = "tcp://0.0.0.0:1317"
max-open-connections = 1000
rpc-read-timeout = 10
rpc-write-timeout = 0
rpc-max-body-bytes = 1000000
enabled-unsafe-cors = true

[grpc]
enable = true
address = "0.0.0.0:9090"

[state-sync]
snapshot-interval = 0
snapshot-keep-recent = 2

[evm]
tracer = ""
max-tx-gas-wanted = 500000

[json-rpc]
enable = true
address = "0.0.0.0:8545"
ws-address = "0.0.0.0:8546"
api = "eth,net,web3"
gas-cap = 25000000
evm-timeout = "5s"
txfee-cap = 1
filter-cap = 200
feehistory-cap = 100
logs-cap = 10000
block-range-cap = 10000
http-timeout = "30s"
http-idle-timeout = "2m0s"

[tls]
certificate-path = ""
key-path = ""
`
