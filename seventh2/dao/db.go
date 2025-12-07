package dao

import (
	models "seventh2/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {

	dsn := "root:123456@tcp(127.0.0.1:3306)/school?charset=utf8mb4&parseTime=True&loc=Local"

	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	return nil
}
