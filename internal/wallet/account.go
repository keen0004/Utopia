package wallet

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"utopia/internal/excel"
)

var (
	ACCOUNTS_SHEET_NAME = "accounts"
	ACCOUNT_LIST_HEADER = []string{"index", "address", "key", "keystore", "password", "notes"}
	AccountList         = make(map[string]Wallet)
)

func GetWallet(address string) (Wallet, error) {
	w, ok := AccountList[address]
	if !ok {
		return nil, errors.New("Address is not exist")
	}

	return w, nil
}

func AddWallet(address string, wallet Wallet) error {
	_, ok := AccountList[address]
	if ok {
		return errors.New("Address is exist")
	}

	AccountList[address] = wallet
	return nil
}

func LoadAccountList(path string) error {
	// read row data from excel file in accounts sheet
	file, err := excel.NewExcel(path)
	if err != nil {
		return err
	}

	err = file.Open()
	if err != nil {
		return err
	}
	defer file.Close(false)

	data, err := file.ReadAll(ACCOUNTS_SHEET_NAME)
	if err != nil {
		return err
	}

	// parse data to wallet
	for index, row := range data {
		// skip the header
		if index == 0 {
			continue
		}

		// [index, address, key, keystore, password, notes]
		if len(row) < len(ACCOUNT_LIST_HEADER) {
			return errors.New("Invalid file format")
		}

		// create wallet and load private key
		w := NewWallet(WALLET_ETH, row[3], row[4])
		if row[2] != "" {
			w.SetPrivateKey(row[2])
		} else {
			err = w.LoadKey()
			if err != nil {
				return err
			}
		}

		err = AddWallet(row[1], w)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveAccountList(path string) error {
	file, err := excel.NewExcel(path)
	if err != nil {
		return err
	}

	err = file.Open()
	if err != nil {
		return err
	}
	defer file.Close(true)

	// add excel header data
	index := 1
	data := make([][]string, 0)
	data = append(data, ACCOUNT_LIST_HEADER)

	// [index, address, key, keystore, password, notes]
	for address, wallet := range AccountList {
		row := make([]string, 0, len(ACCOUNT_LIST_HEADER))
		row = append(row, strconv.Itoa(index))
		row = append(row, address)
		row = append(row, wallet.PrivateKey())
		row = append(row, wallet.FilePath())
		row = append(row, wallet.Password())
		row = append(row, "")

		data = append(data, row)
		index++
	}

	return file.WriteAll(ACCOUNTS_SHEET_NAME, data)
}

// read keystore directory for all wallet
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
