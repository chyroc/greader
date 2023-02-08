package greader_api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (r *GReader) Auth(ctx context.Context, reader HttpReader, writer http.ResponseWriter) {
	authInfo, err := r.auth(ctx, reader)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, authInfo)
	}
}

func (r *GReader) auth(ctx context.Context, reader HttpReader) (string, error) {
	username := reader.FormString("Email")
	password := reader.FormString("Passwd")
	r.log.Info(ctx, "[auth] username=%s", username)

	authInfo, err := r.backend.Login(ctx, username, password)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SID=%s\n"+
		"LSID=null\n"+
		"Auth=%s/%s\n", username, username, authInfo), nil
}

func (r *GReader) getHeaderAuth(req HttpReader) (string, string) {
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
