package contract

import (
	"errors"
	"math/big"
	"utopia/contracts/token"
	"utopia/internal/chain"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/common"
)

type ERC721Contract struct {
	chain    chain.Chain    // Chain id which contract deployed
	address  common.Address // Contract address
	contract *token.ERC721
}

type ERC721Attr struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

type ERC721Meta struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Image       string       `json:"image"`
	Attributes  []ERC721Attr `json:"attributes"`
}

func NewERC721(chain chain.Chain, address string) Contract {
	return &ERC721Contract{
		chain:    chain,
		address:  common.HexToAddress(address),
		contract: nil,
	}
}

func (c *ERC721Contract) Address() string {
	return c.address.Hex()
}

func (c *ERC721Contract) Code() (string, error) {
	return c.chain.Code(c.address.Hex())
}

func (c *ERC721Contract) ABI() string {
	return token.ERC20ABI
}

// not support under functions
func (c *ERC721Contract) SetABI(path string) error {
	return errors.New("Not support")
}

func (c *ERC721Contract) EncodeABI(method string, data string, withfunc bool) (string, error) {
	return "", errors.New("Not support")
}

func (c *ERC721Contract) DecodeABI(method string, data string, withfunc bool) (string, error) {
	return "", errors.New("Not support")
}

func (c *ERC721Contract) Deploy(code string, params string, wallet wallet.Wallet, value *big.Int) (string, error) {
	return "", errors.New("Not support")
}

func (c *ERC721Contract) Call(params string, wallet wallet.Wallet, value *big.Int) ([]interface{}, error) {
	return nil, errors.New("Not support")
}

// query token number which owned by address
func (c *ERC721Contract) Balance(address string) (uint64, error) {
	balance, err := c.contract.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		return 0, err
	}

	return balance.Uint64(), nil
}

// this function need enumable 721 contract
func (c *ERC721Contract) TokenIdByIndex(address string, index uint32) (uint64, error) {
	return 0, nil
}

// query owner of token
func (c *ERC721Contract) Owner(tokenid uint64) (string, error) {
	address, err := c.contract.OwnerOf(nil, new(big.Int).SetUint64(tokenid))
	if err != nil {
		return "", err
	}

	return address.Hex(), nil
}

// query token url
func (c *ERC721Contract) TokenUrl(tokenid uint64) (string, error) {
	return c.contract.TokenURI(nil, new(big.Int).SetUint64(tokenid))
}

// transfer token from owner to receiver
func (c *ERC721Contract) Transfer(to string, tokenid uint64, wallet wallet.Wallet) (string, error) {
	opts, err := c.chain.(*chain.EthChain).GenTransOpts(wallet, nil)
	if err != nil {
		return "", err
	}

	tx, err := c.contract.TransferFrom(opts, common.HexToAddress(wallet.Address()), common.HexToAddress(to), new(big.Int).SetUint64(tokenid))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

// approve token to receiver, tokenid 0 means all tokens
func (c *ERC721Contract) Approve(to string, tokenid uint64, approve bool, wallet wallet.Wallet) (string, error) {
	opts, err := c.chain.(*chain.EthChain).GenTransOpts(wallet, nil)
	if err != nil {
		return "", err
	}

	if tokenid == 0 {
		tx, err := c.contract.SetApprovalForAll(opts, common.HexToAddress(to), approve)
		if err != nil {
			return "", err
		}

		return tx.Hash().Hex(), nil
	} else {
		tx, err := c.contract.Approve(opts, common.HexToAddress(to), new(big.Int).SetUint64(tokenid))
		if err != nil {
			return "", err
		}

		return tx.Hash().Hex(), nil
	}
}
