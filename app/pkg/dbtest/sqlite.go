package dbtest

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Setup создает подключение к SQLite в памяти и проводит автомиграцию предоставленных моделей.
func Setup(t testing.TB, models ...interface{}) *gorm.DB {
	t.Helper()

	// Используем cache=shared для корректной работы in-memory SQLite с GORM
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	if len(models) > 0 {
		err = db.AutoMigrate(models...)
		if err != nil {
			t.Fatalf("failed to migrate test models: %v", err)
		}
	}

	return db
}
