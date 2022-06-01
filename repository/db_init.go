package repository

import (
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

var db *gorm.DB

func Init() error {
	var err error
	dsn := "root:mxyyyds107@tcp(127.0.0.1:3306)/tiktoktest?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err

}
