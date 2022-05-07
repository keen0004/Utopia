package wallet

import (
	"io/ioutil"
	"path/filepath"
)

var (
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

func ListWallet(wallettype int, keydir string, password string) ([]Wallet, error) {
	files, err := ioutil.ReadDir(keydir)
	if err != nil {
		return nil, err
	}

	wallet := make([]Wallet, 0, len(files))
	for _, fi := range files {
		w := NewWallet(wallettype, filepath.Join(keydir, fi.Name()), password)
		if !w.IsKeyFile(fi) {
			continue
		}

		err := w.LoadKey()
		if err != nil {
			continue
		}

		wallet = append(wallet, w)
	}

	return wallet, nil
}
