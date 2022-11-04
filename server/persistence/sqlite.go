package persistence

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const IN_MEMORY_DB string = "file::memory:?cache=shared"

func newSqliteDB(dbPath string, migrationModels ...interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not open and connect to database at %s: %w", dbPath, err)
	}

	// auto-migrate provided models
	for _, aModel := range migrationModels {
		if err := db.AutoMigrate(aModel); err != nil {
			return nil, fmt.Errorf("could not migrate model %T: %w", aModel, err)
		}
	}

	// For sqlite, have to turn on foreign keys
	if res := db.Exec("PRAGMA foreign_keys = ON", nil); res.Error != nil {
		return nil, fmt.Errorf("unable to turn on foreign keys in sqlite db: %w", res.Error)
	}
	return db, nil
}
