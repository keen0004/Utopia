package chain

import (
	"context"
	"errors"
	"math/big"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthChain struct {
	Id        uint64            // The id of chain
	Currency  string            // Symbol of chain currency
	Name      string            // Name of chain
	Rpc       []string          // List of rpc server
	Client    *ethclient.Client // Connection of chain
	index     int               // Connect index of server list
	connected bool              // Is connect to server
}

func NewEthChain(id uint64, currency string, name string) Chain {
	server := make([]string, 0)
	meta, err := ChainMetaById(id)
	if err == nil {
		server = meta.RpcServer
	}

	return &EthChain{
		Id:        id,
		Currency:  currency,
		Name:      name,
		Rpc:       server,
		Client:    nil,
		index:     0,
		connected: false,
	}
}

func (chain *EthChain) Connect(server []string, checkid bool) error {
	if len(server) == 0 {
		server = chain.Rpc
		if len(server) == 0 {
			return errors.New("RPC server list is empty")
		}
	}

	// reset rpc server list and timeout
	chain.DisConnect()
	chain.Rpc = server
	chain.index = 0

	// refresh connect
	if err := chain.refresh(); err != nil {
		return err
	}

	// check chain id
	if checkid {
		id, err := chain.ChainId()
		if err != nil {
			return err
		}

		if id.Cmp(big.NewInt(int64(chain.Id))) != 0 {
			return errors.New("Not match chain id")
		}
	}

	return nil
}

func (chain *EthChain) DisConnect() {
	if chain.connected {
		chain.Client.Close()
	}

	chain.Client = nil
	chain.connected = false
}

func (chain *EthChain) ChainId() (big.Int, error) {
	chain.refresh()
	if !chain.connected {
		return *big.NewInt(0), errors.New("Chain not connected")
	}

	id, err := chain.Client.ChainID(context.Background())
	if err != nil {
		return *big.NewInt(0), err
	}

	return *id, nil
}

func (chain *EthChain) GasPrice() (big.Int, error) {
	chain.refresh()
	if !chain.connected {
		return *big.NewInt(0), errors.New("Chain not connected")
	}

	gas, err := chain.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return *big.NewInt(0), err
	}

	return *gas, nil
}

func (chain *EthChain) BlockNumber() (uint64, error) {
	chain.refresh()
	if !chain.connected {
		return 0, errors.New("Chain not connected")
	}

	number, err := chain.Client.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}

	return number, nil
}

func (chain *EthChain) BlockByNumber(number uint64) (*types.Block, error) {
	chain.refresh()
	if !chain.connected {
		return nil, errors.New("Chain not connected")
	}

	return chain.Client.BlockByNumber(context.Background(), big.NewInt(int64(number)))
}

func (chain *EthChain) BlockByHash(hash []byte) (*types.Block, error) {
	chain.refresh()
	if !chain.connected {
		return nil, errors.New("Chain not connected")
	}

	return chain.Client.BlockByHash(context.Background(), common.BytesToHash(hash))
}

func (chain *EthChain) Transaction(hash []byte) (*types.Transaction, bool, error) {
	chain.refresh()
	if !chain.connected {
		return nil, false, errors.New("Chain not connected")
	}

	return chain.Client.TransactionByHash(context.Background(), common.BytesToHash(hash))
}

func (chain *EthChain) Receipt(hash []byte) (*types.Receipt, error) {
	chain.refresh()
	if !chain.connected {
		return nil, errors.New("Chain not connected")
	}

	return chain.Client.TransactionReceipt(context.Background(), common.BytesToHash(hash))
}

func (chain *EthChain) SendTransaction(tx *types.Transaction, wallet wallet.Wallet) error {
	chain.refresh()
	if !chain.connected {
		return errors.New("Chain not connected")
	}

	key, err := crypto.ToECDSA(wallet.PrivateKey())
	if err != nil {
		return err
	}

	signTx, err := types.SignTx(tx, types.NewLondonSigner(big.NewInt(int64(chain.Id))), key)
	if err != nil {
		return err
	}

	return chain.Client.SendTransaction(context.Background(), signTx)
}

func (chain *EthChain) EstimateGas(tx *types.Transaction) (uint64, error) {
	chain.refresh()
	if !chain.connected {
		return 0, errors.New("Chain not connected")
	}

	from, _ := types.Sender(types.NewLondonSigner(big.NewInt(int64(chain.Id))), tx)
	return chain.Client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:       from,
		To:         tx.To(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		GasFeeCap:  tx.GasFeeCap(),
		GasTipCap:  tx.GasTipCap(),
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	})
}

func (chain *EthChain) Balance(address string) (big.Int, error) {
	chain.refresh()
	if !chain.connected {
		return *big.NewInt(0), errors.New("Chain not connected")
	}

	balance, err := chain.Client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	return *balance, err
}

func (chain *EthChain) Transfer(to []string, value []*big.Int, wallet wallet.Wallet) error {
	if len(to) != len(value) {
		return errors.New("No matched address and values")
	}

	chain.refresh()
	if !chain.connected {
		return errors.New("Chain not connected")
	}

	nonce, err := chain.Nonce(wallet.Address())
	if err != nil {
		return err
	}

	gasprice, err := chain.GasPrice()
	if err != nil {
		return err
	}

	balance, err := chain.Balance(wallet.Address())
	if err != nil {
		return err
	}

	total := big.NewInt(0)
	for _, v := range value {
		total = new(big.Int).Add(total, v)
		total = new(big.Int).Add(total, new(big.Int).Mul(&gasprice, big.NewInt(21000)))
	}

	if balance.Cmp(total) < 0 {
		return errors.New("Not enough balance")
	}

	key, err := crypto.ToECDSA(wallet.PrivateKey())
	if err != nil {
		return err
	}

	for index, address := range to {
		if address == wallet.Address() {
			continue
		}

		if value[index].Cmp(big.NewInt(0)) == 0 {
			continue
		}

		tx := types.NewTransaction(nonce, common.HexToAddress(address), value[index], 21000, &gasprice, nil)
		signTx, err := types.SignTx(tx, types.NewLondonSigner(big.NewInt(int64(chain.Id))), key)
		if err != nil {
			return err
		}

		err = chain.Client.SendTransaction(context.Background(), signTx)
		if err != nil {
			return err
		}

		nonce++
	}

	return nil
}

func (chain *EthChain) Nonce(address string) (uint64, error) {
	chain.refresh()
	if !chain.connected {
		return 0, errors.New("Chain not connected")
	}

	return chain.Client.PendingNonceAt(context.Background(), common.HexToAddress(address))
}

func (chain *EthChain) Code(address string) ([]byte, error) {
	chain.refresh()
	if !chain.connected {
		return []byte{}, errors.New("Chain not connected")
	}

	return chain.Client.CodeAt(context.Background(), common.HexToAddress(address), nil)
}

func (chain *EthChain) refresh() error {
	if chain.connected {
		return nil
	}

	for i := chain.index; i < len(chain.Rpc); i++ {
		client, err := ethclient.Dial(chain.Rpc[i])
		if err != nil {
			continue
		}

		chain.Client = client
		chain.index = i
		chain.connected = true

		return nil
	}

	for i := 0; i < len(chain.Rpc) && i < chain.index; i++ {
		client, err := ethclient.Dial(chain.Rpc[i])
		if err != nil {
			continue
		}

		chain.Client = client
		chain.index = i
		chain.connected = true

		return nil
	}

	return errors.New("Can not connect any server")
}

func (c *EthChain) GenTransOpts(wallet wallet.Wallet, value *big.Int) (*bind.TransactOpts, error) {
	key, err := crypto.ToECDSA(wallet.PrivateKey())
	if err != nil {
		return nil, err
	}
	opts, _ := bind.NewKeyedTransactorWithChainID(key, new(big.Int).SetUint64(c.Id))

	opts.GasPrice, err = c.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	opts.Value = value
	return opts, nil
}
