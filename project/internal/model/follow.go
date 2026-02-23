package model

import "time"

type Follow struct {
	ID          uint64 `gorm:"primaryKey"`
	FollowerID  uint64 `gorm:"index;not null"` // 谁关注
	FollowingID uint64 `gorm:"index;not null"` // 被关注者
	CreatedAt   time.Time
}

// 建议在数据库上建立联合唯一索引(follower_id, following_id)以避免重复关注
