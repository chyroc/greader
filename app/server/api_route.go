package server

import (
	"github.com/gin-gonic/gin"

	"github.com/chyroc/greader/greader_api"
)

func registerAPiRoute(ginIns *gin.Engine, greaderIns *greader_api.GReader) {
	// gin router
	api := ginIns.Group("/api")
	{
		greaderAPI := api.Group("/greader")
		{
			for path, handler := range greaderIns.Routers() {
				handler := handler
				greaderAPI.Handle(path[0], path[1], func(c *gin.Context) {
					handler(c, NewGinHttpReader(c), c.Writer)
				})
			}
		}
	}
}
