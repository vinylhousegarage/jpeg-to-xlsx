package session

import "context"

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
