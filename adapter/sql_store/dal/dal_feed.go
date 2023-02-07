package dal

type ModelSubscription struct {
	Id          int64  `gorm:"column:id; primary_key; auto_increment"`
	FeedURL     string `gorm:"column:feed_url"`
	HomePageURL string `gorm:"column:home_url"`
}

func (ModelSubscription) TableName() string {
	return "feed"
}

func (r *Client) MGetSubscription(ids []int64) (map[int64]*ModelSubscription, error) {
	var pos []*ModelSubscription
	err := r.db.Where("id in (?)", ids).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	res := make(map[int64]*ModelSubscription)
	for _, v := range pos {
		res[v.Id] = v
	}
	return res, nil
}

func (r *Client) CreateSubscription(feedURL, homeURL string) (*ModelSubscription, error) {
	po := &ModelSubscription{
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
