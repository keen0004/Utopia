package main

import (
	"fmt"
	"os"
	"utopia/internal/helper"

	"gopkg.in/urfave/cli.v1"
)

var (
	version string = "1.0.0"
	usage   string = "Simple contract call"
	app     *cli.App

	commandDeploy cli.Command = cli.Command{
		Name:    "deploy",
		Aliases: []string{"deploy", "d"},
		Usage:   "depoly simple contract",
		Action:  DeployContract,
		Flags: []cli.Flag{
			KeyFileFlag,
			PasswordFlag,
			RpcFlag,
		},
	}
	commandCall cli.Command = cli.Command{
		Name:    "call",
		Aliases: []string{"call", "c"},
		Usage:   "call simple contract",
		Action:  CallContract,
		Flags: []cli.Flag{
			MethodFlag,
			KeyFileFlag,
			PasswordFlag,
			RpcFlag,
			AddressFlag,
		},
	}
)

func init() {
	app = helper.NewApp(version, usage)
	app.Commands = []cli.Command{
		commandDeploy,
		commandCall,
	}
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
