package main

import (
	"fmt"
	"os"
	"utopia/internal/helper"

	"gopkg.in/urfave/cli.v1"
)

var (
	version string = "1.0.0"
	usage   string = "Tool box for smart contract operations"
	app     *cli.App
)

func init() {
	app = helper.NewApp(version, usage)
	app.Commands = []cli.Command{
		cmdDeploy,
		cmdCall,
		cmdList,
		cmdERC20,
		cmdERC721,
	}
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
