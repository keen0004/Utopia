package wallet

import "os"

// define the wallet type
const (
	WALLET_BTC = 1
	WALLET_ETH = 2
)

type Wallet interface {
	Address() string
	PrivateKey() []byte
	PublicKey() []byte
	GenerateKey() error
	SaveKey() error
	LoadKey() error
	IsKeyFile(fi os.FileInfo) bool
}
