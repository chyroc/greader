package greader_api

import (
	"context"
	"net/http"
	"sort"
	"strconv"
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

	res, err := r.s.LoadEntry(ctx, entryIDs)
	if err != nil {
		return nil, err
	}

	sort.Slice(res, func(i, j int) bool {
		// TODO
		return true
	})

	return &loadEntryList{
		ID: readingList,
		// Updated: 0,
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
