package main

import (
	"fmt"
	"os"

	"github.com/MOSSV2/dimo-sdk-go/app/cmd"
	"github.com/MOSSV2/dimo-sdk-go/build"
	"github.com/urfave/cli/v2"
)

func main() {
	local := make([]*cli.Command, 0, 1)
	local = append(local, cmd.InitCmd)
	local = append(local, serverCmd)

	app := cli.App{
		Name:    "store-edge",
		Version: build.UserVersion(),
		Usage:   "Unibase DA Storage Node",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  cmd.RepoStr,
				Value: "~/.store",
				Usage: "storage node home dir",
			},
			&cli.StringFlag{
				Name:    cmd.PasswordStr,
				Aliases: []string{"pwd"},
				Value:   "",
			},
		},
		Commands: local,
	}
	app.Setup()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err) // nolint:errcheck
		os.Exit(1)
	}
}