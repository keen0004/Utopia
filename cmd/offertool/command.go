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

	cmdIXO = cli.Command{
		Name:   "ixo",
		Usage:  "IXO query",
		Action: IXOList,
		Flags: []cli.Flag{
			xxxFlag,
		},
	}
	cmdAirdrop = cli.Command{
		Name:   "airdrop",
		Usage:  "Airdrop query",
		Action: AirdropList,
		Flags: []cli.Flag{
			xxxFlag,
		},
	}
)

func IXOList(ctx *cli.Context) error {
	return nil
}

func AirdropList(ctx *cli.Context) error {
	return nil
}
