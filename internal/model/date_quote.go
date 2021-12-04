package model

import (
	"gorm.io/datatypes"
	"time"
)

// DateQuote data model
type DateQuote struct {
	Date      datatypes.Date `gorm:"primarykey,type=date"`
	QuoteID   uint
	Quote     Quote
	LikeCount uint
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `sql:"DEFAULT:current_timestamp"`
}
