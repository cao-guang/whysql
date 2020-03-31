package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)
func main()  {
	app := cli.NewApp()
	app.Name = "caog"
	app.Usage = "caog工具"
	app.Version = Version
	app.Commands = []*cli.Command{
		{
			Name:            "tool",
			Aliases:         []string{"t"},
			Usage:           "caog tool",
			Action:          toolAction,
			SkipFlagParsing: true,
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "caog version",
			Action: func(c *cli.Context) error {
				fmt.Println(getVersion())
				return nil
			},
		},
		{
			Name:   "self-upgrade",
			Usage:  "caog self-upgrade",
			Action: upgradeAction,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}