package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	contracts "utopia/contracts/test"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/urfave/cli.v1"
)

var (
	MethodFlag = cli.StringFlag{
		Name:  "method",
		Usage: "Specfiles the method for call.",
		Value: "",
	}
	KeyFileFlag = cli.StringFlag{
		Name:  "key",
		Usage: "Specfiles the key file for sign tx.",
		Value: "",
	}
	PasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "Specfiles the password of key file.",
		Value: "",
	}
	RpcFlag = cli.StringFlag{
		Name:  "rpc",
		Usage: "Specfiles the url of rpc node.",
		Value: "",
	}
	AddressFlag = cli.StringFlag{
		Name:  "address",
		Usage: "Specfiles the address of contract.",
		Value: "",
	}
)

func DeployContract(ctx *cli.Context) error {
	keyfile := ctx.String(KeyFileFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	rpc := ctx.String(RpcFlag.Name)

	// check parameters
	if keyfile == "" || password == "" || rpc == "" {
		return errors.New("Invalid parameters for deploy contract")
	}

	// connect blockchain
	chain, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}

	key, _ := ioutil.ReadFile(keyfile)
	auth, err := bind.NewTransactor(strings.NewReader(string(key)), password)
	if err != nil {
		return err
	}

	address, tx, _, err := contracts.DeploySimple(auth, chain, "init message")
	if err != nil {
		return err
	}

	fmt.Printf("Deploy simple on address: %s, hash: %s\n", address.Hex(), tx.Hash().Hex())

	_, err = bind.WaitDeployed(context.Background(), chain, tx)
	if err != nil {
		return err
	}

	return nil
}

func CallContract(ctx *cli.Context) error {
	keyfile := ctx.String(KeyFileFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	rpc := ctx.String(RpcFlag.Name)
	address := ctx.String(AddressFlag.Name)

	// check parameters
	if keyfile == "" || password == "" || rpc == "" {
		return errors.New("Invalid parameters for call contract")
	}

	chain, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}

	key, _ := ioutil.ReadFile(keyfile)
	auth, err := bind.NewTransactor(strings.NewReader(string(key)), password)
	if err != nil {
		return err
	}

	simple, err := contracts.NewSimple(common.HexToAddress(address), chain)
	if err != nil {
		return err
	}

	message, err := simple.GetMessage(nil)
	if err != nil {
		return err
	}

	fmt.Printf("Get message %s\n", message)

	tx, err := simple.SetMessage(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, Value: nil}, "New message")
	if err != nil {
		return err
	}

	fmt.Printf("Send transaction %s\n", tx.Hash().Hex())

	receipt, err := bind.WaitMined(context.Background(), chain, tx)
	if err != nil {
		return err
	}

	fmt.Printf("used gas: %d\n", receipt.GasUsed)
	return nil
}
