package database

import (
	"context"
	"log"
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
			quotes := []model.Quote{
				{Author: "Thomas Edison", Quote: "Genius is one percent inspiration and ninety-nine percent perspiration."},
				{Author: "Yogi Berra", Quote: "You can observe a lot just by watching."},
				{Author: "Abraham Lincoln", Quote: "A house divided against itself cannot stand."},
				{Author: "Johann Wolfgang von Goethe", Quote: "Difficulties increase the nearer we get to the goal."},
				{Author: "Byron Pulsifer", Quote: "Fate is in your hands and no one elses"},
				{Author: "Thomas Edison", Quote: "Genius is one percent inspiration and ninety-nine percent perspiration."},
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
