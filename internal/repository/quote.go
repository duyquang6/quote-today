package repository

import (
	"context"

	"github.com/duyquang6/quote-today/internal/model"
	"gorm.io/gorm"
)

type quoteRepo struct{}

// QuoteRepository provide interface interact with Quote model
type QuoteRepository interface {
	GetRandomQuote(ctx context.Context, tx *gorm.DB) (model.Quote, error)
}

// NewQuoteRepository create QuoteRepository concrete object
func NewQuoteRepository() QuoteRepository {
	return &quoteRepo{}
}

func (s *quoteRepo) GetRandomQuote(ctx context.Context, tx *gorm.DB) (model.Quote, error) {
	res := model.Quote{}
	tx.Raw(`
		SELECT
			id, author, quote
		FROM
			quotes OFFSET floor(random() * (
				SELECT
					COUNT(*)
					FROM quotes))
		LIMIT 1;
`).Scan(&res)
	return res, nil
}
