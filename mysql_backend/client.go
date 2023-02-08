package mysql_backend

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/chyroc/greader/greader_api"
	"github.com/chyroc/greader/mysql_backend/dal"
	"github.com/chyroc/greader/mysql_backend/internal"
	"github.com/chyroc/greader/mysql_backend/pack"
)

type MySQLClient struct {
	db  *dal.Client
	log greader_api.ILogger
}

var _ greader_api.IGReaderBackend = (*MySQLClient)(nil)

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

	tagNames = internal.Unique(tagNames)

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

func (r *MySQLClient) AddSubscription(ctx context.Context, username, feedURL, homeURL, title string) (*greader_api.Subscription, error) {
	r.log.Info(ctx, "add subscription, username=%s, feedURL=%s", username, feedURL)

	userID, err := r.validAuth(ctx, username)
	if err != nil {
		return nil, err
	}

	subscriptionPO, err := r.db.CreateFeed(feedURL, homeURL)
	if err != nil {
		return nil, err
	}

	err = r.db.CreateUserFeed(userID, subscriptionPO.ID, title)
	if err != nil {
		return nil, err
	}

	return &greader_api.Subscription{
		FeedID:      strconv.FormatInt(subscriptionPO.ID, 10),
		Name:        title,
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

	pos, err := r.db.ListUserEntry(userID, readed, starred, feedID, int(count))
	if err != nil {
		return "", nil, err
	}

	entryIDs := internal.MapNoneEmpty(pos, func(item *dal.ModeUserEntryRelation) int64 { return item.EntryID })

	return "", entryIDs, nil
}

func (r *MySQLClient) ListFeedURL(ctx context.Context) ([]string, error) {
	// r.log.Info(ctx, "[ListFeedURL]")
	return r.db.ListFeedURL()
}

func (r *MySQLClient) AddFeedEntry(ctx context.Context, username *string, feedURL string, entryList []*greader_api.Entry) error {
	// r.log.Info(ctx, "[AddFeedEntry] feedURL=%s, entryList.len=%d", feedURL, len(entryList))
	var userID int64
	var err error
	if username != nil {
		if userID, err = r.validAuth(ctx, *username); err != nil {
			return err
		}
	}

	feedPO, err := r.db.GetFeedByURL(feedURL)
	if err != nil {
		return err
	}

	entryPOs := pack.EntryToModel(feedPO.ID, entryList)
	if err = r.db.CreateEntries(entryPOs); err != nil {
		r.log.Error(ctx, "[AddFeedEntry] CreateEntry err=%s", err)
	}

	if userID > 0 {
		return r.db.CreateUserEntries(pack.UserEntryToRelation([]int64{userID}, entryPOs))
	} else {
		userIDs, err := r.db.ListFeedUserIDs(feedPO.ID)
		if err != nil {
			return err
		}
		return r.db.CreateUserEntries(pack.UserEntryToRelation(userIDs, entryPOs))
	}

	return nil
}

func (r *MySQLClient) validAuth(ctx context.Context, username string) (int64, error) {
	if username == "" {
		return 0, fmt.Errorf("auth error, username empty")
	}
	user, _ := r.db.GetUser(username)
	if user == nil {
		return 0, fmt.Errorf("auth error, cannot get user")
	}
	return user.ID, nil
}
