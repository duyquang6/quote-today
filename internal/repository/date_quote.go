package repository

import (
	"context"

	"github.com/duyquang6/quote-today/internal/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type dateQuoteRepo struct{}

// DateQuoteRepository provide interface interact with DateQuote model
type DateQuoteRepository interface {
	Create(ctx context.Context, tx *gorm.DB, purchase *model.DateQuote) error
	IncreaseLikeByOne(ctx context.Context, tx *gorm.DB, date datatypes.Date) error
	DecreaseLikeByOne(ctx context.Context, tx *gorm.DB, date datatypes.Date) error
	Get(ctx context.Context, tx *gorm.DB, date datatypes.Date) (model.DateQuote, error)
}

// NewDateQuoteRepository create DateQuoteRepository concrete object
func NewDateQuoteRepository() DateQuoteRepository {
	return &dateQuoteRepo{}
}

// Create new purchase with given DateQuote
func (s *dateQuoteRepo) Create(ctx context.Context, tx *gorm.DB, dateQuote *model.DateQuote) error {
	return tx.Select("Date", "QuoteID", "LikeCount", "CreatedAt", "UpdatedAt").Create(&dateQuote).Error
}

// IncreaseLikeByOne increase like count
func (s *dateQuoteRepo) IncreaseLikeByOne(ctx context.Context, tx *gorm.DB, date datatypes.Date) error {
	return tx.Exec(`
UPDATE date_quotes SET like_count = like_count + 1
WHERE date = ?`, date).Error
}

// DecreaseLikeByOne decrease like count
func (s *dateQuoteRepo) DecreaseLikeByOne(ctx context.Context, tx *gorm.DB, date datatypes.Date) error {
	return tx.Exec(`
UPDATE date_quotes SET like_count = like_count - 1
WHERE date = ?`, date).Error
}

// Get increase like count
func (s *dateQuoteRepo) Get(ctx context.Context, tx *gorm.DB, date datatypes.Date) (model.DateQuote, error) {
	res := model.DateQuote{Date: date}
	err := tx.Preload("Quote").First(&res).Error
	return res, err
}
