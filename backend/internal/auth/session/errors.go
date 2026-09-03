package session

import "errors"

var (
	ErrNotFound = errors.New(
		"session not found",
	)

	ErrUnauthenticated = errors.New(
		"session is unauthenticated",
	)
)
