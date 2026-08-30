package oauthstate

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New(
	"OAuth state not found",
)

type Store interface {
	Save(
		ctx context.Context,
		state State,
	) error

	Consume(
		ctx context.Context,
		value string,
	) (
		State,
		error,
	)
}
