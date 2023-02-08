package greader_api

import (
	"context"
	"net/http"
)

// api: https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPICaller.swift#L217

// api path: /reader/api/0/rename-tag
//
// arg: T=\(token)&s=\(oldTagName)&dest=\(newTagName)
//
// &s=user/-/label/%E5%88%86%E7%B1%BB2&dest=user/-/label/%E5%88%86%E7%B1%BB2asdf

func (r *Client) TagRename(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	username, _ := r.getHeaderAuth(req)

	oldTagName := getUserLabelName(req.FormString("s"))
	newTagName := getUserLabelName(req.FormString("dest"))
	r.log.Info(ctx, "[TagRename] username=%s, rename: %s -> %s", username, oldTagName, newTagName)

	err := r.s.RenameTag(ctx, username, oldTagName, newTagName)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, nil)
	}
}
