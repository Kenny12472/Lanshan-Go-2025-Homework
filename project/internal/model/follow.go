package model

import "time"

type Follow struct {
	ID          uint64 `gorm:"primaryKey"`
	FollowerID  uint64 `gorm:"index;not null"`
	FollowingID uint64 `gorm:"index;not null"`
	CreatedAt   time.Time
}
