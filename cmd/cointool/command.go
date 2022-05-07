package main

import (
	"gopkg.in/urfave/cli.v1"
)

var (
	xxxFlag = cli.StringFlag{
		Name:  "xxx",
		Usage: "xxx",
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
					xxxFlag,
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
					xxxFlag,
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
					xxxFlag,
				},
			},
		},
	}
)

func QueryPrice(ctx *cli.Context) error {
	return nil
}

func QueryCoin(ctx *cli.Context) error {
	return nil
}

func QueryCollection(ctx *cli.Context) error {
	return nil
}
