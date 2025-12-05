package dao

import "sync"

// 简单的内存“数据库”
var (
	db  = map[string]string{}
	dbm sync.RWMutex
)

func AddUser(username, password string) {
	dbm.Lock()
	defer dbm.Unlock()
	db[username] = password
}

func FindUser(username, password string) bool {
	dbm.RLock()
	defer dbm.RUnlock()
	if p, ok := db[username]; ok {
		return p == password
	}
	return false
}

func IsUserExist(username string) bool {
	dbm.RLock()
	defer dbm.RUnlock()
	_, ok := db[username]
	return ok
}

func UpdatePassword(username, newPassword string) bool {
	dbm.Lock()
	defer dbm.Unlock()
	if _, ok := db[username]; !ok {
		return false
	}
	db[username] = newPassword
	return true
}
