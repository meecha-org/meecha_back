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
	// db, err := gorm.Open(sqlite.Open(os.Getenv("DBPATH")), &gorm.Config{})
	// dsn := "root:root@tcp(mysql:3306)/meecha?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// マイグレーション
	db.AutoMigrate(&FriendRequest{})
	db.AutoMigrate(&Friend{})

	// グローバル変数に格納
	dbconn = db
}
