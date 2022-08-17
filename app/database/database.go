package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var instant *gorm.DB

func New() *gorm.DB {

	DbUser := os.Getenv("MYSQL_USER")
	DbPassword := os.Getenv("MYSQL_PASSWORD")
	DbHost := os.Getenv("MYSQL_HOST")
	DbPort := os.Getenv("MYSQL_PORT")
	DbName := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func Instant() *gorm.DB {
	if instant == nil {
		instant = New()
	}
	return instant
}
