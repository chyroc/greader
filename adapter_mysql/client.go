package adapter_mysql

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/chyroc/greader/adapter_mysql/dal"
	"github.com/chyroc/greader/adapter_mysql/internal"
	"github.com/chyroc/greader/greader_api"
)

type MySQLClient struct {
	db  *dal.Client
	log greader_api.ILogger
}

var _ greader_api.IGReaderStore = (*MySQLClient)(nil)

func New(dsn string, logger greader_api.ILogger) (*MySQLClient, error) {
	db, err := newDB(dsn)
	if err != nil {
		return nil, err
	}
	return &MySQLClient{db: dal.New(db), log: logger}, nil
}

func (r *MySQLClient) Login(ctx context.Context, username, password string) (string, error) {
	return r.db.Login(username, internal.CalSha1(username+":"+password))
}

func (r *MySQLClient) ListTag(ctx context.Context, username string) ([]string, error) {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return nil, err
	}

	tagNames, err := r.db.ListUserTagNames(userID)
	if err != nil {
		return nil, err
	}

	return tagNames, nil
}

func (r *MySQLClient) RenameTag(ctx context.Context, username, oldName, newName string) error {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return err
	}

	return r.db.RenameUserFeedTag(userID, oldName, newName)
}

func (r *MySQLClient) DeleteTag(ctx context.Context, username, name string) error {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return err
	}

	return r.db.DeleteUserFeedTag(userID, name)
}

func (r *MySQLClient) ListSubscription(ctx context.Context, username string) ([]*greader_api.Subscription, error) {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return nil, err
	}

	userFeedPOs, err := r.db.ListUserFeed(userID)
	if err != nil {
		return nil, err
	}

	feedIDs := internal.Map(userFeedPOs, func(item *dal.ModelUserFeedRelation) int64 { return item.FeedID })

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

func (r *MySQLClient) AddSubscription(ctx context.Context, username, url string) (*greader_api.Subscription, error) {
	r.log.Info(ctx, "add subscription, username=%s, url=%s", username, url)

	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return nil, err
	}

	feed, err := feedParser.ParseURL(url)
	if err != nil {
		return nil, err
	}

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

func (r *MySQLClient) DeleteSubscription(ctx context.Context, username, feedID string) error {
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

func (r *MySQLClient) UpdateSubscriptionTitle(ctx context.Context, username, feedID, title string) error {
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

func (r *MySQLClient) UpdateSubscriptionTag(ctx context.Context, username, feedID string, addTag, removeTag string) error {
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

func (r *MySQLClient) LoadEntry(ctx context.Context, username string, entryIDs []string) ([]*greader_api.Entry, error) {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return nil, err
	}

	ids := internal.StringListToInt(entryIDs)

	entryPOs, err := r.db.MGetEntry(ids)
	if err != nil {
		return nil, err
	}

	feedIDs := []int64{}
	for _, v := range entryPOs {
		feedIDs = append(feedIDs, v.FeedID)
	}
	feedPOs, err := r.db.MGetFeed(feedIDs)
	if err != nil {
		return nil, err
	}

	tagNames, err := r.db.ListUserFeedTagNames(userID, feedIDs)
	if err != nil {
		return nil, err
	}

	res := internal.MapNoneEmpty(ids, func(id int64) *greader_api.Entry {
		item := entryPOs[id]
		if item == nil {
			return nil
		}
		feedPO := feedPOs[item.FeedID]
		if feedPO == nil {
			return nil
		}

		alternates := []*greader_api.AlternateLocation{{URL: item.URL}}
		categories := []string{}
		if tagNames[item.FeedID] != "" {
			categories = append(categories, tagNames[item.FeedID])
		}

		return &greader_api.Entry{
			ID:                 strconv.FormatInt(item.ID, 10),
			Title:              item.Title,
			Author:             item.Author,
			PublishedTimestamp: item.CreatedAt.Unix(),
			CrawledTimestamp:   strconv.FormatInt(item.CreatedAt.UnixMilli(), 10),
			TimestampUsec:      strconv.FormatInt(item.CreatedAt.UnixMicro(), 10),
			Summary:            &greader_api.EntrySummary{Content: ""}, // TODO
			Alternates:         alternates,
			Categories:         categories,
			Origin: &greader_api.EntryOrigin{
				StreamID: fmt.Sprintf("feed/%d", feedPO.ID),
				Title:    "",
			},
		}
	})

	return res, nil
}

func (r *MySQLClient) UpdateEntry(ctx context.Context, username string, entryIDs []string, read, starred *bool) error {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return err
	}

	ids := internal.StringListToInt(entryIDs)

	return r.db.UpdateUserEntryStatus(userID, ids, read, starred)
}

func (r *MySQLClient) ListEntryIDs(ctx context.Context, username string, readed, starred *bool, feedID *string, since time.Time, count int64, continuation string) (string, []int64, error) {
	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return "", nil, err
	}

	latestEntryID, err := r.db.GetUserEntryLatestID(userID)
	if err != nil {
		return "", nil, err
	}
	r.log.Info(ctx, "[ListEntryIDs] latest_entry_id=%d", latestEntryID)

	var feedIDs []int64
	if feedID != nil {
		id, err := strconv.ParseInt(*feedID, 10, 64)
		if err != nil {
			return "", nil, err
		}
		feedIDs = append(feedIDs, id)
	} else {
		feedIDs, err = r.db.ListUserFeedIDs(userID)
		if err != nil {
			return "", nil, err
		}
	}
	r.log.Info(ctx, "[ListEntryIDs] feed_ids=%+v", feedIDs)

	unreadEntryList, err := r.db.ListEntryByLatestID(feedIDs, latestEntryID, int(count))
	if err != nil {
		return "", nil, err
	}

	go func() {
		for _, v := range unreadEntryList {
			err = r.db.CreateUserEntry(userID, v.FeedID, v.ID)
			if err != nil {
				r.log.Error(ctx, "[ListEntryIDs] CreateUserEntry username=%s, feed_id=%d, err=%s", username, v.FeedID, err)
			}
		}
	}()

	pos, err := r.db.ListUserEntry(userID, readed, starred, feedID, int(count))
	if err != nil {
		return "", nil, err
	}

	entryIDs := internal.MapNoneEmpty(pos, func(item *dal.ModeUserEntryRelation) int64 { return item.EntryID })

	return "", entryIDs, nil
}

func (r *MySQLClient) ListFeedURL(ctx context.Context) ([]string, error) {
	r.log.Info(ctx, "[ListFeedURL]")
	return r.db.ListFeedURL()
}

func (r *MySQLClient) AddFeedEntry(ctx context.Context, feedURL string, entryList []*greader_api.Entry) error {
	r.log.Info(ctx, "[AddFeedEntry] feedURL=%s, entryList.len=%d", feedURL, len(entryList))

	feedPO, err := r.db.GetFeedByURL(feedURL)
	if err != nil {
		return err
	}

	for _, v := range entryList {
		url := ""
		for _, v := range v.Alternates {
			if v.URL != "" {
				url = v.URL
				break
			}
		}
		if url == "" {
			continue
		}
		entryPO := &dal.ModelEntry{
			FeedID: feedPO.ID,
			URL:    url,
			Title:  v.Title,
			Author: v.Author,
		}
		if err = r.db.CreateEntry(entryPO); err != nil {
			r.log.Error(ctx, "[AddFeedEntry] CreateEntry err=%s", err)
		}
	}

	return nil
}

func (r *MySQLClient) validAuth(ctx context.Context, username string) (int64, error) {
	if username == "" {
		return 0, fmt.Errorf("auth error")
	}
	user, _ := r.db.GetUser(username)
	if user == nil {
		return 0, fmt.Errorf("auth error")
	}
	return user.ID, nil
}
