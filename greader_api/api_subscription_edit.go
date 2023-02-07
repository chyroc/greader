package greader_api

import (
	"context"
	"net/http"
)

// api: https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPICaller.swift#L297

// api path: /reader/api/0/subscription/edit

// delete
// T=\(token)&s=\(subscriptionID)&ac=unsubscribe
//
// rename
// T=\(token)&s=\(subscriptionID)&ac=edit&t=\(encodedTitle)
//
// add tag
// T=\(token)&s=\(subscriptionID)&ac=edit&a=user/-/label/\(toLabel)
//
// delete tag
// T=\(token)&s=\(subscriptionID)&ac=edit&r=user/-/label/\(fromLabel)
//
// move tag
// T=\(token)&s=\(subscriptionID)&ac=edit&r=user/-/label/\(fromLabel)&a=user/-/label/\(toLabel)

func (r *Client) EditSubscription(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	err := r.editSubscription(ctx, req)
	if err != nil {
		r.renderErr(writer, err)
	} else {
		r.renderData(writer, map[string]interface{}{})
	}
}

func (r *Client) editSubscription(ctx context.Context, req HttpReader) error {
	if err := r.mustJson(req); err != nil {
		return err
	}

	feedID := req.FormString("s")

	switch req.FormString("ac") {
	case "unsubscribe":
		return r.s.DeleteSubscription(ctx, feedID)
	case "edit":
		addTag := req.FormString("a")
		removeTag := req.FormString("r")
		title := req.FormString("t")

		if addTag != "" || removeTag != "" {
			err := r.s.ChangeSubscriptionTagging(ctx, feedID, addTag, removeTag)
			if err != nil {
				return err
			}
		}

		if title != "" {
			err := r.s.RenameSubscription(ctx, feedID, title)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
