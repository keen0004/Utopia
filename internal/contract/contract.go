package contract

import (
	"math/big"
	"utopia/internal/wallet"
)

const (
	COMMON_CRONTACT = 1
	ERC20_CONTRACT  = 2
	ERC721_CONTRACT = 3
)

// contract interface
type Contract interface {
	Address() string
	Code() (string, error)
	ABI() string
	SetABI(path string) error
	EncodeABI(method string, data string, withfunc bool) (string, error)
	DecodeABI(method string, data string, withfunc bool) (string, error)
	Deploy(code string, params string, wallet wallet.Wallet, value *big.Int) (string, error)
	Call(params string, wallet wallet.Wallet, value *big.Int) ([]interface{}, error)
}
