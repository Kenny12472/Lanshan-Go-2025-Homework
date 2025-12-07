package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"type:varchar(100);unique;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
