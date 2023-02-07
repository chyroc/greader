package dal

func (r *Client) ListUserFeedTagNames(userID int64) ([]string, error) {
	var pos []*ModelUserFeedRelation
	err := r.db.Where("user_id = ?", userID).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	tags := make([]string, 0, len(pos))
	for _, v := range pos {
		tags = append(tags, v.TagName)
	}
	return tags, nil
}

func (r *Client) RenameUserFeedTag(userID int64, oldTag, newTag string) error {
	err := r.db.Model(&ModelUserFeedRelation{}).
		Where("user_id = ? and tag_name = ?", userID, oldTag).
		Update("tag_name", newTag).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Client) MoveUserFeedTag(userID, feedID int64, oldTag, newTag string) error {
	if newTag == "" {
		newTag = defaultTagName
	}
	err := r.db.Model(&ModelUserFeedRelation{}).
		Where("user_id = ? and feed_id = ? and tag_name = ?", userID, feedID, oldTag).
		Update("tag_name", newTag).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Client) DeleteUserFeedTag(userID int64, tagName string) error {
	err := r.db.
		Where("user_id = ? and tag_name = ?", userID, tagName).
		Delete(&ModelUserFeedRelation{}).Error
	if err != nil {
		return err
	}
	return nil
}
