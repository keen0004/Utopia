package wallet

var (
	// wallet's create factory
	WalletMap = map[int]func(string, string) Wallet{
		WALLET_ETH: NewEthWallet,
		WALLET_BTC: NewBtcWallet,
	}
)

func NewWallet(wallettype int, path string, password string) Wallet {
	creator, ok := WalletMap[wallettype]
	if !ok {
		return nil
	}

	return creator(path, password)
}
