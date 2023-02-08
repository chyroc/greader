package greader_api

import (
	"context"
	"net/http"
)

func (r *GReader) Token(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	_, authInfo := r.getHeaderAuth(req)

	writer.Header().Set("Content-Type", "text/html; charset=UTF-8")
	writer.Write([]byte(pad(authInfo, 57, 'Z')))
}

func pad(s string, size int, pad rune) string {
	res := []rune(s)
	for len(res) < size {
		res = append(res, pad)
	}
	return string(res)
}

const (
	keyUsername = "username"
	keyAuth     = "auth"
	keyPath     = "path"
)

func getContextString(ctx context.Context, key string) string {
	val := ctx.Value(key)
	if val == nil {
		return ""
	}
	s, _ := val.(string)
	return s
}

func getContextPath(ctx context.Context) string {
	return getContextString(ctx, keyPath)
}

func setContextPath(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, keyPath, s)
}
