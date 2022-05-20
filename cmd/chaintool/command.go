package main

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
	"utopia/internal/chain"
	"utopia/internal/excel"
	"utopia/internal/helper"
	"utopia/internal/logger"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
)

var (
	AddressFlag = cli.StringFlag{
		Name:  "address",
		Usage: "Address in hex mode",
		Value: "",
	}
	ChainFlag = cli.StringFlag{
		Name:  "chain",
		Usage: "Chain name (ie: btc, eth, bsc)",
		Value: "",
	}
	KeyFlag = cli.StringFlag{
		Name:  "key",
		Usage: "The file path of key store ",
		Value: "",
	}
	KeyDirFlag = cli.StringFlag{
		Name:  "keydir",
		Usage: "The directory of key store",
		Value: "",
	}
	PasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "The password of key store",
		Value: "",
	}
	ToFlag = cli.StringFlag{
		Name:  "to",
		Usage: "The dest address list",
		Value: "",
	}
	ValueFlag = cli.StringFlag{
		Name:  "value",
		Usage: "The value list in eth",
		Value: "",
	}
	GasPriceFlag = cli.Uint64Flag{
		Name:  "gas",
		Usage: "The gas price for tx",
		Value: 0,
	}
	HashFlag = cli.StringFlag{
		Name:  "hash",
		Usage: "The transaction hash in hex mode",
		Value: "",
	}
	FileFlag = cli.StringFlag{
		Name:  "file",
		Usage: "The excel file path with address list",
		Value: "",
	}

	cmdBalance = cli.Command{
		Name:   "balance",
		Usage:  "Query balance of accounts on chain",
		Action: QueryBalance,
		Flags: []cli.Flag{
			ChainFlag,
			AddressFlag,
			KeyDirFlag,
			PasswordFlag,
		},
	}
	cmdTransfer = cli.Command{
		Name:   "transfer",
		Usage:  "Transfer balance of accounts on chain",
		Action: TransferBalance,
		Flags: []cli.Flag{
			ChainFlag,
			KeyFlag,
			PasswordFlag,
			ToFlag,
			ValueFlag,
			FileFlag,
		},
	}
	cmdMerge = cli.Command{
		Name:   "merge",
		Usage:  "Merge balance of accounts on chain",
		Action: MergeBalance,
		Flags: []cli.Flag{
			ChainFlag,
			KeyDirFlag,
			PasswordFlag,
			ToFlag,
		},
	}
	cmdSpeedup = cli.Command{
		Name:   "speedup",
		Usage:  "Speedup transaction on chain",
		Action: Speedup,
		Flags: []cli.Flag{
			ChainFlag,
			KeyFlag,
			PasswordFlag,
			HashFlag,
			GasPriceFlag,
		},
	}
	cmdRpcServer = cli.Command{
		Name:   "rpc",
		Usage:  "Query rpc server list and test",
		Action: ListRpc,
		Flags: []cli.Flag{
			ChainFlag,
		},
	}
	cmdGas = cli.Command{
		Name:   "gas",
		Usage:  "Query current gas price on chain",
		Action: QueryGas,
		Flags: []cli.Flag{
			ChainFlag,
		},
	}
)

func QueryBalance(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	address := ctx.String(AddressFlag.Name)
	keydir := ctx.String(KeyDirFlag.Name)
	password := ctx.String(PasswordFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	addressList := make([]string, 0)
	if address != "" {
		addressList = append(addressList, strings.Split(address, ",")...)
	}

	if keydir != "" && password != "" {
		wallet, err := wallet.ListWallet(wallet.WALLET_ETH, keydir, password)
		if err != nil {
			return err
		}

		for _, w := range wallet {
			addressList = append(addressList, w.Address())
		}
	}

	logger.Debug("Total address is %d", len(addressList))

	chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if chain == nil {
		return errors.New("Connect chain failed")
	}
	defer chain.DisConnect()

	for _, addr := range addressList {
		balance, err := chain.Balance(addr)
		if err != nil {
			logger.Warn("Address[%s].balance=%d, error:%v", addr, 0, err)
			continue
		}

		fmt.Printf("Address[%s].balance=%f\n", addr, helper.WeiToEth(&balance))
	}

	return nil
}

func TransferBalance(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	key := ctx.String(KeyFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	file := ctx.String(FileFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	addresss := strings.Split(to, ",")
	values := strings.Split(value, ",")
	if len(addresss) != len(values) {
		return errors.New("Not match the address and value list")
	}

	if file != "" {
		alist, vlist, err := readTransferFile(file)
		if err != nil {
			return err
		}

		addresss = append(addresss, alist...)
		values = append(values, vlist...)
	}

	logger.Debug("Total address is %d", len(addresss))

	wallet := wallet.NewWallet(wallet.WALLET_ETH, key, password)
	err = wallet.LoadKey()
	if err != nil {
		return err
	}

	chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if chain == nil {
		return errors.New("Connect chain failed")
	}
	defer chain.DisConnect()

	total := new(big.Int)
	valueList := make([]*big.Int, 0)
	for _, v := range values {
		fv, _ := strconv.ParseFloat(v, 64)
		valueList = append(valueList, helper.EthToWei(float32(fv)))
		total = new(big.Int).Add(total, helper.EthToWei(float32(fv)))
	}

	err = chain.Transfer(addresss, valueList, wallet)
	if err != nil {
		return err
	}

	fmt.Printf("Total transfer value %f\n", helper.WeiToEth(total))
	return nil
}

func MergeBalance(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	keydir := ctx.String(KeyDirFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	to := ctx.String(ToFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	wallet, err := wallet.ListWallet(wallet.WALLET_ETH, keydir, password)
	if err != nil {
		return err
	}

	logger.Debug("Total merge number is %d", len(wallet))

	chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if chain == nil {
		return errors.New("Connect chain failed")
	}
	defer chain.DisConnect()

	gasprice, err := chain.GasPrice()
	if err != nil {
		return nil
	}

	toList := make([]string, 0)
	valueList := make([]string, 0)

	total := new(big.Int)
	for _, w := range wallet {
		if w.Address() == to {
			continue
		}

		balance, err := chain.Balance(w.Address())
		if err != nil {
			log.Warn("Merge balance on address %s", w.Address())
			continue
		}

		fee := new(big.Int).Mul(&gasprice, big.NewInt(31500)) // 31500 = 21000 * 1.5
		if fee.Cmp(&balance) >= 0 {
			continue
		}

		value := new(big.Int).Sub(&balance, fee)
		err = chain.Transfer([]string{to}, []*big.Int{value}, w)
		if err != nil {
			log.Warn("Merge balance on address %s", w.Address())
			continue
		}

		total = new(big.Int).Add(total, value)
		toList = append(toList, w.Address())
		valueList = append(valueList, strconv.FormatFloat(float64(helper.WeiToEth(value)), 'f', 5, 32))
	}

	err = saveTransferFile(toList, valueList, "./transfer.xlsx")
	if err != nil {
		log.Warn("Write transfer log failed with %v", err)
	}

	fmt.Printf("Total merge balance %f\n", helper.WeiToEth(total))
	return nil
}

func Speedup(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	key := ctx.String(KeyFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	gas := ctx.Uint64(GasPriceFlag.Name)
	hash := ctx.String(HashFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	wallet := wallet.NewWallet(wallet.WALLET_ETH, key, password)
	err = wallet.LoadKey()
	if err != nil {
		return err
	}

	chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if chain == nil {
		return errors.New("Connect chain failed")
	}
	defer chain.DisConnect()

	gasprice, err := chain.GasPrice()
	if err != nil {
		return err
	}

	if gas != 0 {
		if gas < gasprice.Uint64() {
			return errors.New("Input gas price is less")
		}

		gasprice = *new(big.Int).SetUint64(gas)
	}

	// query pending tx
	tx, pending, err := chain.Transaction(common.FromHex(hash))
	if err != nil {
		return err
	}

	if !pending {
		return errors.New("Transaction is not pending")
	}

	if gasprice.Cmp(tx.GasPrice()) < 0 {
		return errors.New("Input gas less than transaction gas")
	}

	err = chain.SendTransaction(types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), &gasprice, tx.Data()), wallet)
	if err != nil {
		return err
	}

	fmt.Printf("Speed up transaction %s", hash)
	return nil
}

func ListRpc(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)

	metaList := make([]chain.ChainMeta, 0)
	if chainName == "" {
		metaList = chain.ChainList
	} else {
		meta, err := chain.ChainMetaByName(chainName)
		if err != nil {
			return err
		}

		metaList = append(metaList, *meta)
	}

	if len(metaList) == 0 {
		return errors.New("Empty chain list")
	}

	for _, meta := range metaList {
		chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
		if chain == nil {
			return errors.New("Connect chain failed")
		}
		defer chain.DisConnect()

		for index, url := range meta.RpcServer {
			start := time.Now()

			err := chain.Connect([]string{url}, true)
			if err != nil {
				continue
			}

			latency := time.Now().Sub(start)
			fmt.Printf("[%s][%d] url=%s, latency=%d ms\n", meta.Name, index, url, latency.Milliseconds())
		}
	}

	return nil
}

func QueryGas(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)

	metaList := make([]chain.ChainMeta, 0)
	if chainName == "" {
		metaList = chain.ChainList
	} else {
		meta, err := chain.ChainMetaByName(chainName)
		if err != nil {
			return err
		}

		metaList = append(metaList, *meta)
	}

	if len(metaList) == 0 {
		return errors.New("Empty chain list")
	}

	for _, meta := range metaList {
		chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
		if chain == nil {
			return errors.New("Connect chain failed")
		}
		defer chain.DisConnect()

		gasprice, err := chain.GasPrice()
		if err != nil {
			return err
		}

		fmt.Printf("Gas %d on chain %s\n", gasprice.Uint64(), meta.Name)
	}

	return nil
}

func readTransferFile(path string) ([]string, []string, error) {
	to := make([]string, 0)
	value := make([]string, 0)

	file, err := excel.NewExcel(path)
	if err != nil {
		return to, value, err
	}

	err = file.Open()
	if err != nil {
		return to, value, err
	}
	defer file.Close(false)

	data, err := file.ReadAll("transfer")
	if err != nil {
		return to, value, err
	}

	for index, row := range data {
		// skip the header
		if index == 0 {
			continue
		}

		if len(row) != 3 {
			return to, value, errors.New("Invalid file format")
		}

		to = append(to, row[1])
		value = append(value, row[2])
	}

	return to, value, nil
}

func saveTransferFile(to []string, value []string, path string) error {
	if len(to) != len(value) {
		return errors.New("Not match to and value")
	}

	file, err := excel.NewExcel(path)
	if err != nil {
		return err
	}

	err = file.Open()
	if err != nil {
		return err
	}
	defer file.Close(true)

	data := make([][]string, 0)
	header := []string{"index", "address", "value"}
	data = append(data, header)

	for i, key := range to {
		row := make([]string, 0, 3)
		row = append(row, strconv.Itoa(i+1))
		row = append(row, key)
		row = append(row, value[i])

		data = append(data, row)
	}

	return file.WriteAll("transfer", data)
}
