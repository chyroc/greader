package greader_api

import (
	"context"
	"net/http"
	"strings"
)

func (r *Client) Token(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	auth := r.headerAuth(req)

	writer.Header().Set("Content-Type", "text/html; charset=UTF-8")
	writer.Write([]byte(pad(auth, 57, 'Z')))
}

func pad(s string, size int, pad rune) string {
	res := []rune(s)
	for len(res) < size {
		res = append(res, pad)
	}
	return string(res)
}

func (r *Client) headerAuth(req HttpReader) string {
	// Authorization:[GoogleLogin auth=admin/1664cfe9598edbc63fc0adf0d1464d98f42f4840]

	authorization := req.HeaderString("Authorization")
	if authorization == "" {
		return ""
	}
	res := strings.SplitN(authorization, "auth=", 2)
	if len(res) == 2 {
		return res[1]
	}
	return ""
}
