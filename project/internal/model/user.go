package model

import "time"

type User struct {
	ID           uint64 `gorm:"primaryKey"`
	Username     string `gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"` // 存 bcrypt 哈希
	DisplayName  string `gorm:"size:128;not null;default:''"`
	Bio          string `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
