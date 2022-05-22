package main

import (
	"fmt"
	"os"
	"utopia/internal/config"
	"utopia/internal/helper"

	"gopkg.in/urfave/cli.v1"
)

var (
	version string = "1.0.0"
	usage   string = "Tool box for block chain operations"
	app     *cli.App
)

func init() {
	app = helper.NewApp(version, usage)
	app.Commands = []cli.Command{
		cmdBalance,
		cmdTransfer,
		cmdSpeedup,
		cmdRpcServer,
		cmdGas,
	}
}

func main() {
	err := config.Config.LoadConfig("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
