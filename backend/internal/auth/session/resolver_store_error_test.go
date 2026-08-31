package session

import (
	"context"
	"errors"
	"testing"
)

func TestResolverResolveReturnsStoreError(
	t *testing.T,
) {
	t.Parallel()

	cookieReader, store :=
		newResolverTestDependencies(t)

	store.getErr =
		errResolverGetSession

	resolver, err := newResolver(
		cookieReader,
		store,
		resolverTestNow,
	)
	if err != nil {
		t.Fatalf(
			"newResolver() error = %v",
			err,
		)
	}

	ctx := context.Background()

	sessionValue, err := resolver.Resolve(
		ctx,
		newResolverTestRequest(),
	)
	if err == nil {
		t.Fatal(
			"Resolve() error = nil, want an error",
		)
	}

	if sessionValue != (Session{}) {
		t.Errorf(
			"Resolve() session = %+v, want zero value",
			sessionValue,
		)
	}

	if !errors.Is(
		err,
		errResolverGetSession,
	) {
		t.Errorf(
			"Resolve() error = %v, want wrapped %v",
			err,
			errResolverGetSession,
		)
	}

	if cookieReader.readCalls != 1 {
		t.Errorf(
			"CookieReader.Read() calls = %d, want 1",
			cookieReader.readCalls,
		)
	}

	if store.getCalls != 1 {
		t.Errorf(
			"Store.Get() calls = %d, want 1",
			store.getCalls,
		)
	}

	if store.getCtx != ctx {
		t.Errorf(
			"Store.Get() context = %v, want %v",
			store.getCtx,
			ctx,
		)
	}

	wantIDHash := resolverTestSessionIDHash(t)

	if store.idHash != wantIDHash {
		t.Errorf(
			"Store.Get() ID hash = %q, want %q",
			store.idHash,
			wantIDHash,
		)
	}
}
