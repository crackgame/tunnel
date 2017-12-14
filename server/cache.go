package server

import (
	"sync"
)

var userMu sync.RWMutex
var users = map[int]*User{}

func AddUser(user *User) {
	userMu.Lock()
	defer userMu.Unlock()

	users[user.id] = user
}

func RemoveUser(userID int) {
	userMu.Lock()
	defer userMu.Unlock()

	delete(users, userID)
}

func FindUser(userID int) *User {
	userMu.Lock()
	defer userMu.Unlock()

	return users[userID]
}

func hasUser(userID int) bool {
	_, ok := users[userID]
	if !ok {
		return false
	}

	return true
}
