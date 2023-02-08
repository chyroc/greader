package dal

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
	err := r.db.Create(&pos).Error
	return err
}
