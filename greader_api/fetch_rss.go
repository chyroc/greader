package greader_api

import (
	"context"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"
)

var p = gofeed.NewParser()

func (r *Client) FetchRss() error {
	ctx := context.Background()

	feedURLs, err := r.s.ListFeedURL(ctx)
	if err != nil {
		r.fetchLogger.Error(ctx, "[FetchRss] ListFeedURL fail, err=%s", err)
		return err
	}

	for _, feedURL := range feedURLs {
		feed, err := p.ParseURL(feedURL)
		if err != nil {
			r.fetchLogger.Error(ctx, "[FetchRss] ParseURL fail, url=%s, err=%s", feedURL, err)
			continue
		}

		err = r.s.AddFeedEntry(ctx, nil, feedURL, rssItemToEntry(feed.Items))
		if err != nil {
			r.fetchLogger.Error(ctx, "[FetchRss] AddFeedEntry fail, url=%s, err=%s", feedURL, err)
			continue
		}
	}

	return nil
}

func rssItemToEntry(items []*gofeed.Item) []*Entry {
	res := []*Entry{}
	for _, item := range items {
		author := ""
		for i, v := range item.Authors {
			if i > 0 {
				author += ","
			}
			author += v.Name
		}
		content := item.Content
		if content == "" {
			content = item.Description
		}
		alternates := []*AlternateLocation{}
		for _, v := range append([]string{item.Link}, item.Links...) {
			if v != "" {
				alternates = append(alternates, &AlternateLocation{URL: v})
			}
		}
		now := time.Now()
		publish := time.Now()
		if item.PublishedParsed != nil {
			publish = *item.PublishedParsed
		}
		res = append(res, &Entry{
			Title:              item.Title,
			Author:             author,
			PublishedTimestamp: publish.Unix(),                         // 1554845280
			CrawledTimestamp:   strconv.FormatInt(now.UnixMilli(), 10), // 1559362260113
			TimestampUsec:      strconv.FormatInt(now.UnixMicro(), 10), // 1559362260113787
			Summary:            &EntrySummary{Content: content},
			Alternates:         alternates,
		})
	}

	return res
}
