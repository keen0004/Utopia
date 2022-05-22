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

func (w *BtcWallet) PrivateKey() string {
	return ""
}

func (w *BtcWallet) PublicKey() string {
	return ""
}

func (w *BtcWallet) FilePath() string {
	return w.path
}

func (w *BtcWallet) Password() string {
	return w.password
}

func (w *BtcWallet) GenerateKey() error {
	return nil
}

func (w *BtcWallet) SetPrivateKey(key string) error {
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
