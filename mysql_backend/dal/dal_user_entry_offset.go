package dal

import (
	"errors"

	"gorm.io/gorm"
)

type ModeUserEntryOffset struct {
	BaseModel
	UserID int64 `gorm:"column:user_id"`
	// FeedID int64 `gorm:"column:feed_id"`
	Latest int64 `gorm:"column:latest"`
}

func (ModeUserEntryOffset) TableName() string {
	return "user_entry_offset"
}

func (r *Client) GetUserEntryLatestID(userID int64) (int64, error) {
	var exist ModeUserEntryOffset
	err := r.db.Where("user_id = ?", userID).First(&exist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return exist.Latest, nil
}

func (r *Client) UpdateUserEntryOffset(userID, latest int64) error {
	// if exist
	var exist ModeUserEntryOffset
	err := r.db.Where("user_id = ?", userID).First(&exist).Error
	if err != nil {
		// create
		err = r.db.Create(&ModeUserEntryOffset{
			UserID: userID,
			Latest: latest,
		}).Error
		return err
	} else {
		// update
		err = r.db.Model(&ModeUserEntryOffset{}).
			Where("user_id = ?", userID).
			Update("latest", latest).Error
		return err
	}
}
