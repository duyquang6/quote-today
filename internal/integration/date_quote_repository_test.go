//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/duyquang6/quote-today/internal/model"
	"github.com/duyquang6/quote-today/internal/repository"
	"gorm.io/datatypes"
	"time"
)

func (p *PostgresRepositoryTestSuite) TestMySqlDateQuoteRepository_Create() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	dateQuoteRepo := repository.NewDateQuoteRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		dateQuote := model.DateQuote{
			Date:      datatypes.Date(time.Now().UTC()),
			QuoteID:   1,
			LikeCount: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := dateQuoteRepo.Create(ctx, db, &dateQuote)
		p.Assert().NoError(err)
	})
}

func (p *PostgresRepositoryTestSuite) TestMySqlDateQuoteRepository_IncreaseLikeByOne() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	dateQuoteRepo := repository.NewDateQuoteRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		current := time.Now().UTC()
		dateQuote := model.DateQuote{
			Date:      datatypes.Date(current),
			QuoteID:   1,
			LikeCount: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := dateQuoteRepo.Create(ctx, db, &dateQuote)
		p.Assert().NoError(err)

		err = dateQuoteRepo.IncreaseLikeByOne(ctx, db, dateQuote.Date)
		p.Assert().NoError(err)

		res, err := dateQuoteRepo.Get(ctx, tx, dateQuote.Date)
		p.Assert().NoError(err)
		p.Assert().Equal(uint(2), res.LikeCount)
	})
}

func (p *PostgresRepositoryTestSuite) TestMySqlDateQuoteRepository_DecreaseLikeByOne() {
	ctx := context.TODO()
	db := p.env.Database().GetDB()
	dateQuoteRepo := repository.NewDateQuoteRepository()

	p.Run("Success", func() {
		tx := db.Begin()
		defer tx.Rollback()
		current := time.Now().UTC()
		dateQuote := model.DateQuote{
			Date:      datatypes.Date(current),
			QuoteID:   1,
			LikeCount: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := dateQuoteRepo.Create(ctx, db, &dateQuote)
		p.Assert().NoError(err)

		err = dateQuoteRepo.DecreaseLikeByOne(ctx, db, dateQuote.Date)
		p.Assert().NoError(err)

		res, err := dateQuoteRepo.Get(ctx, tx, dateQuote.Date)
		p.Assert().NoError(err)
		p.Assert().Equal(uint(0), res.LikeCount)
	})
}
