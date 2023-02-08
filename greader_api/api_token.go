package greader_api

import (
	"context"
	"net/http"
	"strings"
)

func (r *Client) Token(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	_, auth := GetContextAuth(ctx)

	writer.Header().Set("Content-Type", "text/html; charset=UTF-8")
	writer.Write([]byte(pad(auth, 57, 'Z')))
}

func (r *Client) setAuth(ctx context.Context, req HttpReader) context.Context {
	username, auth := r.headerAuth(req)
	return setContext(ctx, username, auth)
}

func pad(s string, size int, pad rune) string {
	res := []rune(s)
	for len(res) < size {
		res = append(res, pad)
	}
	return string(res)
}

func (r *Client) headerAuth(req HttpReader) (string, string) {
	// Authorization:[GoogleLogin auth=admin/1664cfe9598edbc63fc0adf0d1464d98f42f4840]

	authorization := req.HeaderString("Authorization")
	if authorization == "" {
		return "", ""
	}
	res := strings.SplitN(authorization, "auth=", 2)
	if len(res) == 2 {
		x := strings.SplitN(res[1], "/", 2)
		if len(x) == 2 {
			return x[0], x[1]
		}
		return "", ""
	}
	return "", ""
}

const (
	keyUsername = "username"
	keyAuth     = "auth"
	keyPath     = "path"
)

func setContext(ctx context.Context, username, auth string) context.Context {
	if username != "" {
		ctx = context.WithValue(ctx, keyUsername, username)
	}
	if auth != "" {
		ctx = context.WithValue(ctx, keyAuth, auth)
	}
	return ctx
}

func getContext(ctx context.Context) (string, string) {
	username := getContextString(ctx, keyUsername)
	auth := getContextString(ctx, keyAuth)
	return username, auth
}

func getContextUsername(ctx context.Context) string {
	return getContextString(ctx, keyUsername)
}

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
