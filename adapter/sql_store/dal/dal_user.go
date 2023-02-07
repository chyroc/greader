package dal

import (
	"errors"
	"log"

	"github.com/chyroc/greader/adapter/sql_store/internal"
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

func (r *Client) CheckAuth(username, password string) (hash string, err error) {
	defer func() {
		if err != nil {
			log.Println("[CheckAuth] fail", "username:", username, "err:", err)
		} else {
			log.Println("[CheckAuth] success", "username:", username)
		}
	}()
	hash = internal.CalSha1(username + ":" + password)

	pos := []*ModelUser{}
	if err = r.db.Where("username = ? and hash = ?", username, hash).Find(&pos).Error; err != nil {
		return "", err
	} else if len(pos) == 0 {
		return "", errors.New("invalid userID or password")
	}
	return hash, nil
}
