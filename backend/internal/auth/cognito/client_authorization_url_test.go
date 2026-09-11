package cognito

import (
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"strings"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
)

func TestClientAuthorizationURL(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOAuthServer(t)
	config := newTestClientConfig(server)

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	state := newTestOAuthState()

	authorizationURL, err := client.AuthorizationURL(state)
	if err != nil {
		t.Fatalf("AuthorizationURL() error = %v", err)
	}

	parsedURL, err := url.Parse(authorizationURL)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}

	expectedEndpoint, err := url.Parse(server.authorizationEndpoint())
	if err != nil {
		t.Fatalf("parse expected authorization endpoint: %v", err)
	}

	if parsedURL.Scheme != expectedEndpoint.Scheme {
		t.Errorf(
			"authorization URL scheme = %q, want %q",
			parsedURL.Scheme,
			expectedEndpoint.Scheme,
		)
	}

	if parsedURL.Host != expectedEndpoint.Host {
		t.Errorf(
			"authorization URL host = %q, want %q",
			parsedURL.Host,
			expectedEndpoint.Host,
		)
	}

	if parsedURL.Path != expectedEndpoint.Path {
		t.Errorf(
			"authorization URL path = %q, want %q",
			parsedURL.Path,
			expectedEndpoint.Path,
		)
	}

	query := parsedURL.Query()

	assertAuthorizationQueryValue(t, query, "response_type", "code")
	assertAuthorizationQueryValue(t, query, "client_id", testClientID)
	assertAuthorizationQueryValue(t, query, "redirect_uri", testRedirectURI)
	assertAuthorizationQueryValue(t, query, "scope", strings.Join(config.Scopes, " "))
	assertAuthorizationQueryValue(t, query, "state", state.Value)
	assertAuthorizationQueryValue(t, query, "nonce", state.Nonce)
	assertAuthorizationQueryValue(t, query, "identity_provider", "Google")
	assertAuthorizationQueryValue(t, query, "prompt", "select_account")
	assertAuthorizationQueryValue(t, query, "code_challenge", codeChallenge(state.CodeVerifier))
	assertAuthorizationQueryValue(t, query, "code_challenge_method", "S256")

	if query.Has("client_secret") {
		t.Errorf("authorization URL contains client_secret query parameter")
	}

	if strings.Contains(
		authorizationURL,
		testClientSecret,
	) {
		t.Errorf("authorization URL contains the client secret")
	}
}

func TestClientAuthorizationURLRejectsInvalidState(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*oauthstate.State)
	}{
		{
			name: "empty state",
			mutate: func(state *oauthstate.State) {
				state.Value = ""
			},
		},
		{
			name: "whitespace state",
			mutate: func(state *oauthstate.State) {
				state.Value = "   "
			},
		},
		{
			name: "empty code verifier",
			mutate: func(state *oauthstate.State) {
				state.CodeVerifier = ""
			},
		},
		{
			name: "whitespace code verifier",
			mutate: func(state *oauthstate.State) {
				state.CodeVerifier = "   "
			},
		},
		{
			name: "code verifier shorter than 43 characters",
			mutate: func(state *oauthstate.State) {
				state.CodeVerifier = strings.Repeat("v", 42)
			},
		},
		{
			name: "code verifier longer than 128 characters",
			mutate: func(state *oauthstate.State) {
				state.CodeVerifier = strings.Repeat("v", 129)
			},
		},
		{
			name: "code verifier contains invalid character",
			mutate: func(state *oauthstate.State) {
				state.CodeVerifier = strings.Repeat("v", 42) + "*"
			},
		},
		{
			name: "empty nonce",
			mutate: func(state *oauthstate.State) {
				state.Nonce = ""
			},
		},
		{
			name: "whitespace nonce",
			mutate: func(state *oauthstate.State) {
				state.Nonce = "   "
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				server := newTestOAuthServer(t)

				client, err := NewClient(
					newTestClientConfig(server),
				)
				if err != nil {
					t.Fatalf("NewClient() error = %v", err)
				}

				state := newTestOAuthState()
				test.mutate(&state)

				authorizationURL, err := client.AuthorizationURL(state)
				if err == nil {
					t.Fatal("AuthorizationURL() error = nil, want an error")
				}

				if authorizationURL != "" {
					t.Errorf("AuthorizationURL() URL = %q, want empty", authorizationURL)
				}
			},
		)
	}
}

func assertAuthorizationQueryValue(
	t *testing.T,
	query url.Values,
	name string,
	want string,
) {
	t.Helper()

	values, exists := query[name]
	if !exists {
		t.Errorf("authorization URL query %q is missing", name)

		return
	}

	if len(values) != 1 {
		t.Errorf(
			"authorization URL query %q has %d values, want 1",
			name,
			len(values),
		)

		return
	}

	if values[0] != want {
		t.Errorf(
			"authorization URL query %q = %q, want %q",
			name,
			values[0],
			want,
		)
	}
}

func codeChallenge(
	codeVerifier string,
) string {
	digest := sha256.Sum256([]byte(codeVerifier))

	return base64.RawURLEncoding.EncodeToString(digest[:])
}
