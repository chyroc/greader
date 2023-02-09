package server

import (
	"github.com/gin-gonic/gin"

	"github.com/chyroc/greader/greader_api"
)

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func apiRegister(backend greader_api.IGReaderBackend) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := new(registerReq)
		err := c.BindJSON(req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err = backend.Register(c, req.Username, req.Password)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"error": ""})
	}
}
