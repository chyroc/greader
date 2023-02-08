package greader_api

// https://github.com/Ranchero-Software/NetNewsWire/blob/e59e29eb43804052424770b85d499006d6bde00c/Account/Sources/Account/ReaderAPI/ReaderAPIEntry.swift

/*
{
  "id": "tag:google.com,2005:reader/item/00058a3b5197197b",
  "crawlTimeMsec": "1559362260113",
  "timestampUsec": "1559362260113787",
  "published": 1554845280,
  "title": "",
  "summary": {
    "content": "\n<p>Found an old screenshot of NetNewsWire 1.0 for iPhone!</p>\n\n<p><img src=\"https://nnw.ranchero.com/uploads/2019/c07c0574b1.jpg\" alt=\"Netnewswire 1.0 for iPhone screenshot showing the list of feeds.\" title=\"NewsGator got renamed to Sitrion, years later, and then renamed again as Limeade.\" border=\"0\" width=\"260\" height=\"320\"></p>\n"
  },
  "alternate": [
    {
      "href": "https://nnw.ranchero.com/2019/04/09/found-an-old.html"
    }
  ],
  "categories": [
    "user/-/state/com.google/reading-list",
    "user/-/label/Uncategorized"
  ],
  "origin": {
    "streamId": "feed/130",
    "title": "NetNewsWire"
  }
}
*/

type Entry struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`

	PublishedTimestamp int64  `json:"published"`
	CrawledTimestamp   string `json:"crawlTimeMsec"`
	TimestampUsec      string `json:"timestampUsec"`

	Summary    *EntrySummary        `json:"summary"`
	Alternates []*AlternateLocation `json:"alternate"`
	Categories []string             `json:"categories"`
	Origin     *EntryOrigin         `json:"origin"`
}

type EntrySummary struct {
	Content string `json:"content,omitempty"`
}

type AlternateLocation struct {
	URL string `json:"href,omitempty"`
}

type EntryOrigin struct {
	StreamID string `json:"streamId,omitempty"`
	Title    string `json:"title,omitempty"`
}

type loadEntryList struct {
	ID      string   `json:"id"`
	Updated int      `json:"updated"`
	Entries []*Entry `json:"items"`
}
