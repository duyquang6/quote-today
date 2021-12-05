package database

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/duyquang6/quote-today/internal/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormMigrationScripts = []*gormigrate.Migration{
	{
		ID: "00001-CreateQuoteTbl",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Quote{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&model.Quote{})
		},
	},
	{
		ID: "00002-CreateDateQuoteTbl",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.DateQuote{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&model.DateQuote{})
		},
	},
	{
		ID: "00003-InsertQuotes",
		Migrate: func(tx *gorm.DB) error {
			resp, err := http.Get("https://gist.githubusercontent.com/duyquang6/a0def025bdc27969b8cf0890a6b1bf86/raw/963b5a9355f04741239407320ac973a6096cd7b6/quotes.csv")
			if err != nil {
				return err
			}
			r := csv.NewReader(resp.Body)
			var quotes []model.Quote
			var author, quote string
			lineCount := 0
			for {
				record, err := r.Read()

				if err == io.EOF {
					break
				}

				if err != nil {
					log.Fatal(err)
				}
				lineCount++
				if lineCount == 1 {
					// skip line 1 because of header
					continue
				}
				author, quote = record[0], record[1]
				quotes = append(quotes, model.Quote{Author: author, Quote: quote})
			}
			return tx.Create(&quotes).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("TRUNCATE quotes RESTART IDENTITY CASCADE;").Error
		},
	},
}

// Migrate migrate schema
func (_db *DB) Migrate(ctx context.Context) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db := _db.GetDB().Session(&gorm.Session{Logger: newLogger})

	m := gormigrate.New(db, gormigrate.DefaultOptions, gormMigrationScripts)
	return m.Migrate()
}

// MigrateDown rollback schema
func (_db *DB) MigrateDown(ctx context.Context) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db := _db.GetDB().Session(&gorm.Session{Logger: newLogger})

	m := gormigrate.New(db, gormigrate.DefaultOptions, gormMigrationScripts)

	for i := range gormMigrationScripts {
		migration := gormMigrationScripts[len(gormMigrationScripts)-i-1]
		err := m.RollbackMigration(migration)
		if err != nil {
			newLogger.Error(ctx, "cannot rollback script", migration.ID)
			return err
		}
	}

	return nil
}
