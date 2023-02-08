package dal

import (
	"github.com/chyroc/greader/adapter_mysql/internal"
)

type ModeUserEntryRelation struct {
	BaseModel
	UserID  int64 `gorm:"column:user_id"`
	FeedID  int64 `gorm:"column:feed_id"`
	EntryID int64 `gorm:"column:entry_id"`
	Readed  bool  `gorm:"column:readed"`
	Starred bool  `gorm:"column:starred"`
}

func (ModeUserEntryRelation) TableName() string {
	return "user_entry"
}

func (r *Client) UpdateUserEntryStatus(userID int64, ids []int64, read, star *bool) error {
	ids = internal.Unique(ids)
	updated := map[string]interface{}{}
	if read != nil {
		updated["readed"] = *read
	}
	if star != nil {
		updated["starred"] = *star
	}

	if len(updated) > 0 {
		err := r.db.Model(&ModeUserEntryRelation{}).Where("user_id = ? and entry_id in (?)", userID, ids).Updates(updated).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Client) ListUserEntry(userID int64, readed, starred *bool, feedID *string, count int) ([]*ModeUserEntryRelation, error) {
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
	if count > 0 {
		db = db.Limit(count)
	}

	var pos []*ModeUserEntryRelation
	err := db.Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return pos, nil
}

func (r *Client) CreateUserEntries(pos []*ModeUserEntryRelation) error {
	return r.db.Clauses(ignoreInsertClause).Create(&pos).Error
}
