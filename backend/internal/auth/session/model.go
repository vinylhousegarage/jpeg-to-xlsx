package session

import "time"

type Session struct {
	IDHash     string
	CognitoSub string
	CreatedAt  time.Time
	ExpiresAt  time.Time
}
