package session

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New(
	"session not found",
)

type Store interface {
	Save(
		ctx context.Context,
		session Session,
	) error

	Get(
		ctx context.Context,
		idHash string,
	) (Session, error)

	Delete(
		ctx context.Context,
		idHash string,
	) error
}
