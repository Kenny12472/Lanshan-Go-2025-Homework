package model

import "time"

const (
	ArticleDraft   = 0
	ArticlePublish = 1
	ArticleDeleted = 2
)

type Article struct {
	ID        uint64 `gorm:"primaryKey"`
	AuthorID  uint64 `gorm:"index;not null"`
	Title     string `gorm:"size:255;not null;index"`
	Content   string `gorm:"type:text;not null"`
	Status    int    `gorm:"not null;default:0"`
	ViewCount int64  `gorm:"default:0"`
	LikeCount int64  `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
