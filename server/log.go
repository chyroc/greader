package server

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/chyroc/greader/greader_api"
)

func Log(logger greader_api.ILogger) func(c *gin.Context) {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		query := c.Request.URL.RawQuery
		headerAuth := c.Request.Header.Get("Authorization")

		buf := new(bytes.Buffer)
		bs, _ := io.ReadAll(io.TeeReader(c.Request.Body, buf))
		c.Request.Body = io.NopCloser(buf)

		logger.Info(c, "[GIN][req_] %s\t%s query=%v Authorization=%v body=%v", method, path, query, headerAuth, string(bs))

		writerBuf := new(bytes.Buffer)
		c.Writer = &writerWrap{c.Writer, writerBuf}

		c.Next()

		logger.Info(c, "[GIN][resp] %s\t%s status=%v body=%s", method, path, c.Writer.Status(), writerBuf.String())
	}
}

type writerWrap struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (r *writerWrap) Write(bs []byte) (int, error) {
	r.buf.Write(bs)
	return r.ResponseWriter.Write(bs)
}
