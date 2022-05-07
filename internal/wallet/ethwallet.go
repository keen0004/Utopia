package wallet

import (
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type EthWallet struct {
	path     string
	password string
	key      *keystore.Key
}

func NewEthWallet(path string, password string) Wallet {
	return &EthWallet{path: path, password: password, key: nil}
}

func (w *EthWallet) Address() string {
	if w.key == nil {
		return common.BigToAddress(big.NewInt(0)).Hex()
	}

	return w.key.Address.Hex()
}

func (w *EthWallet) PrivateKey() []byte {
	if w.key == nil {
		return []byte{}
	}

	return crypto.FromECDSA(w.key.PrivateKey)
}

func (w *EthWallet) PublicKey() []byte {
	if w.key == nil {
		return nil
	}

	return crypto.FromECDSAPub(&w.key.PrivateKey.PublicKey)
}

func (w *EthWallet) GenerateKey() error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	UUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	w.key = &keystore.Key{
		Id:         UUID,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	return nil
}

func (w *EthWallet) SaveKey() error {
	if w.path == "" {
		return errors.New("Not set the key file path")
	}

	if w.key == nil {
		err := w.GenerateKey()
		if err != nil {
			return err
		}
	}

	keyjson, err := keystore.EncryptKey(w.key, w.password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(w.path, keyjson, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (w *EthWallet) LoadKey() error {
	keyjson, err := ioutil.ReadFile(w.path)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyjson, w.password)
	if err != nil {
		return err
	}

	w.key = key
	return nil
}

func (w *EthWallet) IsKeyFile(fi os.FileInfo) bool {
	// Skip editor backups and UNIX-style hidden files.
	if strings.HasSuffix(fi.Name(), "~") || strings.HasPrefix(fi.Name(), ".") {
		return false
	}

	// Skip misc special files, directories (yes, symlinks too).
	if fi.IsDir() || fi.Mode()&os.ModeType != 0 {
		return false
	}

	return true
}
