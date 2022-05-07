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

type ERC721Contract struct {
	chain    chain.Chain    // Chain id which contract deployed
	address  common.Address // Contract address
	contract *token.ERC721
}

func (c *ERC721Contract) Address() string {
	return c.address.Hex()
}

func (c *ERC721Contract) Code() ([]byte, error) {
	return c.chain.Code(c.address.Hex())
}

func (c *ERC721Contract) ABI() string {
	return token.ERC20ABI
}

func (c *ERC721Contract) SetABI(path string) error {
	return errors.New("Not support")
}

func (c *ERC721Contract) Call(params string) error {
	return errors.New("Not support")
}

func (c *ERC721Contract) Balance(address string) (uint64, error) {
	balance, err := c.contract.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		return 0, err
	}

	return balance.Uint64(), nil
}

func (c *ERC721Contract) TokenIdByIndex(address string, index uint32) (uint64, error) {
	return 0, nil
}

func (c *ERC721Contract) Owner(tokenid uint64) (string, error) {
	address, err := c.contract.OwnerOf(nil, new(big.Int).SetUint64(tokenid))
	if err != nil {
		return "", err
	}

	return address.Hex(), nil
}

func (c *ERC721Contract) TokenUrl(tokenid uint64) (string, error) {
	return c.contract.TokenURI(nil, new(big.Int).SetUint64(tokenid))
}

func (c *ERC721Contract) Transfer(to string, tokenid uint64, wallet wallet.Wallet) (string, error) {
	key, err := crypto.ToECDSA(wallet.PrivateKey())
	if err != nil {
		return "", err
	}

	tx, err := c.contract.TransferFrom(bind.NewKeyedTransactor(key), common.HexToAddress(wallet.Address()), common.HexToAddress(to), new(big.Int).SetUint64(tokenid))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (c *ERC721Contract) Approve(to string, tokenid uint64, approve bool, wallet wallet.Wallet) (string, error) {
	key, err := crypto.ToECDSA(wallet.PrivateKey())
	if err != nil {
		return "", err
	}

	if tokenid == 0 {
		tx, err := c.contract.SetApprovalForAll(bind.NewKeyedTransactor(key), common.HexToAddress(to), approve)
		if err != nil {
			return "", err
		}

		return tx.Hash().Hex(), nil
	} else {
		tx, err := c.contract.Approve(bind.NewKeyedTransactor(key), common.HexToAddress(to), new(big.Int).SetUint64(tokenid))
		if err != nil {
			return "", err
		}

		return tx.Hash().Hex(), nil
	}
}
