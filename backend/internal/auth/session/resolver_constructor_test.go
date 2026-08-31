package session

import (
	"strings"
	"testing"
	"time"
)

func TestNewResolver(
	t *testing.T,
) {
	t.Parallel()

	cookieReader, store :=
		newResolverTestDependencies(t)

	resolver, err := NewResolver(
		cookieReader,
		store,
	)
	if err != nil {
		t.Fatalf(
			"NewResolver() error = %v",
			err,
		)
	}

	if resolver == nil {
		t.Fatal(
			"NewResolver() resolver = nil, want non-nil",
		)
	}
}

func TestNewResolverRejectsNilDependency(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name          string
		cookieReader  CookieReader
		store         Store
		expectedError string
	}{
		{
			name:          "nil cookie reader",
			cookieReader:  nil,
			store:         &fakeResolverSessionStore{},
			expectedError: "cookie reader is nil",
		},
		{
			name:          "nil session store",
			cookieReader:  &fakeResolverCookieReader{},
			store:         nil,
			expectedError: "session store is nil",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				resolver, err := NewResolver(
					test.cookieReader,
					test.store,
				)
				if err == nil {
					t.Fatal(
						"NewResolver() error = nil, want an error",
					)
				}

				if resolver != nil {
					t.Errorf(
						"NewResolver() resolver = %v, want nil",
						resolver,
					)
				}

				if !strings.Contains(
					err.Error(),
					test.expectedError,
				) {
					t.Errorf(
						"NewResolver() error = %q, want to contain %q",
						err,
						test.expectedError,
					)
				}
			},
		)
	}
}

func TestNewResolverRejectsNilClock(
	t *testing.T,
) {
	t.Parallel()

	cookieReader, store :=
		newResolverTestDependencies(t)

	resolver, err := newResolver(
		cookieReader,
		store,
		nil,
	)
	if err == nil {
		t.Fatal(
			"newResolver() error = nil, want an error",
		)
	}

	if resolver != nil {
		t.Errorf(
			"newResolver() resolver = %v, want nil",
			resolver,
		)
	}

	const expectedError = "clock is nil"

	if !strings.Contains(
		err.Error(),
		expectedError,
	) {
		t.Errorf(
			"newResolver() error = %q, want to contain %q",
			err,
			expectedError,
		)
	}
}

func TestNewResolverSetsClock(
	t *testing.T,
) {
	t.Parallel()

	cookieReader, store :=
		newResolverTestDependencies(t)

	clock := func() time.Time {
		return resolverTestNow()
	}

	resolver, err := newResolver(
		cookieReader,
		store,
		clock,
	)
	if err != nil {
		t.Fatalf(
			"newResolver() error = %v",
			err,
		)
	}

	if resolver == nil {
		t.Fatal(
			"newResolver() resolver = nil, want non-nil",
		)
	}

	if resolver.now == nil {
		t.Fatal(
			"newResolver() clock = nil, want non-nil",
		)
	}

	if got := resolver.now(); !got.Equal(
		resolverTestNow(),
	) {
		t.Errorf(
			"resolver clock = %v, want %v",
			got,
			resolverTestNow(),
		)
	}
}
