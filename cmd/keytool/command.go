package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"sync"
	"utopia/internal/config"
	"utopia/internal/excel"
	"utopia/internal/helper"
	"utopia/internal/logger"
	"utopia/internal/wallet"

	"github.com/cheggaaa/pb/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/urfave/cli.v1"
)

const (
	MAX_KEY_SIZE    = 10000
	MAX_THREAD_SIZE = 8
)

type KeyInfo struct {
	privatekey string
	address    string
}

var (
	KeyDirFlag = cli.StringFlag{
		Name:  "keydir",
		Usage: "Specfiles the directory for store key files.",
		Value: "",
	}
	KeyNumFlag = cli.UintFlag{
		Name:  "keynum",
		Usage: "Specfiles the number of key files.(max size 10000)",
		Value: 1,
	}
	PasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "Specfiles the password of key file.",
		Value: "",
	}
	ThreadNumFlag = cli.UintFlag{
		Name:  "thread",
		Usage: "Specfiles the thread number to generate keys.(max size 8)",
		Value: 1,
	}
	KeyFlag = cli.StringFlag{
		Name:  "key",
		Usage: "Specfiles the sign key data",
		Value: "",
	}
	KeyStoreFlag = cli.StringFlag{
		Name:  "keystore",
		Usage: "Specfiles the key store file",
		Value: "",
	}
	DataFlag = cli.StringFlag{
		Name:  "data",
		Usage: "Specfiles the sign data in hex mode",
		Value: "",
	}
	SignFlag = cli.StringFlag{
		Name:  "sign",
		Usage: "Specfiles the signature in hex mode or rsv sperate by |",
		Value: "",
	}

	cmdGenerate = cli.Command{
		Name:   "gen",
		Usage:  "Batch generate key store files",
		Action: GenKeyFiles,
		Flags: []cli.Flag{
			KeyDirFlag,
			KeyNumFlag,
			PasswordFlag,
			ThreadNumFlag,
		},
	}
	cmdList = cli.Command{
		Name:   "list",
		Usage:  "list key information in directory",
		Action: ListKey,
		Flags: []cli.Flag{
			KeyDirFlag,
			PasswordFlag,
		},
	}
	cmdSign = cli.Command{
		Name:   "sign",
		Usage:  "sign message by private key",
		Action: SignMessage,
		Flags: []cli.Flag{
			DataFlag,
		},
	}
	cmdVerify = cli.Command{
		Name:   "verify",
		Usage:  "verify signature",
		Action: VerifySig,
		Flags: []cli.Flag{
			SignFlag,
			DataFlag,
		},
	}
	cmdHash = cli.Command{
		Name:   "hash",
		Usage:  "hash data",
		Action: HashData,
		Flags: []cli.Flag{
			DataFlag,
		},
	}
)

func GenKeyFiles(ctx *cli.Context) error {
	keydir := ctx.String(KeyDirFlag.Name)
	keynum := ctx.Int(KeyNumFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	thread := ctx.Int(ThreadNumFlag.Name)

	if keydir == "" || keynum <= 0 || password == "" || thread <= 0 {
		return errors.New("Invalid parameters for generate keys")
	}

	// check the max input
	if keynum > MAX_KEY_SIZE {
		keynum = MAX_KEY_SIZE
	}

	if thread > MAX_THREAD_SIZE {
		thread = MAX_THREAD_SIZE
	}

	// create key store files directory
	info, err := os.Stat(keydir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(keydir, 0700)
		if err != nil {
			return err
		}
	} else if !info.IsDir() {
		return errors.New("Data dir is exist")
	}

	outexcel := path.Join(keydir, "keys.xlsx")
	_, err = os.Stat(outexcel)
	if !os.IsNotExist(err) {
		return errors.New("Out excel file[" + outexcel + "] is exist")
	}

	logger.Debug("Generate %d key files to %s by %d thread", keynum, keydir, thread)

	var lock sync.Mutex
	walletlist := make([]wallet.Wallet, 0, keynum)
	bar := pb.StartNew(keynum)
	wg := sync.WaitGroup{}
	wg.Add(int(thread))

	// start go routine for generate key files
	for i := 0; i < thread; i++ {
		// caclute the generate number by this thread
		num := keynum / thread
		start := num*i + 1
		if i == thread-1 {
			num = keynum - num*i
		}

		go func(size int, index int, bar *pb.ProgressBar) {
			defer wg.Done()

			wallet, err := batchGenKey(index, keydir, size, password, bar)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			lock.Lock()
			walletlist = append(walletlist, wallet...)
			lock.Unlock()
		}(num, start, bar)
	}

	wg.Wait()
	bar.Finish()

	logger.Debug("Write excel %s with key size %d", outexcel, len(walletlist))

	err = writeExcel(outexcel, walletlist)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Success generate all key files\n")
	return nil
}

func ListKey(ctx *cli.Context) error {
	keydir := ctx.String(KeyDirFlag.Name)
	password := ctx.String(PasswordFlag.Name)

	if keydir == "" || password == "" {
		return errors.New("Invalid parameters for list keys")
	}

	info, err := os.Stat(keydir)
	if os.IsNotExist(err) {
		return errors.New("key dir not exist")
	}

	// read one key file
	if !info.IsDir() {
		wallet := wallet.NewWallet(wallet.WALLET_ETH, keydir, password)
		if !wallet.IsKeyFile(info) {
			return errors.New("Not a valid key store")
		}

		err := wallet.LoadKey()
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "key %d: Address = %s, Private = %s\n", 1, wallet.Address(), wallet.PrivateKey())
	} else {
		wallet, err := wallet.ListWallet(wallet.WALLET_ETH, keydir, password)
		if err != nil {
			return err
		}

		for index, w := range wallet {
			fmt.Fprintf(os.Stderr, "key %d: Address = %s, Private = %s\n", index+1, w.Address(), w.PrivateKey())
		}
	}

	return nil
}

func SignMessage(ctx *cli.Context) error {
	data := ctx.String(DataFlag.Name)
	if data == "" {
		return errors.New("Invalid parameters for sign data")
	}

	// load configs
	w, err := wallet.GetWallet(config.Config.From)
	if err != nil {
		return err
	}

	// sign data
	key, _ := crypto.ToECDSA(common.FromHex(w.PrivateKey()))
	sign, err := crypto.Sign(crypto.Keccak256(common.FromHex(data)), key)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Sign result: 0x%s\n", hex.EncodeToString(sign))
	return nil
}

func VerifySig(ctx *cli.Context) error {
	sign := ctx.String(SignFlag.Name)
	data := ctx.String(DataFlag.Name)

	if sign == "" || data == "" {
		return errors.New("Invalid parameters for verify sign")
	}

	var hash []byte
	if len(data) == 66 {
		hash = common.FromHex(data)
	} else {
		hash = crypto.Keccak256(common.FromHex(data))
	}

	// recover public key from signature
	recoveredPubkey, err := crypto.SigToPub(hash, common.FromHex(sign))
	if err != nil {
		return err
	}

	// calculate the address from public key
	address := crypto.PubkeyToAddress(*recoveredPubkey)
	fmt.Fprintf(os.Stderr, "Data signed by address: %s\n", address.Hex())

	return nil
}

func HashData(ctx *cli.Context) error {
	data := ctx.String(DataFlag.Name)
	if data == "" {
		return errors.New("Invalid parameters for hash data")
	}

	hash := crypto.Keccak256(common.FromHex(data))

	// another impl
	// var hash []byte
	// h := sha3.NewLegacyKeccak256()
	// h.Write(common.FromHex(data))
	// hash = h.Sum(hash)

	fmt.Fprintf(os.Stderr, "Hash result: 0x%s\n", common.Bytes2Hex(hash))
	return nil
}

func batchGenKey(start int, dir string, size int, password string, bar *pb.ProgressBar) ([]wallet.Wallet, error) {
	wallets := make([]wallet.Wallet, 0, size)

	for i := 0; i < size; i++ {
		wallet := wallet.NewWallet(wallet.WALLET_ETH, path.Join(dir, fmt.Sprintf("key_%d", start+i)), password)

		// generate new wallet
		err := wallet.GenerateKey()
		if err != nil {
			logger.Error("generate key_%d error %v", start+i, err)

			bar.Add(1)
			continue
		}

		// save key data to file
		err = wallet.SaveKey()
		if err != nil {
			logger.Error("save key_%d error %v", start+i, err)

			bar.Add(1)
			continue
		}

		wallets = append(wallets, wallet)
		bar.Add(1)
	}

	return wallets, nil
}

func writeExcel(path string, data []wallet.Wallet) error {
	excel, err := excel.NewExcel(path)
	if err != nil {
		return err
	}

	err = excel.Open()
	if err != nil {
		return err
	}
	defer excel.Close(true)

	// generate write data
	values := make([][]string, 0)
	values = append(values, wallet.ACCOUNT_LIST_HEADER)

	for i, key := range data {
		row := make([]string, 0, 3)
		row = append(row, strconv.Itoa(i+1))
		row = append(row, helper.DefaultVlue(key.Address(), "0x"))
		row = append(row, helper.DefaultVlue(key.PrivateKey(), "0x"))
		row = append(row, helper.DefaultVlue(key.FilePath(), "x"))
		row = append(row, helper.DefaultVlue(key.Password(), "x"))
		row = append(row, "x")

		values = append(values, row)
	}

	// write excel data
	err = excel.WriteAll(wallet.ACCOUNTS_SHEET_NAME, values)
	if err != nil {
		return err
	}

	return nil
}
