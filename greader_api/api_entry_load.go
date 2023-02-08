package greader_api

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// api: https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPICaller.swift#L502

// api path:
// /reader/api/0/stream/items/contents

// T=\(token)&output=json&\(idsToFetch)
// i=tag:google.com,2005:reader/item/\(idHexString)

func (r *Client) LoadItem(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	res, err := r.loadItem(ctx, req)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, res)
	}
}

func (r *Client) loadItem(ctx context.Context, req HttpReader) (*loadEntryList, error) {
	username := getContextUsername(ctx)
	entryIDs := getEntryHexIDs(req.FormList("i"))
	r.log.Info(ctx, "[LoadItem], username=%s, entryIDs=%+v", username, entryIDs)

	if err := r.mustJson(req); err != nil {
		return nil, err
	}

	res, err := r.s.LoadEntry(ctx, username, entryIDs)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		id, err := strconv.ParseInt(v.ID, 10, 64)
		if err == nil {
			v.ID = "tag:google.com,2005:reader/item/" + intToHex16(id)
		}
		tmp := []string{"user/-/state/com.google/reading-list"}
		for _, v := range v.Categories {
			tmp = append(tmp, buildUserLabelName(v))
		}
		v.Categories = tmp
	}

	return &loadEntryList{
		ID:      readingList,
		Updated: time.Now().Unix(),
		Entries: res,
	}, nil
}

func getEntryHexIDs(list []string) []string {
	entryIDs := []string{}
	for _, v := range list {
		idHex := getTaggedItemHexID(v)
		id, err := hex16ToInt(idHex)
		if err != nil {
			// ignore
		} else {
			entryIDs = append(entryIDs, strconv.FormatInt(id, 10))
		}
	}
	return entryIDs
}
