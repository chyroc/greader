package adapter_mysql

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"

	dal2 "github.com/chyroc/greader/adapter_mysql/dal"
	"github.com/chyroc/greader/adapter_mysql/internal"
	"github.com/chyroc/greader/greader_api"
)

type Client struct {
	db  *dal2.Client
	log greader_api.ILogger
}

var _ greader_api.IGReaderStore = (*Client)(nil)

func New(dsn string, logger greader_api.ILogger) (*Client, error) {
	db, err := newDB(dsn)
	if err != nil {
		return nil, err
	}
	return &Client{db: dal2.New(db), log: logger}, nil
}

func (r *Client) Login(ctx context.Context, username, password string) (string, error) {
	return r.db.Login(username, password)
}

func (r *Client) ListTag(ctx context.Context, username string) ([]string, error) {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return nil, err
	}

	tagNames, err := r.db.ListUserFeedTagNames(userID)
	if err != nil {
		return nil, err
	}

	return tagNames, nil
}

func (r *Client) RenameTag(ctx context.Context, username, oldName, newName string) error {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return err
	}

	return r.db.RenameUserFeedTag(userID, oldName, newName)
}

func (r *Client) DeleteTag(ctx context.Context, username, name string) error {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return err
	}

	return r.db.DeleteUserFeedTag(userID, name)
}

func (r *Client) ListSubscription(ctx context.Context, username string) ([]*greader_api.Subscription, error) {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return nil, err
	}

	userFeedPOs, err := r.db.ListUserFeed(userID)
	if err != nil {
		return nil, err
	}

	feedIDs := internal.Map(userFeedPOs, func(item *dal2.ModelUserFeedRelation) int64 { return item.FeedID })

	feedPOMap, err := r.db.MGetFeed(feedIDs)
	if err != nil {
		return nil, err
	}

	res := make([]*greader_api.Subscription, 0, len(userFeedPOs))
	for _, v := range userFeedPOs {
		feedPO := feedPOMap[v.FeedID]
		if feedPO == nil {
			continue
		}
		res = append(res, &greader_api.Subscription{
			FeedID:      strconv.FormatInt(v.FeedID, 10),
			Name:        v.Title,
			Categories:  greader_api.BuildCategories([]string{v.TagName}),
			FeedURL:     feedPO.FeedURL,
			HomePageURL: feedPO.HomePageURL,
		})
	}

	return res, nil
}

var feedParser = gofeed.NewParser()

func (r *Client) AddSubscription(ctx context.Context, username, url string) (*greader_api.Subscription, error) {
	r.log.Info(ctx, "add subscription, username=%s, url=%s", username, url)

	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return nil, err
	}

	feed, err := feedParser.ParseURL(url)
	if err != nil {
		return nil, err
	}
	r.log.Info(ctx, "add subscription, feed=%s", internal.Json(feed))

	subscriptionPO, err := r.db.CreateFeed(url, feed.Link)
	if err != nil {
		return nil, err
	}

	err = r.db.CreateUserFeed(userID, subscriptionPO.ID, feed.Title)
	if err != nil {
		return nil, err
	}

	return &greader_api.Subscription{
		FeedID:      strconv.FormatInt(subscriptionPO.ID, 10),
		Name:        feed.Title,
		FeedURL:     subscriptionPO.FeedURL,
		HomePageURL: subscriptionPO.HomePageURL,
	}, nil
}

func (r *Client) DeleteSubscription(ctx context.Context, username, feedID string) error {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(feedID, 10, 64)
	if err != nil {
		return err
	}

	return r.db.DeleteUserFeed(userID, id)
}

func (r *Client) UpdateSubscriptionTitle(ctx context.Context, username, feedID, title string) error {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(feedID, 10, 64)
	if err != nil {
		return err
	}

	return r.db.UpdateUserFeedTitle(userID, id, title)
}

func (r *Client) UpdateSubscriptionTag(ctx context.Context, username, feedID string, addTag, removeTag string) error {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(feedID, 10, 64)
	if err != nil {
		return err
	}

	err = r.db.UpdateUserFeedTag(userID, id, addTag)
	if err != nil {
		return err
	}

	return nil
}

func (r *Client) LoadEntry(ctx context.Context, entryIDs []string) ([]*greader_api.Entry, error) {
	ids := internal.StringListToInt(entryIDs)

	pos, err := r.db.MGetEntry(ids)
	if err != nil {
		return nil, err
	}

	res := internal.MapNoneEmpty(ids, func(id int64) *greader_api.Entry {
		item := pos[id]
		if item == nil {
			return nil
		}
		return &greader_api.Entry{
			ID:     strconv.FormatInt(item.ID, 10),
			Title:  item.Title,
			Author: item.Author,
		}
	})

	return res, nil
}

func (r *Client) UpdateEntry(ctx context.Context, username string, entryIDs []string, read, starred *bool) error {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return err
	}

	ids := internal.StringListToInt(entryIDs)

	return r.db.UpdateUserEntryStatus(userID, ids, read, starred)
}

func (r *Client) ListEntryIDs(ctx context.Context, username string, readed, starred *bool, feedID *string, since time.Time, count int64, continuation string) (string, []int64, error) {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return "", nil, err
	}

	pos, err := r.db.ListUserEntry(userID, readed, starred, feedID)
	if err != nil {
		return "", nil, err
	}

	entryIDs := internal.MapNoneEmpty(pos, func(item *dal2.ModeUserEntryRelation) int64 { return item.EntryID })

	return "", entryIDs, nil
}

func (r *Client) validAuth(ctx context.Context, username string) (int64, error) {
	if username == "" {
		return 0, fmt.Errorf("auth error")
	}
	user, _ := r.db.GetUser(username)
	if user == nil {
		return 0, fmt.Errorf("auth error")
	}
	return user.ID, nil
}
