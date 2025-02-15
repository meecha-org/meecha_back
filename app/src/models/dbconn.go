package models

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbconn *gorm.DB = nil
)

func Init() {
	// データベースを開く
	db, err := gorm.Open(sqlite.Open(os.Getenv("DBPATH")), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}


	// マイグレーション
	db.AutoMigrate(&FriendRequest{})
	db.AutoMigrate(&Friend{})

	// グローバル変数に格納
	dbconn = db
}