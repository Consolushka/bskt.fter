package database

import (
	"IMP/app/internal/infra/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetDB(t *testing.T) {
	// Сбрасываем глобальное состояние перед тестом
	oldDb := db
	db = nil
	defer func() { db = oldDb }()

	t.Run("should return nil if db is not initiated", func(t *testing.T) {
		res := GetDB()
		assert.Nil(t, res)
	})

	t.Run("should return db instance if initiated", func(t *testing.T) {
		mockDB := &gorm.DB{}
		db = mockDB
		res := GetDB()
		assert.Equal(t, mockDB, res)
	})
}

func TestOpenDbConnection_Singleton(t *testing.T) {
	// Сбрасываем глобальное состояние
	oldDb := db
	mockDB := &gorm.DB{}
	db = mockDB
	defer func() { db = oldDb }()

	// Если db уже не nil, функция должна вернуть её сразу, не пытаясь подключиться
	res := OpenDbConnection(config.DatabaseConfig{})
	assert.Equal(t, mockDB, res)
}

func TestOpenDbConnection_Panic(t *testing.T) {
	// Сбрасываем глобальное состояние
	oldDb := db
	db = nil
	defer func() { db = oldDb }()

	assert.Panics(t, func() {
		OpenDbConnection(config.DatabaseConfig{
			Host: "",
			User: "",
		})
	})
}
