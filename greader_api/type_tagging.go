package greader_api

// https://github.com/Ranchero-Software/NetNewsWire/blob/e59e29eb43804052424770b85d499006d6bde00c/Account/Sources/Account/ReaderAPI/ReaderAPITagging.swift

type ReaderAPITagging struct {
	TaggingID int    `json:"id"`
	FeedID    int    `json:"feed_id"`
	Name      string `json:"name"`
}

type ReaderAPIDeleteTagging struct {
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}
