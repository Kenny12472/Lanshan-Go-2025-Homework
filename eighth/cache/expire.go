package cache

import (
	"math/rand"
	"time"
)

func RandomExpire() time.Duration {
	return time.Minute*5 + time.Duration(rand.Intn(60))*time.Second
}
