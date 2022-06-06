package main

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"
	"utopia/internal/chain"
	"utopia/internal/config"
	"utopia/internal/helper"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gopkg.in/urfave/cli.v1"
)

var (
	ChainFlag = cli.StringFlag{
		Name:  "chain",
		Usage: "Chain name",
		Value: "",
	}
	AddressFlag = cli.StringFlag{
		Name:  "address",
		Usage: "Address in hex mode",
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
			AddressFlag,
		},
	}
	cmdTransfer = cli.Command{
		Name:   "transfer",
		Usage:  "Transfer balance of accounts on chain",
		Action: TransferBalance,
		Flags: []cli.Flag{
			ToFlag,
			ValueFlag,
			FileFlag,
		},
	}
	cmdSpeedup = cli.Command{
		Name:   "speedup",
		Usage:  "Speedup transaction on chain",
		Action: Speedup,
		Flags: []cli.Flag{
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
	address := ctx.String(AddressFlag.Name)

	// get chain meta and connect it
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	balance, err := c.Balance(address)
	if err != nil {
		return err
	}

	fmt.Printf("Address[%s].balance=%f\n", address, helper.WeiToEth(balance))
	return nil
}

func TransferBalance(ctx *cli.Context) error {
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	file := ctx.String(FileFlag.Name)

	translist := make([]helper.TransferInfo, 0)
	if to != "" && value != "" {
		translist = append(translist, helper.TransferInfo{
			From:  config.Config.Chain.From,
			To:    to,
			Value: value,
		})
	}

	// read file transfer list
	if file != "" {
		list, err := helper.ReadTransferFile(file)
		if err != nil {
			return err
		}

		translist = append(translist, list...)
	}

	// get chain meta and connect it
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	for _, info := range translist {
		// get wallet for sign transaction
		wallet, err := wallet.GetWallet(info.From)
		if err != nil {
			return err
		}

		fv, _ := strconv.ParseFloat(info.Value, 64)
		tx, err := c.Transfer(info.To, helper.EthToWei(float32(fv)), wallet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Tranfer %s to %s failed with err: %v\n", info.Value, info.To, err)
		} else {
			fmt.Fprintf(os.Stderr, "Tranfer %s to %s with transaction %s\n", info.Value, info.To, tx)
		}
	}

	return nil
}

func Speedup(ctx *cli.Context) error {
	gas := ctx.Uint64(GasPriceFlag.Name)
	hash := ctx.String(HashFlag.Name)

	// get chain meta and connect it
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	// query current gasprice on chain
	gasprice, err := c.GasPrice()
	if err != nil {
		return err
	}

	if gas != 0 {
		if gas < gasprice.Uint64() {
			return errors.New("Input gas price is less")
		}

		gasprice = new(big.Int).SetUint64(gas)
	}

	// query pending tx
	tx, pending, err := c.Transaction(common.FromHex(hash))
	if err != nil {
		return err
	}

	if !pending {
		return errors.New("Transaction is not pending")
	}

	if gasprice.Cmp(tx.GasPrice()) < 0 {
		return errors.New("Input gas less than transaction gas")
	}

	// get wallet for sign transaction
	wallet, err := wallet.GetWallet(config.Config.Chain.From)
	if err != nil {
		return err
	}

	// resend transaction with higher gas
	hash, err = c.SendTransaction(types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), gasprice, tx.Data()), wallet)
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

	// ping all rpc server
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

	// query gas from all chains
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
