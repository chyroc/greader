package greader_api

import (
	"context"
	"fmt"
	"log"
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
	log.Println("[auth]", "username:", username, "password:", password)

	authInfo, err := r.s.Auth(ctx, username, password)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SID=%s\n"+
		"LSID=null\n"+
		"Auth=%s/%s\n", username, username, authInfo), nil
}
