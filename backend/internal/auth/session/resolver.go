package session

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CookieReader interface {
	Read(
		request *http.Request,
	) (
		string,
		error,
	)
}

type Resolver struct {
	cookieReader CookieReader
	store        Store
	now          func() time.Time
}

func NewResolver(
	cookieReader CookieReader,
	store Store,
) (
	*Resolver,
	error,
) {
	return newResolver(
		cookieReader,
		store,
		time.Now,
	)
}

func newResolver(
	cookieReader CookieReader,
	store Store,
	now func() time.Time,
) (
	*Resolver,
	error,
) {
	if cookieReader == nil {
		return nil, fmt.Errorf(
			"create session resolver: cookie reader is nil",
		)
	}

	if store == nil {
		return nil, fmt.Errorf(
			"create session resolver: session store is nil",
		)
	}

	if now == nil {
		return nil, fmt.Errorf(
			"create session resolver: clock is nil",
		)
	}

	return &Resolver{
		cookieReader: cookieReader,
		store:        store,
		now:          now,
	}, nil
}

func (
	resolver *Resolver,
) Resolve(
	ctx context.Context,
	request *http.Request,
) (
	Session,
	error,
) {
	if resolver == nil {
		return Session{}, fmt.Errorf(
			"resolve session: resolver is nil",
		)
	}

	if ctx == nil {
		return Session{}, fmt.Errorf(
			"resolve session: context is nil",
		)
	}

	if request == nil {
		return Session{}, fmt.Errorf(
			"resolve session: request is nil",
		)
	}

	if resolver.cookieReader == nil {
		return Session{}, fmt.Errorf(
			"resolve session: cookie reader is nil",
		)
	}

	if resolver.store == nil {
		return Session{}, fmt.Errorf(
			"resolve session: session store is nil",
		)
	}

	if resolver.now == nil {
		return Session{}, fmt.Errorf(
			"resolve session: clock is nil",
		)
	}

	rawSessionID, err :=
		resolver.cookieReader.Read(request)
	if err != nil {
		return Session{}, fmt.Errorf(
			"resolve session: read cookie: %w",
			err,
		)
	}

	if strings.TrimSpace(rawSessionID) == "" {
		return Session{}, fmt.Errorf(
			"resolve session: session ID is empty",
		)
	}

	idHash, err := HashID(rawSessionID)
	if err != nil {
		return Session{}, fmt.Errorf(
			"resolve session: hash session ID: %w",
			err,
		)
	}

	sessionValue, err := resolver.store.Get(
		ctx,
		idHash,
	)
	if err != nil {
		return Session{}, fmt.Errorf(
			"resolve session: get session: %w",
			err,
		)
	}

	if !resolver.now().Before(
		sessionValue.ExpiresAt,
	) {
		return Session{}, fmt.Errorf(
			"resolve session: session is expired",
		)
	}

	return sessionValue, nil
}
