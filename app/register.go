package app

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/chyroc/greader/greader_api"
	"github.com/chyroc/greader/mysql_backend"
)

func Register() *cli.Command {
	return &cli.Command{
		Name:        "register",
		Description: "register a new user",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Usage: "greader app host"},
			&cli.StringFlag{Name: "username", Usage: "username"},
			&cli.StringFlag{Name: "password", Usage: "password"},
		},
		Action: func(c *cli.Context) error {
			conf, err := loadConf()
			if err != nil {
				return err
			}
			username, password := c.String("username"), c.String("password")

			backend, err := mysql_backend.New(conf.DSN(), greader_api.NewDefaultLogger())
			if err != nil {
				return err
			}

			if err := backend.Register(c.Context, username, password); err != nil {
				return err
			}
			fmt.Println("register success")
			return nil
		},
	}
}
