package genesis

import "fs.video/blockchain/core"

var ClientToml = `
chain-id = "` + core.ChainID + `"
keyring-backend = "os"
output = "text"
node = "tcp://localhost:26657"
broadcast-mode = "sync"
`
