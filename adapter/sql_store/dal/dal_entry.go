package dal

type ModelEntry struct {
	BaseModel

	FeedID int64  `gorm:"column:feed_id"`
	Title  string `gorm:"column:title"`
	Author string `gorm:"column:author"`
}

func (ModelEntry) TableName() string {
	return "entry"
}

func (r *Client) MGetEntry(ids []int64) (map[int64]*ModelEntry, error) {
	pos := []*ModelEntry{}
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

func (r *Client) CreateEntry(feedID int64, title, author string) (*ModelEntry, error) {
	po := &ModelEntry{
		FeedID: feedID,
		Title:  title,
		Author: author,
	}
	err := r.db.Create(po).Error
	return po, err
}
