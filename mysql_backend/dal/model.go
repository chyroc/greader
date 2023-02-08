package dal

import (
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseModel struct {
	ID        int64          `gorm:"column:id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

type Client struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Client {
	return &Client{db}
}

const defaultTagName = "default"

var ignoreInsertClause = clause.Insert{Modifier: "IGNORE"}

func isDuplicateErr(err error) bool {
	if err == nil {
		return false
	}
	m, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return m.Number == 1062
}
