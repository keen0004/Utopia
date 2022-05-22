package chain

import (
	"math/big"
	"time"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/core/types"
)

type BtcChain struct {
	Id        uint64        // The id of chain
	Currency  string        // Symbol of chain currency
	Name      string        // Name of chain
	Rpc       []string      // List of rpc server
	Timeout   time.Duration // Call timeout in millsecond
	index     int           // Connect index of server list
	connected bool          // Is connect to server
}

func NewBtcChain(id uint64, currency string, name string) Chain {
	return &BtcChain{
		Id:        id,
		Currency:  currency,
		Name:      name,
		Rpc:       make([]string, 0),
		index:     0,
		connected: false,
	}
}

func (chain *BtcChain) Connect(server []string, checkid bool) error {
	return nil
}

func (chain *BtcChain) DisConnect() {

}

func (chain *BtcChain) ChainId() (*big.Int, error) {
	return big.NewInt(0), nil
}

func (chain *BtcChain) GasPrice() (*big.Int, error) {
	return big.NewInt(0), nil
}

func (chain *BtcChain) BlockNumber() (uint64, error) {
	return 0, nil
}

func (chain *BtcChain) BlockByNumber(number uint64) (*types.Block, error) {
	return nil, nil
}

func (chain *BtcChain) BlockByHash(hash []byte) (*types.Block, error) {
	return nil, nil
}

func (chain *BtcChain) Transaction(hash []byte) (*types.Transaction, bool, error) {
	return nil, false, nil
}

func (chain *BtcChain) Receipt(hash []byte) (*types.Receipt, error) {
	return nil, nil
}

func (chain *BtcChain) SendTransaction(tx *types.Transaction, wallet wallet.Wallet) (string, error) {
	return "", nil
}

func (chain *BtcChain) EstimateGas(tx *types.Transaction) (uint64, error) {
	return 0, nil
}

func (chain *BtcChain) Balance(address string) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (chain *BtcChain) Transfer(to string, value *big.Int, wallet wallet.Wallet) (string, error) {
	return "", nil
}

func (chain *BtcChain) Nonce(address string) (uint64, error) {
	return 0, nil
}

func (chain *BtcChain) Code(address string) (string, error) {
	return "", nil
}
