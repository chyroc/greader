package pack

import (
	"github.com/chyroc/greader/greader_api"
	"github.com/chyroc/greader/mysql_backend/dal"
)

func EntryToModel(feedID int64, entryList []*greader_api.Entry) []*dal.ModelEntry {
	entryPOs := []*dal.ModelEntry{}
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
			FeedID: feedID,
			URL:    url,
			Title:  v.Title,
			Author: v.Author,
		}
		entryPOs = append(entryPOs, entryPO)
	}
	return entryPOs
}
