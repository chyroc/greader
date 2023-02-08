package greader_api

import (
	"context"
	"net/http"
)

// api: https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPICaller.swift#L297

// api path: /reader/api/0/subscription/list

func (r *Client) ListSubscription(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	res, err := r.listSubscription(ctx, req)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, res)
	}
}

func (r *Client) listSubscription(ctx context.Context, req HttpReader) (*subscriptionList, error) {
	username, _ := r.getHeaderAuth(req)
	r.log.Info(ctx, "[ListSubscription] username=%s", username)

	if err := r.mustJson(req); err != nil {
		return nil, err
	}

	res, err := r.s.ListSubscription(ctx, username)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		v.reBuild()
	}

	return &subscriptionList{
		Subscriptions: res,
	}, nil
}
