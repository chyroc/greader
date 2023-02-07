package greader_api

import (
	"context"
	"time"
)

func GetContextAuth(ctx context.Context) (string, string) {
	return getContext(ctx)
}

type IStore interface {
	Auth(ctx context.Context, username, password string) (string, error)

	ListTag(ctx context.Context) ([]string, error)
	RenameTag(ctx context.Context, oldName, newName string) error
	DeleteTag(ctx context.Context, name string) error

	ListSubscription(ctx context.Context) ([]*Subscription, error)
	AddSubscription(ctx context.Context, url string) (*Subscription, error)
	DeleteSubscription(ctx context.Context, feedID string) error
	RenameSubscription(ctx context.Context, feedID, title string) error

	ChangeSubscriptionTagging(ctx context.Context, feedID string, addTag, removeTag string) error

	LoadEntry(ctx context.Context, articleIDs []string) ([]*Entry, error)
	ListEntryIDs(ctx context.Context, readed, starred *bool, feedID *string, since time.Time, count int64, continuation string) (string, []int64, error)
	ChangeEntryStatus(ctx context.Context, articleIDs []string, read, starred *bool) error
}
