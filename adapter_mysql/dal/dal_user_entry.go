package dal

type ModeUserEntryRelation struct {
	Id      int64 `gorm:"column:id; primary_key; auto_increment"`
	UserID  int64 `gorm:"column:user_id; index:username"`
	FeedID  int64 `gorm:"column:feed_id"`
	EntryID int64 `gorm:"column:entry_id"`
	Readed  bool  `gorm:"column:readed; default:false"`
	Starred bool  `gorm:"column:starred; default:false"`
}

func (ModeUserEntryRelation) TableName() string {
	return "user_entry"
}

func (r *Client) UpdateUserEntryStatus(userID int64, ids []int64, read, star *bool) error {
	updated := map[string]interface{}{}
	if read != nil {
		updated["read"] = *read
	}
	if star != nil {
		updated["star"] = *star
	}

	if len(updated) > 0 {
		err := r.db.Model(&ModeUserEntryRelation{}).Where("user_id = ? and entry_id in (?)", userID, ids).Updates(updated).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Client) ListUserEntry(userID int64, readed, starred *bool, feedID *string) ([]*ModeUserEntryRelation, error) {
	db := r.db.Where("user_id = ?", userID)
	if readed != nil {
		db = db.Where("readed = ?", *readed)
	}
	if starred != nil {
		db = db.Where("starred = ?", *starred)
	}
	if feedID != nil {
		db = db.Where("feed_id = ?", *feedID)
	}

	var pos []*ModeUserEntryRelation
	err := db.Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return pos, nil
}
