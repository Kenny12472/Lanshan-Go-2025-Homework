package model

import "time"

type Like struct {
	ID         uint64 `gorm:"primaryKey"`
	TargetType string `gorm:"size:20;index;not null"` // "article" or "comment"
	TargetID   uint64 `gorm:"index;not null"`
	UserID     uint64 `gorm:"index;not null;uniqueIndex:idx_like_user_target"`
	CreatedAt  time.Time
}
