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

type Contract interface {
	Address() string
	Code() ([]byte, error)
	ABI() string
	SetABI(path string) error
	EncodeABI(method string, data string, withfunc bool) ([]byte, error)
	DecodeABI(method string, data []byte, withfunc bool) (string, error)
	Deploy(code []byte, params string, wallet wallet.Wallet, value *big.Int) (string, error)
	Call(params string, wallet wallet.Wallet, value *big.Int) ([]interface{}, error)
}
