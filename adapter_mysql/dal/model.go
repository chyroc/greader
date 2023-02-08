package dal

import (
	"time"

	"gorm.io/gorm"
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
