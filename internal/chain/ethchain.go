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
	meta, err := ChainMetaById(id)
	if err != nil {
		return nil
	}

	return &EthChain{
		Id:        id,
		Currency:  currency,
		Name:      name,
		Rpc:       meta.RpcServer,
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
	err := chain.refresh()
	if err != nil {
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

func (chain *EthChain) ChainId() (*big.Int, error) {
	chain.refresh()
	if !chain.connected {
		return nil, errors.New("Chain not connected")
	}

	// query chain id from server
	id, err := chain.Client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (chain *EthChain) GasPrice() (*big.Int, error) {
	chain.refresh()
	if !chain.connected {
		return nil, errors.New("Chain not connected")
	}

	// query gas price from chain
	gas, err := chain.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	return gas, nil
}

func (chain *EthChain) BlockNumber() (uint64, error) {
	chain.refresh()
	if !chain.connected {
		return 0, errors.New("Chain not connected")
	}

	// query block number from chain
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

	// query block data from chain
	return chain.Client.BlockByNumber(context.Background(), big.NewInt(int64(number)))
}

func (chain *EthChain) BlockByHash(hash []byte) (*types.Block, error) {
	chain.refresh()
	if !chain.connected {
		return nil, errors.New("Chain not connected")
	}

	// query block data from chain
	return chain.Client.BlockByHash(context.Background(), common.BytesToHash(hash))
}

func (chain *EthChain) Transaction(hash []byte) (*types.Transaction, bool, error) {
	chain.refresh()
	if !chain.connected {
		return nil, false, errors.New("Chain not connected")
	}

	// query transaction from chain
	return chain.Client.TransactionByHash(context.Background(), common.BytesToHash(hash))
}

func (chain *EthChain) Receipt(hash []byte) (*types.Receipt, error) {
	chain.refresh()
	if !chain.connected {
		return nil, errors.New("Chain not connected")
	}

	// query transaction receipt from chain
	return chain.Client.TransactionReceipt(context.Background(), common.BytesToHash(hash))
}

func (chain *EthChain) SendTransaction(tx *types.Transaction, wallet wallet.Wallet) (string, error) {
	chain.refresh()
	if !chain.connected {
		return "", errors.New("Chain not connected")
	}

	// sign transaction and send
	key, err := crypto.ToECDSA(common.FromHex(wallet.PrivateKey()))
	if err != nil {
		return "", err
	}

	signTx, err := types.SignTx(tx, types.NewLondonSigner(big.NewInt(int64(chain.Id))), key)
	if err != nil {
		return "", err
	}

	return signTx.Hash().Hex(), chain.Client.SendTransaction(context.Background(), signTx)
}

func (chain *EthChain) EstimateGas(tx *types.Transaction) (uint64, error) {
	chain.refresh()
	if !chain.connected {
		return 0, errors.New("Chain not connected")
	}

	// estimate gas cost from chain
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

func (chain *EthChain) Balance(address string) (*big.Int, error) {
	chain.refresh()
	if !chain.connected {
		return nil, errors.New("Chain not connected")
	}

	// query acount balance from chain
	balance, err := chain.Client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	return balance, err
}

// batch transfer value
func (chain *EthChain) Transfer(to string, value *big.Int, wallet wallet.Wallet) (string, error) {
	chain.refresh()
	if !chain.connected {
		return "", errors.New("Chain not connected")
	}

	if to == wallet.Address() {
		return "", errors.New("Can not transfer value to self")
	}

	// get sender nonce from chain
	nonce, err := chain.Nonce(wallet.Address())
	if err != nil {
		return "", err
	}

	gasprice, err := chain.GasPrice()
	if err != nil {
		return "", err
	}

	// check balance is enough
	balance, err := chain.Balance(wallet.Address())
	if err != nil {
		return "", err
	}

	if balance.Cmp(value) < 0 {
		return "", errors.New("Not enough balance")
	}

	key, err := crypto.ToECDSA(common.FromHex(wallet.PrivateKey()))
	if err != nil {
		return "", err
	}

	// gen transaction and sign it
	tx := types.NewTransaction(nonce, common.HexToAddress(to), value, 21000, gasprice, nil)
	signTx, err := types.SignTx(tx, types.NewLondonSigner(big.NewInt(int64(chain.Id))), key)
	if err != nil {
		return "", err
	}

	// send transaction
	err = chain.Client.SendTransaction(context.Background(), signTx)
	if err != nil {
		return "", err
	}

	return signTx.Hash().Hex(), nil
}

func (chain *EthChain) Nonce(address string) (uint64, error) {
	chain.refresh()
	if !chain.connected {
		return 0, errors.New("Chain not connected")
	}

	// read account nonce from chain
	return chain.Client.PendingNonceAt(context.Background(), common.HexToAddress(address))
}

func (chain *EthChain) Code(address string) (string, error) {
	chain.refresh()
	if !chain.connected {
		return "", errors.New("Chain not connected")
	}

	// read contract code from chain
	code, err := chain.Client.CodeAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return "", err
	}

	return "0x" + common.Bytes2Hex(code), nil
}

func (chain *EthChain) refresh() error {
	if chain.connected {
		return nil
	}

	// connect rpc server by index
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
	// set sign function
	key, err := crypto.ToECDSA(common.FromHex(wallet.PrivateKey()))
	if err != nil {
		return nil, err
	}
	opts, _ := bind.NewKeyedTransactorWithChainID(key, new(big.Int).SetUint64(c.Id))

	// set gasprice
	opts.GasPrice, err = c.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// set transfer value
	opts.Value = value
	return opts, nil
}
