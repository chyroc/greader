package dal

import (
	"errors"
)

type ModelUser struct {
	BaseModel
	Username string `gorm:"column:username"`
	Hash     string `gorm:"column:hash"`
}

func (ModelUser) TableName() string {
	return "user"
}

func (r *Client) GetUser(username string) (*ModelUser, error) {
	var pos []*ModelUser
	err := r.db.Where("username = ?", username).Find(&pos).Error
	if err != nil {
		return nil, err
	} else if len(pos) == 0 {
		return nil, nil
	}
	return pos[0], nil
}

func (r *Client) Login(username, hash string) (string, error) {
	var pos []*ModelUser
	if err := r.db.
		Where("username = ? and hash = ?", username, hash).
		Find(&pos).Error; err != nil {
		return "", err
	} else if len(pos) == 0 {
		return "", errors.New("invalid userID or password")
	}
	return hash, nil
}
