package session

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestResolverResolveRejectsExpiredSession(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name      string
		expiresAt time.Time
	}{
		{
			name:      "expiration equals current time",
			expiresAt: resolverTestNow(),
		},
		{
			name: "expiration is before current time",
			expiresAt: resolverTestNow().Add(
				-time.Second,
			),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				cookieReader, store :=
					newResolverTestDependencies(t)

				store.session.ExpiresAt =
					test.expiresAt

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

				sessionValue, err :=
					resolver.Resolve(
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

				const expectedError = "session is expired"

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

				if store.getCalls != 1 {
					t.Errorf(
						"Store.Get() calls = %d, want 1",
						store.getCalls,
					)
				}

				wantIDHash :=
					resolverTestSessionIDHash(t)

				if store.idHash != wantIDHash {
					t.Errorf(
						"Store.Get() ID hash = %q, want %q",
						store.idHash,
						wantIDHash,
					)
				}

				if store.saveCalls != 0 {
					t.Errorf(
						"Store.Save() calls = %d, want 0",
						store.saveCalls,
					)
				}

				if store.deleteCalls != 0 {
					t.Errorf(
						"Store.Delete() calls = %d, want 0",
						store.deleteCalls,
					)
				}
			},
		)
	}
}
