package greader_api

import (
	"context"
	"net/http"
)

// api: https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPICaller.swift#L667
//
// arg: T=\(token)&\(idsToFetch)&\(actionIndicator)=\(state.rawValue)
//
// i=tag:google.com,2005:reader/item/\(idHexString)

func (r *GReader) EditSubscriptionStatus(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	err := r.editSubscriptionStatus(ctx, req)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, nil)
	}
}

func (r *GReader) editSubscriptionStatus(ctx context.Context, req HttpReader) error {
	username, _ := r.getHeaderAuth(req)

	entryIDs := getEntryHexIDs(req.FormList("i"))
	add := req.FormString("a")    // add
	remove := req.FormString("r") // remove
	r.log.Info(ctx, "[EditSubscriptionStatus] username=%s, entryIDs=%+v, add=%s, remove=%s", username, entryIDs, add, remove)

	state := add
	if state == "" {
		state = remove
	}
	isAdd := add != ""

	var starred *bool
	var readed *bool

	switch state {
	case readerStateRead:
		if isAdd {
			// delete unread
			readed = &[]bool{true}[0]
		} else {
			// add unread
			readed = &[]bool{false}[0]
		}
	case readerStateStarred:
		if isAdd {
			// add starred
			starred = &[]bool{true}[0]
		} else {
			// delete starred
			starred = &[]bool{false}[0]
		}
	}

	return r.backend.UpdateEntry(ctx, username, entryIDs, readed, starred)
}
