package chain

var (
	ChainMap = map[uint64]func(uint64, string, string) Chain{
		ETH_MAINNET: NewEthChain,
		BSC_MAINNET: NewEthChain,
		DEV_NETWORK: NewEthChain,
	}
)

func NewChain(id uint64, currency string, name string) Chain {
	creator, ok := ChainMap[id]
	if !ok {
		return nil
	}

	chain := creator(id, currency, name)
	if chain == nil {
		return nil
	}

	err := chain.Connect([]string{}, true)
	if err != nil {
		return nil
	}

	return chain
}
