package wallet

import "os"

// define the wallet type
const (
	WALLET_BTC = 1
	WALLET_ETH = 2
)

// wallet interface
type Wallet interface {
	Address() string
	PrivateKey() string
	PublicKey() string
	FilePath() string
	Password() string
	GenerateKey() error
	SetPrivateKey(key string) error
	SaveKey() error
	LoadKey() error
	IsKeyFile(fi os.FileInfo) bool
}
