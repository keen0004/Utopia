package main

import (
	"errors"
	"fmt"
	"utopia/contracts/token"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/urfave/cli.v1"
)

var (
	TypeFlag = cli.StringFlag{
		Name:  "type",
		Usage: "Token type (erc20 or erc721)",
		Value: "erc20",
	}
	RpcFlag = cli.StringFlag{
		Name:  "rpc",
		Usage: "Specfiles the url of rpc node.",
		Value: "",
	}

	commandToken cli.Command = cli.Command{
		Name:  "token",
		Usage: "Token operations",
		Subcommands: []cli.Command{
			{
				Name:   "balance",
				Usage:  "Query balance of account",
				Action: QueryBalance,
				Flags: []cli.Flag{
					RpcFlag,
					ContractFlag,
					AccountFlag,
					TypeFlag,
				},
			},
		},
	}
)

func QueryBalance(ctx *cli.Context) error {
	rpc := ctx.String(RpcFlag.Name)
	contract := ctx.String(ContractFlag.Name)
	account := ctx.String(AccountFlag.Name)
	ctype := ctx.String(TypeFlag.Name)

	// check parameters
	if contract == "" || account == "" || ctype == "" || rpc == "" {
		return errors.New("Invalid parameters for call contract")
	}

	chain, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}

	if ctype == "erc20" {
		erc20, err := token.NewERC20(common.HexToAddress(contract), chain)
		if err != nil {
			return err
		}

		balance, err := erc20.BalanceOf(nil, common.HexToAddress(account))
		if err != nil {
			return err
		}

		fmt.Printf("Get ERC20 balance: %d\n", balance)
	} else if ctype == "erc721" {
		erc721, err := token.NewERC721(common.HexToAddress(contract), chain)
		if err != nil {
			return err
		}

		balance, err := erc721.BalanceOf(nil, common.HexToAddress(account))
		if err != nil {
			return err
		}

		fmt.Printf("Get ERC20 balance: %d\n", balance)
	} else {
		return errors.New("Not support contract type")
	}

	return nil
}
