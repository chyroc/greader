package greader_api

import (
	"context"
	"fmt"
	"net/http"
)

func (r *Client) Auth(ctx context.Context, reader HttpReader, writer http.ResponseWriter) {
	fmt.Println("auth")

	authInfo, err := r.auth(ctx, reader)
	if err != nil {
		r.renderErr(writer, err)
	} else {
		r.renderData(writer, authInfo)
	}
}

func (r *Client) auth(ctx context.Context, reader HttpReader) (string, error) {
	username := reader.FormString("Email")
	password := reader.FormString("Passwd")
	authInfo, err := r.s.Auth(ctx, username, password)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SID=%s\n"+
		"LSID=null\n"+
		"Auth=%s\n", username, authInfo), nil
}
