package server

import (
	"github.com/gin-gonic/gin"

	"github.com/chyroc/greader/greader_api"
)

type ginHttpReader struct {
	c *gin.Context
}

func NewGinHttpReader(c *gin.Context) greader_api.HttpReader {
	return &ginHttpReader{c}
}

func (r *ginHttpReader) FormList(key string) []string {
	return r.c.PostFormArray(key)
}

func (r *ginHttpReader) FormString(key string) string {
	return r.c.PostForm(key)
}

func (r *ginHttpReader) HeaderString(key string) string {
	return r.c.GetHeader(key)
}

func (r *ginHttpReader) QueryString(key string) string {
	return r.c.Query(key)
}
