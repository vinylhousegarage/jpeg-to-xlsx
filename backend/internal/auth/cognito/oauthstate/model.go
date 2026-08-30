package oauthstate

import "time"

type State struct {
	Value        string
	CodeVerifier string
	Nonce        string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}
