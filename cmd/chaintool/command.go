package main

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
	"utopia/internal/chain"
	"utopia/internal/helper"
	"utopia/internal/logger"
	"utopia/internal/wallet"

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

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	addresss := strings.Split(to, ",")
	values := strings.Split(value, ",")
	if len(addresss) != len(values) {
		return errors.New("Not match the address and value list")
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
	}

	fmt.Printf("Total merge balance %f\n", helper.WeiToEth(total))
	return nil
}

func Speedup(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	key := ctx.String(KeyFlag.Name)
	password := ctx.String(PasswordFlag.Name)

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

	// query pending tx
	txs, err := chain.PendingTransaction()
	if err != nil {
		return err
	}

	log.Debug("Total speed up transaction size %d", len(txs))

	count := 0
	for _, tx := range txs {
		if gasprice.Cmp(tx.GasPrice()) < 0 {
			continue
		}

		err = chain.SendTransaction(types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), &gasprice, tx.Data()), wallet)
		if err != nil {
			log.Warn("Speed up tx %s failed with error %v", tx.Hash().Hex(), err)
			continue
		}

		count++
	}

	fmt.Printf("Speed up transaction size %d", count)
	return nil
}

func ListRpc(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if chain == nil {
		return errors.New("Connect chain failed")
	}
	defer chain.DisConnect()

	for index, url := range meta.RpcServer {
		start := time.Now()

		err = chain.Connect([]string{url}, true)
		if err != nil {
			continue
		}

		latency := time.Now().Sub(start)
		fmt.Printf("[%d] url=%s, latency=%d ms\n", index, url, latency.Milliseconds())
	}

	return nil
}

func QueryGas(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
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

	fmt.Printf("Gas %d on chain %s\n", gasprice.Uint64(), chainName)
	return nil
}
