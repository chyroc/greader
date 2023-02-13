package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/chyroc/greader/app"
)

func main() {
	appCli := &cli.App{
		Name:        "greader",
		Description: "RSS service, providing api similar to google reader.",
		Commands: []*cli.Command{
			app.StartServer(),
			app.Register(),
		},
	}
	if err := appCli.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
