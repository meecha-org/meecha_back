package models

import (
	"log"
	"os"
	// "os"

	// "gorm.io/driver/sqlite"
	// "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbconn *gorm.DB = nil
)

func Init() {
	// データベースを開く
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// グローバル変数に格納
	dbconn = db
}
