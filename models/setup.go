package models

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DeletedAt sql.NullTime

type Model struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		fmt.Printf("Failed opening db\nError:%s", err.Error());
		return
	}

	// sqlite performance tuning
	// https://phiresky.github.io/blog/2020/sqlite-performance-tuning/
	db.Exec("pragma journal_mode = WAL;")
	db.Exec("pragma synchronous = normal;")
	db.Exec("pragma temp_store = memory;")
	db.Exec("pragma mmap_size = 30000000000;")

	db.AutoMigrate(&ErLog{})

	DB = db
}

// todo: make this actually check if the db is connected
func IsConnected() bool {
	if DB == nil {
		return false
	}

	return true
}