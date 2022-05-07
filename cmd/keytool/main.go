package main

import (
	"os"
	"utopia/internal/helper"
	"utopia/internal/logger"

	"gopkg.in/urfave/cli.v1"
)

var (
	version string = "1.0.0"
	usage   string = "Tool box for account key operations"
	app     *cli.App
)

func init() {
	app = helper.NewApp(version, usage)
	app.Commands = []cli.Command{
		cmdGenerate,
		cmdList,
		cmdSign,
		cmdVerify,
		cmdHash,
	}
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		logger.Error("Application run with error: %v", err)
		os.Exit(1)
	}
}
