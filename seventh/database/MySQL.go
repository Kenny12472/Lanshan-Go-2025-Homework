package database

import (
	"fmt"
	"seventh/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	DB = db

	if err := DB.AutoMigrate(&model.Todo{}); err != nil {
		return fmt.Errorf("自动建表失败: %v", err)
	}

	fmt.Println("✅ MySQL 连接成功")
	return nil
}
