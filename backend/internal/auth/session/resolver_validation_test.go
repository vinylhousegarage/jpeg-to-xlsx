package session

import (
	"context"
	"strings"
	"testing"
)

func TestResolverResolveRejectsNilResolver(
	t *testing.T,
) {
	t.Parallel()

	var resolver *Resolver

	sessionValue, err := resolver.Resolve(
		context.Background(),
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

	const expectedError = "resolver is nil"

	if !strings.Contains(
		err.Error(),
		expectedError,
	) {
		t.Errorf(
			"Resolve() error = %q, want to contain %q",
			err,
			expectedError,
		)
	}
}

func TestResolverResolveRejectsNilContext(
	t *testing.T,
) {
	t.Parallel()

	cookieReader, store :=
		newResolverTestDependencies(t)

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

	var ctx context.Context

	sessionValue, err := resolver.Resolve(
		ctx, //nolint:staticcheck // Verify defensive validation of a nil context.
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

	const expectedError = "context is nil"

	if !strings.Contains(
		err.Error(),
		expectedError,
	) {
		t.Errorf(
			"Resolve() error = %q, want to contain %q",
			err,
			expectedError,
		)
	}

	if cookieReader.readCalls != 0 {
		t.Errorf(
			"CookieReader.Read() calls = %d, want 0",
			cookieReader.readCalls,
		)
	}

	if store.getCalls != 0 {
		t.Errorf(
			"Store.Get() calls = %d, want 0",
			store.getCalls,
		)
	}
}

func TestResolverResolveRejectsNilRequest(
	t *testing.T,
) {
	t.Parallel()

	cookieReader, store :=
		newResolverTestDependencies(t)

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

	sessionValue, err := resolver.Resolve(
		context.Background(),
		nil,
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

	const expectedError = "request is nil"

	if !strings.Contains(
		err.Error(),
		expectedError,
	) {
		t.Errorf(
			"Resolve() error = %q, want to contain %q",
			err,
			expectedError,
		)
	}

	if cookieReader.readCalls != 0 {
		t.Errorf(
			"CookieReader.Read() calls = %d, want 0",
			cookieReader.readCalls,
		)
	}

	if store.getCalls != 0 {
		t.Errorf(
			"Store.Get() calls = %d, want 0",
			store.getCalls,
		)
	}
}
