package session

import (
	"testing"
	"time"
)

func TestNewCookieManager(
	t *testing.T,
) {
	t.Parallel()

	manager, err := NewCookieManager(
		testCookieLifetime,
	)
	if err != nil {
		t.Fatalf(
			"NewCookieManager() error = %v",
			err,
		)
	}

	if manager == nil {
		t.Fatal(
			"NewCookieManager() manager = nil, want non-nil",
		)
	}

	if manager.lifetime !=
		testCookieLifetime {
		t.Errorf(
			"CookieManager lifetime = %v, want %v",
			manager.lifetime,
			testCookieLifetime,
		)
	}

	if manager.now == nil {
		t.Error(
			"CookieManager now = nil, want non-nil",
		)
	}
}

func TestNewCookieManagerRejectsInvalidLifetime(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name     string
		lifetime time.Duration
	}{
		{
			name:     "zero",
			lifetime: 0,
		},
		{
			name:     "negative",
			lifetime: -time.Hour,
		},
		{
			name:     "less than one second",
			lifetime: time.Millisecond,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				manager, err :=
					NewCookieManager(
						test.lifetime,
					)
				if err == nil {
					t.Fatal(
						"NewCookieManager() error = nil, want an error",
					)
				}

				if manager != nil {
					t.Errorf(
						"NewCookieManager() manager = %v, want nil",
						manager,
					)
				}
			},
		)
	}
}
