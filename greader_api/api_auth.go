package greader_api

import (
	"context"
	"fmt"
	"net/http"
)

func (r *Client) Auth(ctx context.Context, reader HttpReader, writer http.ResponseWriter) {
	authInfo, err := r.auth(ctx, reader)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, authInfo)
	}
}

func (r *Client) auth(ctx context.Context, reader HttpReader) (string, error) {
	username := reader.FormString("Email")
	password := reader.FormString("Passwd")
	r.log.Info(ctx, "[auth] username=%s", username)

	authInfo, err := r.s.Login(ctx, username, password)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SID=%s\n"+
		"LSID=null\n"+
		"Auth=%s/%s\n", username, username, authInfo), nil
}
