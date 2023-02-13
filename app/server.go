package app

import (
	"github.com/urfave/cli/v2"

	"github.com/chyroc/greader/app/server"
)

func StartServer() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "start greader server",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			conf, err := loadConf()
			if err != nil {
				return err
			}

			app, err := server.New(conf.DSN())
			if err != nil {
				return err
			}

			return app.Start(":8081")
		},
	}
}
