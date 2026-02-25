package model

import "time"

type Comment struct {
	ID        uint64 `gorm:"primaryKey"`
	ArticleID uint64 `gorm:"index;not null"`
	UserID    uint64 `gorm:"index;not null"`
	Content   string `gorm:"type:text;not null"`
	LikeCount int64  `gorm:"default:0"`
	CreatedAt time.Time
}
