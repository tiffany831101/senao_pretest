package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(user, password, host, dbName string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, dbName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Successfully connect to database!")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get underlying DB:", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)
}

func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Fatal("Failed to get underlying DB:", err)
		}
		sqlDB.Close()
	}
}
