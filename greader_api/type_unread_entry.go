package greader_api

import "strconv"

// https://github.com/Ranchero-Software/NetNewsWire/blob/e59e29eb43804052424770b85d499006d6bde00c/Account/Sources/Account/ReaderAPI/ReaderAPIUnreadEntry.swift

type entryIDs struct {
	EntryID string `json:"id,omitempty"`
}

type listEntryIDsResponse struct {
	EntryIDs     []*entryIDs `json:"itemRefs,omitempty"`
	Continuation string      `json:"continuation,omitempty"`
}

func buildEntryIDs(ids []int64) []*entryIDs {
	res := make([]*entryIDs, 0, len(ids))
	for _, v := range ids {
		res = append(res, &entryIDs{EntryID: strconv.FormatInt(v, 10)})
	}
	return res
}
