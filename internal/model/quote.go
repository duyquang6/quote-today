package model

import (
	"time"
)

// Quote data model
type Quote struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	Author    string
	Quote     string
}
