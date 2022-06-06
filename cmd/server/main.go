package main

import (
	"fmt"
	"os"
	"utopia/internal/config"
	"utopia/internal/helper"

	"github.com/gin-gonic/gin"
	"gopkg.in/urfave/cli.v1"
)

var (
	version string = "1.0.0"
	usage   string = "Utopia service"
	app     *cli.App
)

func init() {
	app = helper.NewApp(version, usage)
	app.Action = service
	app.Commands = []cli.Command{}
}

func service(*cli.Context) error {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/user/login")
	}

	return nil
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
