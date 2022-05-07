package wallet

import "os"

type BtcWallet struct {
	path     string
	password string
}

func NewBtcWallet(path string, password string) Wallet {
	return &BtcWallet{path: path, password: password}
}

func (w *BtcWallet) Address() string {
	return ""
}

func (w *BtcWallet) PrivateKey() []byte {
	return []byte{}
}

func (w *BtcWallet) PublicKey() []byte {
	return []byte{}
}

func (w *BtcWallet) GenerateKey() error {
	return nil
}

func (w *BtcWallet) SaveKey() error {
	return nil
}

func (w *BtcWallet) LoadKey() error {
	return nil
}

func (w *BtcWallet) IsKeyFile(fi os.FileInfo) bool {
	return false
}
