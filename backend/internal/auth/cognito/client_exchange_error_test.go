package cognito

import (
	"context"
	"strings"
	"testing"
)

func TestClientExchangeRejectsInvalidInput(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name         string
		code         string
		codeVerifier string
	}{
		{
			name:         "empty authorization code",
			code:         "",
			codeVerifier: newTestOAuthState().CodeVerifier,
		},
		{
			name:         "whitespace authorization code",
			code:         "   ",
			codeVerifier: newTestOAuthState().CodeVerifier,
		},
		{
			name:         "empty code verifier",
			code:         testCode,
			codeVerifier: "",
		},
		{
			name:         "whitespace code verifier",
			code:         testCode,
			codeVerifier: "   ",
		},
		{
			name: "code verifier shorter than 43 characters",
			code: testCode,
			codeVerifier: strings.Repeat(
				"v",
				42,
			),
		},
		{
			name: "code verifier longer than 128 characters",
			code: testCode,
			codeVerifier: strings.Repeat(
				"v",
				129,
			),
		},
		{
			name: "code verifier contains invalid character",
			code: testCode,
			codeVerifier: strings.Repeat(
				"v",
				42,
			) + "*",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				server := newTestOAuthServer(t)

				client, err := NewClient(
					newTestClientConfig(
						server,
					),
				)
				if err != nil {
					t.Fatalf(
						"NewClient() error = %v",
						err,
					)
				}

				idToken, err := client.Exchange(
					context.Background(),
					test.code,
					test.codeVerifier,
				)
				if err == nil {
					t.Fatal(
						"Exchange() error = nil, want an error",
					)
				}

				if idToken != "" {
					t.Errorf(
						"Exchange() ID token = %q, want empty",
						idToken,
					)
				}

				request := server.tokenRequest()

				if request.Calls != 0 {
					t.Errorf(
						"token endpoint calls = %d, want 0",
						request.Calls,
					)
				}
			},
		)
	}
}

func TestClientExchangeReturnsTokenEndpointError(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOAuthServer(t)

	server.setTokenResponse(
		400,
		`{
			"error": "invalid_grant",
			"error_description": "Authorization code is invalid"
		}`,
	)

	client, err := NewClient(
		newTestClientConfig(server),
	)
	if err != nil {
		t.Fatalf(
			"NewClient() error = %v",
			err,
		)
	}

	state := newTestOAuthState()

	idToken, err := client.Exchange(
		context.Background(),
		testCode,
		state.CodeVerifier,
	)
	if err == nil {
		t.Fatal(
			"Exchange() error = nil, want an error",
		)
	}

	if idToken != "" {
		t.Errorf(
			"Exchange() ID token = %q, want empty",
			idToken,
		)
	}

	request := server.tokenRequest()

	if request.Calls != 1 {
		t.Errorf(
			"token endpoint calls = %d, want 1",
			request.Calls,
		)
	}
}

func TestClientExchangeRejectsMalformedTokenResponse(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOAuthServer(t)

	server.setTokenResponse(
		200,
		`{"access_token":`,
	)

	client, err := NewClient(
		newTestClientConfig(server),
	)
	if err != nil {
		t.Fatalf(
			"NewClient() error = %v",
			err,
		)
	}

	state := newTestOAuthState()

	idToken, err := client.Exchange(
		context.Background(),
		testCode,
		state.CodeVerifier,
	)
	if err == nil {
		t.Fatal(
			"Exchange() error = nil, want an error",
		)
	}

	if idToken != "" {
		t.Errorf(
			"Exchange() ID token = %q, want empty",
			idToken,
		)
	}
}

func TestClientExchangeRejectsMissingIDToken(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOAuthServer(t)

	server.setTokenResponse(
		200,
		`{
			"access_token": "test-access-token",
			"refresh_token": "test-refresh-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`,
	)

	client, err := NewClient(
		newTestClientConfig(server),
	)
	if err != nil {
		t.Fatalf(
			"NewClient() error = %v",
			err,
		)
	}

	state := newTestOAuthState()

	idToken, err := client.Exchange(
		context.Background(),
		testCode,
		state.CodeVerifier,
	)
	if err == nil {
		t.Fatal(
			"Exchange() error = nil, want an error",
		)
	}

	if idToken != "" {
		t.Errorf(
			"Exchange() ID token = %q, want empty",
			idToken,
		)
	}
}

func TestClientExchangeRejectsEmptyIDToken(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOAuthServer(t)

	server.setTokenResponse(
		200,
		`{
			"access_token": "test-access-token",
			"id_token": "   ",
			"refresh_token": "test-refresh-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`,
	)

	client, err := NewClient(
		newTestClientConfig(server),
	)
	if err != nil {
		t.Fatalf(
			"NewClient() error = %v",
			err,
		)
	}

	state := newTestOAuthState()

	idToken, err := client.Exchange(
		context.Background(),
		testCode,
		state.CodeVerifier,
	)
	if err == nil {
		t.Fatal(
			"Exchange() error = nil, want an error",
		)
	}

	if idToken != "" {
		t.Errorf(
			"Exchange() ID token = %q, want empty",
			idToken,
		)
	}
}

func TestClientExchangeRejectsNonStringIDToken(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOAuthServer(t)

	server.setTokenResponse(
		200,
		`{
			"access_token": "test-access-token",
			"id_token": 12345,
			"refresh_token": "test-refresh-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`,
	)

	client, err := NewClient(
		newTestClientConfig(server),
	)
	if err != nil {
		t.Fatalf(
			"NewClient() error = %v",
			err,
		)
	}

	state := newTestOAuthState()

	idToken, err := client.Exchange(
		context.Background(),
		testCode,
		state.CodeVerifier,
	)
	if err == nil {
		t.Fatal(
			"Exchange() error = nil, want an error",
		)
	}

	if idToken != "" {
		t.Errorf(
			"Exchange() ID token = %q, want empty",
			idToken,
		)
	}
}
