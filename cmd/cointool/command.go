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

const (
	HTTP_PROXY_URL = "http://127.0.0.1:1087"
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

	// setenv CMC_PRO_API_KEY first
	client := cmc.NewClient("", HTTP_PROXY_URL)
	if client == nil {
		return errors.New("Create cmc client failed")
	}

	infolist := make([]*cmc.Listing, 0)
	if coin != "" {
		result, err := client.LatestQuotes(&cmc.QuoteOptions{Symbol: coin, Convert: "USD"})
		if err != nil {
			return err
		}

		infolist = append(infolist, result...)
	} else {
		result, err := client.LatestListings(&cmc.ListingOptions{Start: 1, Limit: 2000, Convert: "USD"})
		if err != nil {
			return err
		}

		infolist = append(infolist, result...)
	}

	// "id", "symbol", "rank", "current", "total", "pairs", "platform", "address", "price", "marketcap", "lastupdated"
	if file != "" {
		err := helper.WriteCurrencyFile(file, infolist)
		if err != nil {
			return err
		}
	} else {
		for _, info := range infolist {
			fmt.Fprintf(os.Stderr, "id=%d, symbol=%s, rank=%d, current=%d, total=%d, pairs=%d, platform=%s, address=%s, price=%f, marketcap=%f, lastupdated=%s\n",
				int(info.ID), info.Symbol, int(info.CMCRank), int(info.CirculatingSupply), int(info.TotalSupply),
				int(info.NumMarketPairs), info.Platform.Symbol, info.Platform.TokenAddress, info.Quote["USD"].Price,
				info.Quote["USD"].MarketCap, info.Quote["USD"].LastUpdated)
		}
	}

	return nil
}

func QueryCollection(ctx *cli.Context) error {
	contract := ctx.String(CotractFlag.Name)

	client := opensea.New()
	c, err := client.Contract(context.Background(), &opensea.ContractRequest{AssetContractAddress: contract})
	if err != nil {
		return err
	}

	fmt.Printf("Contract: %v\n", c)

	return nil
}
