package contract

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
	Call(params string) error
}
