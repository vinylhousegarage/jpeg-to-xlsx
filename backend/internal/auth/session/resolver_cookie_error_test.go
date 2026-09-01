package session

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestResolverResolveReturnsCookieReaderError(
	t *testing.T,
) {
	t.Parallel()

	cookieReader, store :=
		newResolverTestDependencies(t)

	cookieReader.err =
		errResolverReadCookie

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

	request := newResolverTestRequest()

	sessionValue, err := resolver.Resolve(
		context.Background(),
		request,
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
		errResolverReadCookie,
	) {
		t.Errorf(
			"Resolve() error = %v, want wrapped %v",
			err,
			errResolverReadCookie,
		)
	}

	if cookieReader.readCalls != 1 {
		t.Errorf(
			"CookieReader.Read() calls = %d, want 1",
			cookieReader.readCalls,
		)
	}

	if cookieReader.request != request {
		t.Errorf(
			"CookieReader.Read() request = %p, want %p",
			cookieReader.request,
			request,
		)
	}

	if store.getCalls != 0 {
		t.Errorf(
			"Store.Get() calls = %d, want 0",
			store.getCalls,
		)
	}
}

func TestResolverResolveRejectsEmptySessionID(
	t *testing.T,
) {
	t.Parallel()

	cookieReader, store :=
		newResolverTestDependencies(t)

	cookieReader.rawSessionID = ""

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

	const expectedError = "session ID is empty"

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

	if cookieReader.readCalls != 1 {
		t.Errorf(
			"CookieReader.Read() calls = %d, want 1",
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
