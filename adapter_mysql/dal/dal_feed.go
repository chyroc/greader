package dal

type ModelFeed struct {
	BaseModel
	FeedURL     string `gorm:"column:feed_url"`
	HomePageURL string `gorm:"column:home_url"`
}

func (ModelFeed) TableName() string {
	return "feed"
}

func (r *Client) MGetFeed(ids []int64) (map[int64]*ModelFeed, error) {
	var pos []*ModelFeed
	err := r.db.Where("id in (?)", ids).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	res := make(map[int64]*ModelFeed)
	for _, v := range pos {
		res[v.ID] = v
	}
	return res, nil
}

func (r *Client) CreateFeed(feedURL, homeURL string) (*ModelFeed, error) {
	po := &ModelFeed{
		FeedURL:     feedURL,
		HomePageURL: homeURL,
	}
	err := r.db.Create(po).Error
	if err == nil {
		return po, nil
	}

	r.db.Where("feed_url = ?", feedURL).Find(&po)

	return po, nil
}
