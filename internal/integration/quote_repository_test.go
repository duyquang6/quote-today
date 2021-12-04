//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/duyquang6/quote-today/internal/repository"
)

func (p *PostgresRepositoryTestSuite) TestMySqlQuoteRepository_GetRandomQuote() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	quoteRepo := repository.NewQuoteRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		quote, err := quoteRepo.GetRandomQuote(ctx, tx)
		p.Assert().NoError(err)
		p.Assert().True(len(quote.Quote) > 0)
	})
}
