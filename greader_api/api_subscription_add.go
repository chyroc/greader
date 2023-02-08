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
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, res)
	}
}

func (r *Client) addSubscription(ctx context.Context, req HttpReader) (*AddSubscriptionResult, error) {
	username := getContextUsername(ctx)
	feedURL := req.FormString("quickadd")
	r.log.Info(ctx, "[AddSubscription] username=%s, feedURL=%s", username, feedURL)

	feed, err := p.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}

	subscription, err := r.s.AddSubscription(ctx, username, feedURL, feed.Link, feed.Title)
	if err != nil {
		return nil, err
	} else if subscription == nil {
		return &AddSubscriptionResult{NumResults: 0}, nil
	}

	if err = r.s.AddFeedEntry(ctx, &username, feedURL, rssItemToEntry(feed.Items)); err != nil {
		return nil, err
	}

	return &AddSubscriptionResult{NumResults: 1, FeedID: buildFeedID(subscription.FeedID)}, nil
}
