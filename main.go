package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chyroc/greader/adapter/sql_store"
	"github.com/chyroc/greader/greader_api"
)

func main() {
	server := "127.0.0.1:8082"
	fmt.Println("http://" + server)

	dsn := "root:@tcp(127.0.0.1:3306)/greader?charset=utf8&parseTime=True&loc=Local"
	db, err := sql_store.New(dsn)
	if err != nil {
		log.Fatal(err)
		return
	}

	cli := greader_api.New(&greader_api.ClientConfig{db, greader_api.NewDefaultLogger()})

	log.Println(http.ListenAndServe(server, cli))
}
