package dal

import "github.com/chyroc/greader/adapter/sql_store/internal"

func (r *Client) ListUserFeedTagNames(userID int64) ([]string, error) {
	var pos []*ModelUserFeedRelation
	err := r.db.Where("user_id = ?", userID).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	tags := internal.Map(pos, func(i *ModelUserFeedRelation) string { return i.TagName })
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

func (r *Client) UpdateUserFeedTag(userID, feedID int64, newTag string) error {
	if newTag == "" {
		newTag = defaultTagName
	}
	err := r.db.Model(&ModelUserFeedRelation{}).
		Where("user_id = ? and feed_id = ?", userID, feedID).
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
