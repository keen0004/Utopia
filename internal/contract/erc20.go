package contract

import (
	"errors"
	"math/big"
	"utopia/contracts/token"
	"utopia/internal/chain"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type ERC20Contract struct {
	chain    chain.Chain    // Chain id which contract deployed
	address  common.Address // Contract address
	contract *token.ERC20
}

func (c *ERC20Contract) Address() string {
	return c.address.Hex()
}

func (c *ERC20Contract) Code() ([]byte, error) {
	return c.chain.Code(c.address.Hex())
}

func (c *ERC20Contract) ABI() string {
	return token.ERC20ABI
}

func (c *ERC20Contract) SetABI(path string) error {
	return errors.New("Not support")
}

func (c *ERC20Contract) Call(params string) error {
	return errors.New("Not support")
}

func (c *ERC20Contract) Balance(address string) (*big.Int, error) {
	return c.contract.BalanceOf(nil, common.HexToAddress(address))
}

func (c *ERC20Contract) Transfer(to string, value *big.Int, wallet wallet.Wallet) (string, error) {
	key, err := crypto.ToECDSA(wallet.PrivateKey())
	if err != nil {
		return "", err
	}

	tx, err := c.contract.Transfer(bind.NewKeyedTransactor(key), common.HexToAddress(to), value)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (c *ERC20Contract) Approve(to string, value *big.Int, wallet wallet.Wallet) (string, error) {
	key, err := crypto.ToECDSA(wallet.PrivateKey())
	if err != nil {
		return "", err
	}

	tx, err := c.contract.Approve(bind.NewKeyedTransactor(key), common.HexToAddress(to), value)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
