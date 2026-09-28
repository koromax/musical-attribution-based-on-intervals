package ds

import "sync"

const DefaultUserID uint = 1

type UserSession struct {
	UserID uint
}

var (
	sessionInstance *UserSession
	once            sync.Once
)

func GetCurrentUserSingleton() *UserSession {
	once.Do(func() {
		sessionInstance = &UserSession{UserID: DefaultUserID}
	})
	return sessionInstance
}
