package greader_api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// api: https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPICaller.swift#L557
//
// api path: /reader/api/0/stream/items/ids

const (
	readerStateRead    = "user/-/state/com.google/read"
	readerStateStarred = "user/-/state/com.google/starred"
)
const readingList = "user/-/state/com.google/reading-list"

func (r *Client) ListItemIDs(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	res, err := r.listItemIDs(ctx, req)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, res)
	}
}

func (r *Client) listItemIDs(ctx context.Context, req HttpReader) (*listEntryIDsResponse, error) {
	if err := r.mustJson(req); err != nil {
		return nil, err
	}

	s := req.QueryString("s")
	xt := req.QueryString("xt")
	continuation := req.QueryString("c")
	count, err := queryCount(req)
	if err != nil {
		return nil, err
	}
	// sinceTimeInterval := req.QueryString("ot")
	since := time.Now().Add(-time.Hour * 24 * 30)
	var readed *bool
	var starred *bool
	var feedID *string
	if s == readerStateStarred {
		// starred
		starred = &[]bool{true}[0]
	} else if s == readingList && xt == readerStateRead {
		// unread
		readed = &[]bool{false}[0]
	} else if s == readingList {
		// all for account
	} else if s != "" {
		// for feed, feedID = s
		feedID = &[]string{getFeedID(s)}[0]
	}
	continuationNew, ids, err := r.s.ListEntryIDs(ctx, readed, starred, feedID, since, count, continuation)
	if err != nil {
		return nil, err
	}

	return &listEntryIDsResponse{
		EntryIDs:     buildEntryIDs(ids),
		Continuation: continuationNew,
	}, nil
}

func queryCount(req HttpReader) (int64, error) {
	count := req.QueryString("n")
	if count == "" {
		return 0, fmt.Errorf("missing count(n)")
	}
	n, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid count(n:%s)", count)
	}
	return n, nil
}
