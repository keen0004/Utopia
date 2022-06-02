package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"utopia/internal/cmc"
	"utopia/internal/helper"

	"github.com/pinealctx/opensea-go"
	"gopkg.in/urfave/cli.v1"
)

var (
	CoinFlag = cli.StringFlag{
		Name:  "coin",
		Usage: "Set the coin symbol, empty for all",
		Value: "",
	}
	FileFlag = cli.StringFlag{
		Name:  "file",
		Usage: "Set the file path for input or output",
		Value: "",
	}
	CotractFlag = cli.StringFlag{
		Name:  "contract",
		Usage: "Set the contract address in hex mode",
		Value: "",
	}
	AmountFlag = cli.Float64Flag{
		Name:  "amount",
		Usage: "Amount of coin",
		Value: 0.0,
	}
	MoneyFlag = cli.StringFlag{
		Name:  "money",
		Usage: "cash symbol like USD or CNY",
		Value: "USD",
	}

	cmdDex = cli.Command{
		Name:  "dex",
		Usage: "Dex operations on uniswap",
		Subcommands: []cli.Command{
			{
				Name:   "price",
				Usage:  "Query price of token",
				Action: QueryPrice,
				Flags: []cli.Flag{
					CoinFlag,
				},
			},
		},
	}
	cmdCMC = cli.Command{
		Name:  "cmc",
		Usage: "Coinmarketcap query operations",
		Subcommands: []cli.Command{
			{
				Name:   "query",
				Usage:  "Query coin information on cmc",
				Action: QueryCoin,
				Flags: []cli.Flag{
					CoinFlag,
					FileFlag,
					MoneyFlag,
				},
			},
			{
				Name:   "convert",
				Usage:  "Convert one coin to USDT with same value",
				Action: ConverCoin,
				Flags: []cli.Flag{
					CoinFlag,
					AmountFlag,
					MoneyFlag,
				},
			},
		},
	}
	cmdOpensea = cli.Command{
		Name:  "opensea",
		Usage: "Opensea query operations",
		Subcommands: []cli.Command{
			{
				Name:   "query",
				Usage:  "Query collection information on cmc",
				Action: QueryCollection,
				Flags: []cli.Flag{
					CotractFlag,
				},
			},
		},
	}
)

func QueryPrice(ctx *cli.Context) error {
	return nil
}

func QueryCoin(ctx *cli.Context) error {
	coin := ctx.String(CoinFlag.Name)
	file := ctx.String(FileFlag.Name)
	money := ctx.String(MoneyFlag.Name)

	// setenv CMC_PRO_API_KEY & HTTPS_PROXY first
	client := cmc.NewClient("", "")
	if client == nil {
		return errors.New("Create cmc client failed")
	}

	infolist := make([]*cmc.Listing, 0)
	if coin != "" {
		result, err := client.LatestQuotes(&cmc.QuoteOptions{Symbol: coin, Convert: money})
		if err != nil {
			return err
		}

		infolist = append(infolist, result...)
	} else {
		result, err := client.LatestListings(&cmc.ListingOptions{Start: 1, Limit: 2000, Convert: money})
		if err != nil {
			return err
		}

		infolist = append(infolist, result...)
	}

	// "id", "symbol", "rank", "current", "total", "pairs", "platform", "address", "price", "marketcap", "lastupdated"
	if file != "" {
		err := helper.WriteCurrencyFile(file, money, infolist)
		if err != nil {
			return err
		}
	} else {
		for _, info := range infolist {
			fmt.Fprintf(os.Stderr, "id=%d, symbol=%s, rank=%d, current=%d, total=%d, pairs=%d, platform=%s, address=%s, price=%.5f, marketcap=%.2f, lastupdated=%s\n",
				int(info.ID), info.Symbol, int(info.CMCRank), int(info.CirculatingSupply), int(info.TotalSupply),
				int(info.NumMarketPairs), info.Platform.Symbol, info.Platform.TokenAddress, info.Quote[money].Price,
				info.Quote[money].MarketCap, info.Quote[money].LastUpdated)
		}
	}

	return nil
}

func ConverCoin(ctx *cli.Context) error {
	coin := ctx.String(CoinFlag.Name)
	amount := ctx.Float64(AmountFlag.Name)
	money := ctx.String(MoneyFlag.Name)

	// setenv CMC_PRO_API_KEY & HTTPS_PROXY first
	client := cmc.NewClient("", "")
	if client == nil {
		return errors.New("Create cmc client failed")
	}

	result, err := client.PriceConversion(&cmc.ConvertOptions{Amount: amount, Symbol: coin, Convert: money})
	if err != nil {
		return err
	}

	for _, o := range result {
		fmt.Fprintf(os.Stderr, "%.5f %s = %.5f %s\n", amount, coin, o.Quote[money].Price, money)
	}

	return nil
}

func QueryCollection(ctx *cli.Context) error {
	contract := ctx.String(CotractFlag.Name)

	// shold set env HTTPS_PROXY
	client := opensea.New()
	c, err := client.Contract(context.Background(), &opensea.ContractRequest{AssetContractAddress: contract})
	if err != nil {
		return err
	}

	fmt.Printf("Contract: %v\n", c)

	return nil
}
