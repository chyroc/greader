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

func (r *GReader) EditSubscription(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	err := r.editSubscription(ctx, req)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, nil)
	}
}

func (r *GReader) editSubscription(ctx context.Context, req HttpReader) error {
	username, _ := r.getHeaderAuth(req)

	feedID := getFeedID(req.FormString("s"))
	action := req.FormString("ac")
	addTag := getUserLabelName(req.FormString("a"))
	removeTag := getUserLabelName(req.FormString("r"))
	title := req.FormString("t")
	r.log.Info(ctx, "[EditSubscription], username=%s, feedID=%s, action=%s, removeTag=%s, addTag=%s, title=%s", username, feedID, action, removeTag, addTag, title)

	if err := r.mustJson(req); err != nil {
		return err
	}

	switch action {
	case "unsubscribe":
		return r.backend.DeleteSubscription(ctx, username, feedID)
	case "edit":
		if addTag != "" || removeTag != "" {
			err := r.backend.UpdateSubscriptionTag(ctx, username, feedID, addTag, removeTag)
			if err != nil {
				return err
			}
		}

		if title != "" {
			err := r.backend.UpdateSubscriptionTitle(ctx, username, feedID, title)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
