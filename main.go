package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chyroc/greader/adapter_mysql"
	"github.com/chyroc/greader/greader_api"
)

func main() {
	server := "127.0.0.1:8082"
	fmt.Println("http://" + server)

	logger := greader_api.NewDefaultLogger()

	dsn := "root:@tcp(127.0.0.1:3306)/greader?charset=utf8&parseTime=True&loc=Local"
	db, err := adapter_mysql.New(dsn, logger)
	if err != nil {
		log.Fatal(err)
		return
	}

	cli := greader_api.New(&greader_api.ClientConfig{Store: db, Logger: logger})

	log.Println(http.ListenAndServe(server, cli))
}
