package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/datatypes"

	"github.com/duyquang6/quote-today/internal/database"
	"github.com/duyquang6/quote-today/internal/model"
	"github.com/duyquang6/quote-today/internal/repository"
	"github.com/duyquang6/quote-today/pkg/dto"
	"github.com/duyquang6/quote-today/pkg/exception"
	"github.com/duyquang6/quote-today/pkg/logging"
)

const (
	purchaseServiceLoggingFmt = "PurchaseService.%s"
)

// DateQuoteService provide purchase service functionality
type DateQuoteService interface {
	GetRandomQuoteInsertIfNotExist(ctx context.Context) (dto.GetQuoteResponse, error)
	Like(ctx context.Context, request dto.LikeRequest) (dto.LikeResponse, error)
	Dislike(ctx context.Context, request dto.LikeRequest) (dto.LikeResponse, error)
}

type dateQuoteSvc struct {
	dbFactory     database.DBFactory
	quoteRepo     repository.QuoteRepository
	dateQuoteRepo repository.DateQuoteRepository
}

// NewDateQuoteService create concrete object which implement UserService
func NewDateQuoteService(dbFactory database.DBFactory,
	quoteRepo repository.QuoteRepository,
	dateQuote repository.DateQuoteRepository) DateQuoteService {
	return &dateQuoteSvc{
		dbFactory: dbFactory,
		quoteRepo: quoteRepo, dateQuoteRepo: dateQuote,
	}
}

func (s *dateQuoteSvc) GetRandomQuoteInsertIfNotExist(ctx context.Context) (dto.GetQuoteResponse, error) {
	var (
		tx       = s.dbFactory.GetDB()
		function = "GetRandomQuoteInsertIfNotExist"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(purchaseServiceLoggingFmt, function))
	)

	dateQuote, err := s.dateQuoteRepo.Get(ctx, tx, datatypes.Date(time.Now().UTC()))
	if err == nil {
		return dto.GetQuoteResponse{
			Date:      uint(time.Time(dateQuote.Date).Unix()),
			Quote:     dateQuote.Quote.Quote,
			Author:    dateQuote.Quote.Author,
			LikeCount: dateQuote.LikeCount,
		}, nil
	}
	if database.IsNotFound(err) {
		logger.Info("creating quote of the day")
		quote, err := s.quoteRepo.GetRandomQuote(ctx, tx)
		if err != nil {
			if database.IsNotFound(err) {
				logger.Info("not found any quote")
				return dto.GetQuoteResponse{}, exception.Wrap(exception.ErrNoQuoteFound, err, "create quote of the day failed")
			}
			logger.Error("cannot create quote, error:", err)
			return dto.GetQuoteResponse{}, exception.Wrap(exception.ErrInternalServer, err, "create quote of the day failed")
		}
		dateQuote = model.DateQuote{
			Date:      datatypes.Date(time.Now().UTC().Truncate(time.Hour * 24)),
			QuoteID:   quote.ID,
			Quote:     quote,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		err = s.dateQuoteRepo.Create(ctx, tx, &dateQuote)
		if err != nil {
			logger.Error("cannot create quote, error:", err)
			return dto.GetQuoteResponse{}, exception.Wrap(exception.ErrInternalServer, err, "create quote of the day failed")
		}
		logger.Info("init quote of the day successfully")
		return dto.GetQuoteResponse{
			Date:      uint(time.Time(dateQuote.Date).Unix()),
			Quote:     dateQuote.Quote.Quote,
			Author:    dateQuote.Quote.Author,
			LikeCount: dateQuote.LikeCount,
		}, nil
	}
	logger.Error("cannot get quote of the day, error:", err)
	return dto.GetQuoteResponse{}, exception.Wrap(exception.ErrInternalServer, err, "get quote of the day failed")
}

func (s *dateQuoteSvc) Like(ctx context.Context, request dto.LikeRequest) (dto.LikeResponse, error) {
	var (
		tx       = s.dbFactory.GetDB()
		function = "Like"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(purchaseServiceLoggingFmt, function))
	)
	dateTime := time.Unix(int64(request.Date), 0)
	err := s.dateQuoteRepo.IncreaseLikeByOne(ctx, tx, datatypes.Date(dateTime))
	if err != nil {
		logger.Error("cannot like, req:", request, "error:", err)
		return dto.LikeResponse{}, err
	}
	dateQuote, err := s.dateQuoteRepo.Get(ctx, tx, datatypes.Date(dateTime))
	if err != nil {
		logger.Error("cannot get like count, req:", request, "error:", err)
		return dto.LikeResponse{}, err
	}
	return dto.LikeResponse{
		LikeCount: dateQuote.LikeCount,
	}, nil
}

func (s *dateQuoteSvc) Dislike(ctx context.Context, request dto.LikeRequest) (dto.LikeResponse, error) {
	var (
		tx       = s.dbFactory.GetDB()
		function = "Dislike"
		logger   = logging.FromContext(ctx).Named(fmt.Sprintf(purchaseServiceLoggingFmt, function))
	)
	dateTime := time.Unix(int64(request.Date), 0)
	err := s.dateQuoteRepo.DecreaseLikeByOne(ctx, tx, datatypes.Date(dateTime))
	if err != nil {
		logger.Error("cannot like, req:", request, "error:", err)
		return dto.LikeResponse{}, err
	}
	dateQuote, err := s.dateQuoteRepo.Get(ctx, tx, datatypes.Date(dateTime))
	if err != nil {
		logger.Error("cannot get like count, req:", request, "error:", err)
		return dto.LikeResponse{}, err
	}
	return dto.LikeResponse{
		LikeCount: dateQuote.LikeCount,
	}, nil
}
