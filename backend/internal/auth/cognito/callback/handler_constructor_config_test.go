package callback

import (
	"testing"
	"time"
)

func TestNewHandlerRejectsInvalidRedirectURL(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name        string
		redirectURL string
	}{
		{
			name:        "empty",
			redirectURL: "",
		},
		{
			name:        "whitespace",
			redirectURL: "   ",
		},
		{
			name:        "relative URL",
			redirectURL: "/",
		},
		{
			name:        "HTTP URL",
			redirectURL: "http://example.com/",
		},
		{
			name:        "missing host",
			redirectURL: "https:///callback",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				dependencies :=
					newCallbackTestDependencies()

				config :=
					validTestHandlerConfig()
				config.RedirectURL =
					test.redirectURL

				handler, err := NewHandler(
					dependencies.oauthStateStore,
					dependencies.cognitoClient,
					dependencies.verifier,
					dependencies.idGenerator,
					dependencies.sessionStore,
					dependencies.cookieWriter,
					config,
				)
				if err == nil {
					t.Fatal(
						"NewHandler() error = nil, want an error",
					)
				}

				if handler != nil {
					t.Errorf(
						"NewHandler() handler = %v, want nil",
						handler,
					)
				}
			},
		)
	}
}

func TestNewHandlerRejectsInvalidSessionLifetime(
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
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				dependencies :=
					newCallbackTestDependencies()

				config :=
					validTestHandlerConfig()
				config.SessionLifetime =
					test.lifetime

				handler, err := NewHandler(
					dependencies.oauthStateStore,
					dependencies.cognitoClient,
					dependencies.verifier,
					dependencies.idGenerator,
					dependencies.sessionStore,
					dependencies.cookieWriter,
					config,
				)
				if err == nil {
					t.Fatal(
						"NewHandler() error = nil, want an error",
					)
				}

				if handler != nil {
					t.Errorf(
						"NewHandler() handler = %v, want nil",
						handler,
					)
				}
			},
		)
	}
}
