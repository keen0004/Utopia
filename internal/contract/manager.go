package contract

import "utopia/internal/chain"

var (
	ContractMap = map[int]func(chain.Chain, string) Contract{
		// COMMON_CRONTACT: NewCommonContract,
		// ERC20_CONTRACT:  NewERC20,
		// ERC721_CONTRACT: NewERC721,
	}
)

func NewContract(chain chain.Chain, address string, ctype int) Contract {
	creator, ok := ContractMap[ctype]
	if !ok {
		return nil
	}

	return creator(chain, address)
}
