package contract

import (
	"errors"
	"math/big"
	"utopia/contracts/token"
	"utopia/internal/chain"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/common"
)

type ERC20Contract struct {
	chain    chain.Chain    // Chain id which contract deployed
	address  common.Address // Contract address
	contract *token.ERC20
}

func NewERC20(chain chain.Chain, address string) Contract {
	return &ERC20Contract{
		chain:    chain,
		address:  common.HexToAddress(address),
		contract: nil,
	}
}

func (c *ERC20Contract) Address() string {
	return c.address.Hex()
}

func (c *ERC20Contract) Code() (string, error) {
	return c.chain.Code(c.address.Hex())
}

func (c *ERC20Contract) ABI() string {
	return token.ERC20ABI
}

// not support under functions
func (c *ERC20Contract) SetABI(path string) error {
	return errors.New("Not support")
}

func (c *ERC20Contract) EncodeABI(method string, data string, withfunc bool) (string, error) {
	return "", errors.New("Not support")
}

func (c *ERC20Contract) DecodeABI(method string, data string, withfunc bool) (string, error) {
	return "", errors.New("Not support")
}

func (c *ERC20Contract) Deploy(code string, params string, wallet wallet.Wallet, value *big.Int) (string, error) {
	return "", errors.New("Not support")
}

func (c *ERC20Contract) Call(params string, wallet wallet.Wallet, value *big.Int) ([]interface{}, error) {
	return nil, errors.New("Not support")
}

// query token balance of owner
func (c *ERC20Contract) Balance(address string) (*big.Int, error) {
	return c.contract.BalanceOf(nil, common.HexToAddress(address))
}

// transfer token to receiver
func (c *ERC20Contract) Transfer(to string, value *big.Int, wallet wallet.Wallet) (string, error) {
	opts, err := c.chain.(*chain.EthChain).GenTransOpts(wallet, nil)
	if err != nil {
		return "", err
	}

	tx, err := c.contract.Transfer(opts, common.HexToAddress(to), value)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

// approve token to receiver
func (c *ERC20Contract) Approve(to string, value *big.Int, wallet wallet.Wallet) (string, error) {
	opts, err := c.chain.(*chain.EthChain).GenTransOpts(wallet, nil)
	if err != nil {
		return "", err
	}

	tx, err := c.contract.Approve(opts, common.HexToAddress(to), value)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
