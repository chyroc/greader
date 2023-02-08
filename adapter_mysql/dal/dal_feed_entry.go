package dal

import (
	"github.com/chyroc/greader/adapter_mysql/internal"
)

type ModelEntry struct {
	BaseModel

	FeedID int64  `gorm:"column:feed_id"`
	Title  string `gorm:"column:title"`
	URL    string `gorm:"column:url"`
	Author string `gorm:"column:author"`
}

func (ModelEntry) TableName() string {
	return "entry"
}

func (r *Client) MGetEntry(ids []int64) (map[int64]*ModelEntry, error) {
	var pos []*ModelEntry
	err := r.db.Where("id in (?)", ids).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	res := make(map[int64]*ModelEntry)
	for _, v := range pos {
		res[v.ID] = v
	}
	return res, nil
}

func (r *Client) CreateEntries(pos []*ModelEntry) error {
	return r.db.Clauses(ignoreInsertClause).Create(&pos).Error
}

func (r *Client) ListEntryByLatestID(feedIDs []int64, latestEntryID int64, limit int) ([]*ModelEntry, error) {
	feedIDs = internal.Unique(feedIDs)
	var pos []*ModelEntry
	err := r.db.
		Where("feed_id in (?) and id > ?", feedIDs, latestEntryID).
		Order("id desc").
		Limit(limit).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}

	return pos, nil
}
