package chain

import (
	"math/big"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/core/types"
)

// define the chain id
const (
	ETH_MAINNET      = 1
	BSC_MAINNET      = 56
	POLYGON_MAINNET  = 137
	AVAX_MAINNET     = 43114
	FTM_MAINNET      = 250
	ARBITRUM_MAINNET = 42161
	OPTIMISM_MAINNET = 10
	HECO_MAINNET     = 128
	TLOS_MAINNET     = 40
	DEV_NETWORK      = 12345
)

type Chain interface {
	Connect(rpc []string, checkid bool) error
	DisConnect()

	ChainId() (big.Int, error)
	GasPrice() (big.Int, error)

	BlockNumber() (uint64, error)
	BlockByNumber(number uint64) (*types.Block, error)
	BlockByHash(hash []byte) (*types.Block, error)

	Transaction(hash []byte) (*types.Transaction, error)
	PendingTransaction() (types.Transactions, error)
	Receipt(hash []byte) (*types.Receipt, error)
	SendTransaction(tx *types.Transaction, wallet wallet.Wallet) error
	EstimateGas(tx *types.Transaction) (uint64, error)

	Balance(address string) (big.Int, error)
	Transfer(to []string, value []*big.Int, wallet wallet.Wallet) error
	Nonce(address string) (uint64, error)
	Code(address string) ([]byte, error)
}
