package dal

type ModelUserFeedRelation struct {
	BaseModel
	UserID  int64  `gorm:"column:user_id"`
	FeedID  int64  `gorm:"column:feed_id"`
	TagName string `gorm:"column:tag_name"`
	Title   string `gorm:"column:title"`
}

func (ModelUserFeedRelation) TableName() string {
	return "user_feed"
}

func (r *Client) ListUserFeed(userID int64) ([]*ModelUserFeedRelation, error) {
	var pos []*ModelUserFeedRelation
	err := r.db.
		Where("user_id = ?", userID).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	return pos, nil
}

func (r *Client) ListUserFeedIDs(userID int64) ([]int64, error) {
	var ids []int64
	err := r.db.
		Model(&ModelUserFeedRelation{}).
		Where("user_id = ?", userID).
		Pluck("feed_id", &ids).Error
	return ids, err
}

func (r *Client) CreateUserFeed(userID, feedID int64, title string) error {
	err := r.db.Create(&ModelUserFeedRelation{
		UserID:  userID,
		FeedID:  feedID,
		Title:   title,
		TagName: defaultTagName,
	}).Error
	if err != nil {
		// TODO ignore dup
		return err
	}
	return nil
}

func (r *Client) DeleteUserFeed(userID, feedID int64) error {
	err := r.db.Where("user_id = ? and feed_id = ?", userID, feedID).
		Delete(&ModelUserFeedRelation{}).Error
	return err
}

func (r *Client) UpdateUserFeedTitle(userID, feedID int64, newTitle string) error {
	err := r.db.Model(&ModelUserFeedRelation{}).
		Where("user_id = ? and feed_id = ?", userID, feedID).
		Update("title", newTitle).Error
	return err
}
