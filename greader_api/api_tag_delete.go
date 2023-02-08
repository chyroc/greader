package greader_api

import (
	"context"
	"net/http"
)

// api: https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPICaller.swift#L258
//
// api path: /reader/api/0/disable-tag
//
// arg: T=\(token)&s=\(folderExternalID)"
//
// &s=user/-/label/分类2asdf

func (r *Client) TagDelete(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	err := r.tagDelete(ctx, req)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, nil)
	}
}

func (r *Client) tagDelete(ctx context.Context, req HttpReader) error {
	username := getContextUsername(ctx)
	tagName := getUserLabelName(req.FormString("s"))
	r.log.Info(ctx, "[tagDelete] username=%s, tagName=%s", username, tagName)

	return r.s.DeleteTag(ctx, username, tagName)
}
