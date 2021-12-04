package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/duyquang6/quote-today/internal/database"
	_mockDB "github.com/duyquang6/quote-today/internal/database/mocks"
	"github.com/duyquang6/quote-today/internal/model"
	"github.com/duyquang6/quote-today/internal/repository"
	_mockRepo "github.com/duyquang6/quote-today/internal/repository/mocks"
	"github.com/duyquang6/quote-today/pkg/dto"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestNewDateQuoteService(t *testing.T) {
	t.Parallel()
	type args struct {
		dbFactory     database.DBFactory
		dateQuoteRepo repository.DateQuoteRepository
		quoteRepo     repository.QuoteRepository
	}
	dateQuoteRepoMock := &_mockRepo.DateQuoteRepository{}
	quoteRepoMock := &_mockRepo.QuoteRepository{}
	dbFactoryMock := &_mockDB.DBFactory{}
	tests := []struct {
		name string
		args args
		want DateQuoteService
	}{
		{
			name: "TC1_NewPurchaseServiceSuccess",
			args: args{
				dbFactory:     dbFactoryMock,
				dateQuoteRepo: dateQuoteRepoMock,
				quoteRepo:     quoteRepoMock,
			},
			want: &dateQuoteSvc{dbFactoryMock, quoteRepoMock, dateQuoteRepoMock},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDateQuoteService(tt.args.dbFactory, tt.args.quoteRepo, tt.args.dateQuoteRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDateQuoteService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dateQuoteSvc_GetRandomQuoteInsertIfNotExist(t *testing.T) {
	t.Parallel()

	type fields struct {
		db            database.DBFactory
		dateQuoteRepo repository.DateQuoteRepository
		quoteRepo     repository.QuoteRepository
	}

	type args struct {
		ctx context.Context
	}

	timeNow := time.Now().UTC()

	mockQuote := model.Quote{
		ID:        1,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Author:    "Hello",
		Quote:     "Hole",
	}

	mockDateQuote := model.DateQuote{
		Date:      datatypes.Date(timeNow),
		QuoteID:   1,
		Quote:     mockQuote,
		LikeCount: 1,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	//db := &gorm.DB{}
	dbFactoryMock := &_mockDB.DBFactory{}
	dbFactoryMock.On("GetDB").Return(&gorm.DB{})

	dateQuoteRepoHappyCase := &_mockRepo.DateQuoteRepository{}
	dateQuoteRepoHappyCase.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(mockDateQuote, nil)

	dateQuoteRepoNotFound := &_mockRepo.DateQuoteRepository{}
	dateQuoteRepoNotFound.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(model.DateQuote{}, gorm.ErrRecordNotFound)
	dateQuoteRepoNotFound.On("Create", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	quoteRepoSuccess := &_mockRepo.QuoteRepository{}
	quoteRepoSuccess.On("GetRandomQuote", mock.Anything, mock.Anything).
		Return(mockQuote, nil)

	dateQuoteRepoUnexpectedErr := &_mockRepo.DateQuoteRepository{}
	dateQuoteRepoUnexpectedErr.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(model.DateQuote{}, errors.New("unexpected"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.GetQuoteResponse
		wantErr bool
	}{
		{
			name: "TC1_GetSuccess",
			fields: fields{
				db:            dbFactoryMock,
				dateQuoteRepo: dateQuoteRepoHappyCase,
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
			want: dto.GetQuoteResponse{
				Date:      uint(timeNow.Unix()),
				Quote:     "Hole",
				Author:    "Hello",
				LikeCount: 1,
			},
		},
		{
			name: "TC2_CreateSuccess",
			fields: fields{
				db:            dbFactoryMock,
				dateQuoteRepo: dateQuoteRepoNotFound,
				quoteRepo:     quoteRepoSuccess,
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
			want: dto.GetQuoteResponse{
				Date:      uint(time.Now().Truncate(time.Hour * 24).Unix()),
				Quote:     "Hole",
				Author:    "Hello",
				LikeCount: 0,
			},
		},
		{
			name: "TC3_GetDateQuoteDBFailed",
			fields: fields{
				db:            dbFactoryMock,
				dateQuoteRepo: dateQuoteRepoUnexpectedErr,
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dateQuoteSvc{
				dbFactory:     tt.fields.db,
				dateQuoteRepo: tt.fields.dateQuoteRepo,
				quoteRepo:     tt.fields.quoteRepo,
			}
			got, err := s.GetRandomQuoteInsertIfNotExist(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("dateQuoteSvc.CreatePurchase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dateQuoteSvc.CreatePurchase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dateQuoteSvc_Like(t *testing.T) {
	t.Parallel()

	type fields struct {
		db            database.DBFactory
		dateQuoteRepo repository.DateQuoteRepository
		quoteRepo     repository.QuoteRepository
	}

	type args struct {
		ctx context.Context
		req dto.LikeRequest
	}

	timeNow := time.Now().UTC()

	mockQuote := model.Quote{
		ID:        1,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Author:    "Hello",
		Quote:     "Hole",
	}

	mockDateQuote := model.DateQuote{
		Date:      datatypes.Date(timeNow),
		QuoteID:   1,
		Quote:     mockQuote,
		LikeCount: 1,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	dbFactoryMock := &_mockDB.DBFactory{}
	dbFactoryMock.On("GetDB").Return(&gorm.DB{})

	dateQuoteRepoHappyCase := &_mockRepo.DateQuoteRepository{}
	dateQuoteRepoHappyCase.On("IncreaseLikeByOne", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	dateQuoteRepoHappyCase.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(mockDateQuote, nil)

	dateQuoteRepoUnexpectedErr := &_mockRepo.DateQuoteRepository{}
	dateQuoteRepoUnexpectedErr.On("IncreaseLikeByOne", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("unexpected"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.LikeResponse
		wantErr bool
	}{
		{
			name: "TC1_LikeSuccess",
			fields: fields{
				db:            dbFactoryMock,
				dateQuoteRepo: dateQuoteRepoHappyCase,
			},
			args: args{
				ctx: context.TODO(),
				req: dto.LikeRequest{Date: uint(time.Now().UTC().Truncate(time.Hour * 24).Unix())},
			},
			wantErr: false,
			want: dto.LikeResponse{
				LikeCount: 1,
			},
		},
		{
			name: "TC2_LikeFailedDBError",
			fields: fields{
				db:            dbFactoryMock,
				dateQuoteRepo: dateQuoteRepoUnexpectedErr,
			},
			args: args{
				ctx: context.TODO(),
				req: dto.LikeRequest{Date: uint(time.Now().UTC().Truncate(time.Hour * 24).Unix())},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dateQuoteSvc{
				dbFactory:     tt.fields.db,
				dateQuoteRepo: tt.fields.dateQuoteRepo,
				quoteRepo:     tt.fields.quoteRepo,
			}
			got, err := s.Like(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("dateQuoteSvc.CreatePurchase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dateQuoteSvc.CreatePurchase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dateQuoteSvc_Dislike(t *testing.T) {
	t.Parallel()

	type fields struct {
		db            database.DBFactory
		dateQuoteRepo repository.DateQuoteRepository
		quoteRepo     repository.QuoteRepository
	}

	type args struct {
		ctx context.Context
		req dto.LikeRequest
	}

	timeNow := time.Now().UTC()

	mockQuote := model.Quote{
		ID:        1,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Author:    "Hello",
		Quote:     "Hole",
	}

	mockDateQuote := model.DateQuote{
		Date:      datatypes.Date(timeNow),
		QuoteID:   1,
		Quote:     mockQuote,
		LikeCount: 1,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	dbFactoryMock := &_mockDB.DBFactory{}
	dbFactoryMock.On("GetDB").Return(&gorm.DB{})

	dateQuoteRepoHappyCase := &_mockRepo.DateQuoteRepository{}
	dateQuoteRepoHappyCase.On("DecreaseLikeByOne", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	dateQuoteRepoHappyCase.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(mockDateQuote, nil)

	dateQuoteRepoUnexpectedErr := &_mockRepo.DateQuoteRepository{}
	dateQuoteRepoUnexpectedErr.On("DecreaseLikeByOne", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("unexpected"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.LikeResponse
		wantErr bool
	}{
		{
			name: "TC1_LikeSuccess",
			fields: fields{
				db:            dbFactoryMock,
				dateQuoteRepo: dateQuoteRepoHappyCase,
			},
			args: args{
				ctx: context.TODO(),
				req: dto.LikeRequest{Date: uint(time.Now().UTC().Truncate(time.Hour * 24).Unix())},
			},
			wantErr: false,
			want: dto.LikeResponse{
				LikeCount: 1,
			},
		},
		{
			name: "TC2_LikeFailedDBError",
			fields: fields{
				db:            dbFactoryMock,
				dateQuoteRepo: dateQuoteRepoUnexpectedErr,
			},
			args: args{
				ctx: context.TODO(),
				req: dto.LikeRequest{Date: uint(time.Now().UTC().Truncate(time.Hour * 24).Unix())},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &dateQuoteSvc{
				dbFactory:     tt.fields.db,
				dateQuoteRepo: tt.fields.dateQuoteRepo,
				quoteRepo:     tt.fields.quoteRepo,
			}
			got, err := s.Dislike(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("dateQuoteSvc.CreatePurchase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dateQuoteSvc.CreatePurchase() = %v, want %v", got, tt.want)
			}
		})
	}
}
