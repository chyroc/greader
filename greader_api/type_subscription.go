package greader_api

import "strings"

// struct: https://github.com/Ranchero-Software/NetNewsWire/blob/e59e29eb43804052424770b85d499006d6bde00c/Account/Sources/Account/ReaderAPI/ReaderAPISubscription.swift#L34

/*
{
  "numResults": 0,
  "error": "Already subscribed! https://inessential.com/xml/rss.xml"
}
*/

type AddSubscriptionResult struct {
	NumResults int    `json:"numResults"`
	Error      string `json:"error"`
	FeedID     string `json:"streamId"`
}

/*
{
	"id": "feed/1",
	"title": "Questionable Content",
	"categories": [
	{
		"id": "user/-/label/Comics",
		"label": "Comics"
	}
	],
	"url": "http://www.questionablecontent.net/QCRSS.xml",
	"htmlUrl": "http://www.questionablecontent.net",
	"iconUrl": "https://rss.confusticate.com/f.php?24decabc"
}
*/

type Subscription struct {
	FeedID      string      `json:"id"`
	Name        string      `json:"title"`
	Categories  []*Category `json:"categories"`
	FeedURL     string      `json:"url"`
	HomePageURL string      `json:"htmlUrl"`
	IconURL     string      `json:"iconUrl"`
}

type Category struct {
	ID    string `json:"id"` // "user/-/label/未分类"
	Label string `json:"label"`
}

func BuildCategory(name string) *Category {
	return &Category{
		ID:    "user/-/label/" + name,
		Label: name,
	}
}

func BuildCategories(names []string) []*Category {
	categories := make([]*Category, len(names))
	for i, name := range names {
		categories[i] = BuildCategory(name)
	}
	return categories
}

type subscriptionList struct {
	Subscriptions []*Subscription `json:"subscriptions"`
}

func buildFeedID(id string) string {
	return "feed/" + id
}

func (r *Subscription) reBuild() {
	if !strings.HasPrefix(r.FeedID, "feed/") {
		r.FeedID = buildFeedID(r.FeedID)
	}
	for _, v := range r.Categories {
		if !strings.HasPrefix(v.ID, "user/-/label/") {
			v.ID = "user/-/label/" + v.ID
		}
	}
}
