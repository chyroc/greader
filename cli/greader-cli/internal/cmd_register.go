package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
)

func CmdRegister() *cli.Command {
	return &cli.Command{
		Name:        "register",
		Description: "register a new user",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Usage: "greader server host"},
			&cli.StringFlag{Name: "username", Usage: "username"},
			&cli.StringFlag{Name: "password", Usage: "password"},
		},
		Action: func(c *cli.Context) error {
			err := register(c.String("host"), c.String("username"), c.String("password"))
			if err != nil {
				return err
			}
			fmt.Println("register success")
			return nil
		},
	}
}

func register(host, username, password string) error {
	bs, _ := json.Marshal(map[string]any{
		"username": username,
		"password": password,
	})
	resp, err := http.DefaultClient.Post(strings.TrimRight(host, "/")+"/api/v2/auth/register", "application/json", bytes.NewReader(bs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bs, _ = io.ReadAll(resp.Body)
	var data map[string]string
	_ = json.Unmarshal(bs, &data)
	if data["error"] != "" {
		return fmt.Errorf("register failed: %s", data["error"])
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("register failed: %s", resp.Status)
	}
	return nil
}
