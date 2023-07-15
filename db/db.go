package db

import (
	"github.com/devson2561/eth-tracker/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDatabase(dbFile string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Transaction{}, &models.Address{}, &models.Block{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
