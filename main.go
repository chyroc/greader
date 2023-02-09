package main

import (
	"fmt"
	"log"

	"github.com/chyroc/go-loader"

	"github.com/chyroc/greader/server"
)

type Config struct {
	MySQLHost     string `loader:"env,key=MYSQL_HOST"`
	MySQLUsername string `loader:"env,key=MYSQL_USERNAME"`
	MySQLPassword string `loader:"env,key=MYSQL_PASSWORD"`
	MySQLDatabase string `loader:"env,key=MYSQL_DATABASE"`

	DisableRegister bool `loader:"env,key=DISABLE_REGISTER"`
}

func main() {
	conf := new(Config)
	if err := loader.Load(conf); err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.MySQLUsername, conf.MySQLPassword, conf.MySQLHost, conf.MySQLDatabase)
	app, err := server.New(dsn, conf.DisableRegister)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Start(":8081"); err != nil {
		log.Fatalln(err)
	}
}
