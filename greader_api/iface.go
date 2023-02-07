package greader_api

import (
	"context"
	"time"
)

type ListEntryType int

const (
	ListEntryTypeStarred ListEntryType = 1
	ListEntryTypeUnread  ListEntryType = 2
	ListEntryTypeFeed    ListEntryType = 3
	ListEntryTypeAll     ListEntryType = 4
)

type IStore interface {
	Auth(ctx context.Context, username, password string) (string, error)

	ListTag(ctx context.Context) ([]string, error)
	RenameTag(ctx context.Context, oldName, newName string) error
	DeleteTag(ctx context.Context, name string) error

	ListSubscription(ctx context.Context) ([]*Subscription, error)
	AddSubscription(ctx context.Context, url string) (*Subscription, error)
	DeleteSubscription(ctx context.Context, feedID string) error
	RenameSubscription(ctx context.Context, feedID, title string) error
	ChangeSubscriptionTagging(ctx context.Context, feedID string, addTags, removeTags string) error

	LoadEntry(ctx context.Context, articleIDs []string) ([]*Entry, error)
	ListEntryIDs(ctx context.Context, typ ListEntryType, feedID string, since time.Time, count int64, continuation string) (string, []int64, error)
}
