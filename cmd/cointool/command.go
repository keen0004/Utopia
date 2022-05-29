package main

import (
	"errors"
	"fmt"
	"os"
	"utopia/internal/cmc"

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
					CoinFlag,
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
	//file := ctx.String(FileFlag.Name)

	client := cmc.NewClient(&cmc.Config{ProAPIKey: os.Getenv("CMC_KEY")})
	if client == nil {
		return errors.New("Create cmc client failed")
	}

	result, err := client.Cryptocurrency.Info(&cmc.InfoOptions{Symbol: coin})
	if err != nil {
		return err
	}

	for symbol, info := range result {
		fmt.Fprintf(os.Stderr, "[%s] %v", symbol, info)
	}

	return nil
}

func QueryCollection(ctx *cli.Context) error {
	return nil
}
