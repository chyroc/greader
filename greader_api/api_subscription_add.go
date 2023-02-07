package greader_api

import (
	"context"
	"net/http"
)

// api: https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPICaller.swift#L325

// api path: /reader/api/0/subscription/quickadd

// args: T=\(token)&quickadd=\(encodedFeedURL)

func (r *Client) AddSubscription(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	res, err := r.addSubscription(ctx, req)
	if err != nil {
		r.renderErr(writer, err)
	} else {
		r.renderData(writer, res)
	}
}

func (r *Client) addSubscription(ctx context.Context, req HttpReader) (*AddSubscriptionResult, error) {
	feedURL := req.FormString("quickadd")

	subscription, err := r.s.AddSubscription(ctx, feedURL)
	if err != nil {
		return nil, err
	} else if subscription == nil {
		return &AddSubscriptionResult{NumResults: 0}, nil
	}

	return &AddSubscriptionResult{NumResults: 1, FeedID: subscription.FeedID}, nil
}
