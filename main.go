package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/chyroc/greader/greader_api"
	"github.com/chyroc/greader/mysql_backend"
	"github.com/chyroc/greader/server"
)

func main() {
	logger := greader_api.NewDefaultLogger()

	r := gin.New()
	r.Use(server.Log(logger))

	dsn := "root:@tcp(127.0.0.1:3306)/greader?charset=utf8&parseTime=True&loc=Local"
	db, err := mysql_backend.New(dsn, logger)
	if err != nil {
		log.Fatal(err)
		return
	}

	cli := greader_api.New(&greader_api.ClientConfig{
		Backend:     db,
		Logger:      logger,
		FetchLogger: logger,
	})
	cli.FetchRssBackend()

	api := r.Group("/api")
	{
		greaderAPI := api.Group("/greader")
		{
			for path, handler := range cli.Routers() {
				handler := handler
				greaderAPI.Handle(path[0], path[1], func(c *gin.Context) {
					handler(c, server.NewGinHttpReader(c), c.Writer)
				})
			}
		}
	}

	if err := r.Run("127.0.0.1:8082"); err != nil {
		log.Fatalln(err)
	}
}
