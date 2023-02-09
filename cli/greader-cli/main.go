package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/chyroc/greader/cli/greader-cli/internal"
)

func main() {
	name := "greader-cli"
	app := &cli.App{
		Name: name,
		Commands: []*cli.Command{
			internal.CmdRegister(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
