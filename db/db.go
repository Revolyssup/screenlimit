package db

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	mx sync.Mutex
}

func NewDB() (*DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Roles{})
	db.AutoMigrate(&Events{})
	return &DB{db, sync.Mutex{}}, nil
}
