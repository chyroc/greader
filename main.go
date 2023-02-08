package main

import (
	"log"

	"github.com/chyroc/greader/server_app"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/greader?charset=utf8&parseTime=True&loc=Local"

	app, err := server_app.New(dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Start("127.0.0.1:8082"); err != nil {
		log.Fatalln(err)
	}
}
