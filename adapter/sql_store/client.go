package sql_store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/chyroc/greader/adapter/sql_store/dal"
	"github.com/chyroc/greader/adapter/sql_store/internal"
	"github.com/chyroc/greader/greader_api"
)

type Client struct {
	db  *dal.Client
	log greader_api.ILogger
}

var _ greader_api.IStore = (*Client)(nil)

func New(dsn string) (*Client, error) {
	db, err := newDB(dsn)
	if err != nil {
		return nil, err
	}
	return &Client{db: dal.New(db), log: greader_api.NewDefaultLogger()}, nil
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func (r *Client) Auth(ctx context.Context, username, password string) (string, error) {
	return r.db.CheckAuth(username, password)
}

func (r *Client) ListTag(ctx context.Context) ([]string, error) {
	userID, err := r.validAuth(ctx)
	if err != nil {
		return nil, err
	}

	res, err := r.db.ListUserFeedTagNames(userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Client) RenameTag(ctx context.Context, oldName, newName string) error {
	userID, err := r.validAuth(ctx)
	if err != nil {
		return err
	}

	return r.db.RenameUserFeedTag(userID, oldName, newName)
}

func (r *Client) DeleteTag(ctx context.Context, name string) error {
	userID, err := r.validAuth(ctx)
	if err != nil {
		return err
	}

	return r.db.DeleteUserFeedTag(userID, name)
}

func (r *Client) ListSubscription(ctx context.Context) ([]*greader_api.Subscription, error) {
	userID, err := r.validAuth(ctx)
	if err != nil {
		return nil, err
	}

	pos, err := r.db.ListUserSubscription(userID)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(pos))
	for _, v := range pos {
		ids = append(ids, v.FeedID)
	}

	feeds, err := r.db.MGetSubscription(ids)
	if err != nil {
		return nil, err
	}

	feedIDs := []int64{}
	for _, v := range feeds {
		feedIDs = append(feedIDs, v.Id)
	}

	res := make([]*greader_api.Subscription, 0, len(pos))
	for _, v := range pos {
		feed := feeds[v.FeedID]
		if feed == nil {
			continue
		}
		res = append(res, &greader_api.Subscription{
			FeedID:      strconv.FormatInt(v.FeedID, 10),
			Name:        v.Title,
			Categories:  greader_api.BuildCategories([]string{v.TagName}),
			FeedURL:     feed.FeedURL,
			HomePageURL: feed.HomePageURL,
		})
	}

	return res, nil
}

var feedParser = gofeed.NewParser()

func (r *Client) AddSubscription(ctx context.Context, url string) (*greader_api.Subscription, error) {
	r.log.Info(ctx, "add subscription, url=%s", url)

	userID, err := r.validAuth(ctx)
	if err != nil {
		return nil, err
	}
	r.log.Info(ctx, "add subscription, user=%d", userID)

	feed, err := feedParser.ParseURL(url)
	if err != nil {
		return nil, err
	}
	r.log.Info(ctx, "add subscription, feed=%s", internal.Json(feed))

	subscriptionPO, err := r.db.CreateSubscription(url, feed.Link)
	if err != nil {
		return nil, err
	}

	err = r.db.CreateUserSubscription(userID, subscriptionPO.Id, feed.Title)
	if err != nil {
		return nil, err
	}

	return &greader_api.Subscription{
		FeedID:      strconv.FormatInt(subscriptionPO.Id, 10),
		Name:        feed.Title,
		Categories:  nil,
		FeedURL:     subscriptionPO.FeedURL,
		HomePageURL: subscriptionPO.HomePageURL,
	}, nil
}

func (r *Client) DeleteSubscription(ctx context.Context, feedID string) error {
	userID, err := r.validAuth(ctx)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(feedID, 10, 64)
	if err != nil {
		return err
	}

	return r.db.DeleteUserSubscription(userID, id)
}

func (r *Client) RenameSubscription(ctx context.Context, feedID, title string) error {
	userID, err := r.validAuth(ctx)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(feedID, 10, 64)
	if err != nil {
		return err
	}

	return r.db.RenameUserSubscription(userID, id, title)
}

func (r *Client) ChangeSubscriptionTagging(ctx context.Context, feedID string, addTag, removeTag string) error {
	userID, err := r.validAuth(ctx)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(feedID, 10, 64)
	if err != nil {
		return err
	}

	if removeTag != "" {
		err = r.db.MoveUserFeedTag(userID, id, removeTag, addTag)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Client) LoadEntry(ctx context.Context, articleIDs []string) ([]*greader_api.Entry, error) {
	ids := []int64{}
	for _, v := range articleIDs {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	pos, err := r.db.MGetEntry(ids)
	if err != nil {
		return nil, err
	}

	res := make([]*greader_api.Entry, 0, len(pos))
	for _, v := range pos {
		res = append(res, &greader_api.Entry{
			ArticleID: strconv.FormatInt(v.ID, 10),
			Title:     v.Title,
			Author:    v.Author,
		})
	}

	return res, nil
}

func (r *Client) ChangeEntryStatus(ctx context.Context, articleIDs []string, read, starred *bool) error {
	userID, err := r.validAuth(ctx)
	if err != nil {
		return err
	}

	ids := []int64{}
	for _, v := range articleIDs {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	return r.db.UpdateUserEntryStatus(userID, ids, read, starred)
}

func (r *Client) ListEntryIDs(ctx context.Context, readed, starred *bool, feedID *string, since time.Time, count int64, continuation string) (string, []int64, error) {
	userID, err := r.validAuth(ctx)
	if err != nil {
		return "", nil, err
	}

	pos, err := r.db.ListUserEntry(userID, readed, starred, feedID)
	if err != nil {
		return "", nil, err
	}

	res := make([]int64, 0, len(pos))
	for _, v := range pos {
		res = append(res, v.EntryID)
	}

	return "", res, nil
}

func (r *Client) validAuth(ctx context.Context) (int64, error) {
	username, _ := greader_api.GetContextAuth(ctx)
	if username == "" {
		return 0, fmt.Errorf("auth error")
	}
	user, _ := r.db.GetUser(username)
	if user == nil {
		return 0, fmt.Errorf("auth error")
	}
	return user.ID, nil
}
