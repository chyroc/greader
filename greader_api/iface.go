package greader_api

import (
	"context"
	"time"
)

func GetContextAuth(ctx context.Context) (string, string) {
	return getContext(ctx)
}

type IGReaderStore interface {
	// Login check username and password, return token
	Login(ctx context.Context, username, password string) (string, error)

	// tag

	// ListTag list all tags
	ListTag(ctx context.Context, username string) ([]string, error)
	// RenameTag rename tag, from oldName to newName
	RenameTag(ctx context.Context, username, oldName, newName string) error
	// DeleteTag delete tag, and remove all subscription of this tag
	DeleteTag(ctx context.Context, username, name string) error

	// subscription

	// ListSubscription list all subscription
	ListSubscription(ctx context.Context, username string) ([]*Subscription, error)
	// AddSubscription add subscription
	AddSubscription(ctx context.Context, username string, feedURL, homeURL, title string) (*Subscription, error)
	// DeleteSubscription delete subscription
	DeleteSubscription(ctx context.Context, username string, feedID string) error
	// UpdateSubscriptionTitle update subscription's title
	UpdateSubscriptionTitle(ctx context.Context, username string, feedID, title string) error
	// UpdateSubscriptionTag update subscription's tag
	UpdateSubscriptionTag(ctx context.Context, username string, feedID string, addTag, removeTag string) error

	// user's entry

	// ListEntryIDs list entry's id list
	ListEntryIDs(ctx context.Context, username string, readed, starred *bool, feedID *string, since time.Time, count int64, continuation string) (string, []int64, error)
	// UpdateEntry update entry's status
	UpdateEntry(ctx context.Context, username string, entryIDs []string, read, starred *bool) error
	// LoadEntry load entry by entry ids
	LoadEntry(ctx context.Context, username string, entryIDs []string) ([]*Entry, error)

	// global

	// LoadEntry load entry by entry ids
	// ListFeedURL list all feed url
	ListFeedURL(ctx context.Context) ([]string, error)
	// AddFeedEntry add feed entry
	AddFeedEntry(ctx context.Context, username *string, feedURL string, entryList []*Entry) error
}
