package app

import (
	"fmt"

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

			if conf.AdminUsername != "" && conf.AdminPassword != "" {
				if err := app.Backend.Register(c.Context, conf.AdminUsername, conf.AdminPassword); err != nil {
					return err
				} else {
					fmt.Println("register admin success")
				}
			}

			return app.Gin.Run(":8081")
		},
	}
}
