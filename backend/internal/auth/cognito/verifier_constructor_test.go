package cognito

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestNewVerifier(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOIDCServer(t)

	verifier, err := NewVerifier(
		context.Background(),
		VerifierConfig{
			Issuer:   server.issuer(),
			ClientID: testClientID,
		},
	)
	if err != nil {
		t.Fatalf(
			"NewVerifier() error = %v",
			err,
		)
	}

	if verifier == nil {
		t.Fatal(
			"NewVerifier() verifier = nil, want non-nil",
		)
	}

	requests := server.requests()

	if requests.DiscoveryCalls != 1 {
		t.Errorf(
			"OIDC discovery calls = %d, want 1",
			requests.DiscoveryCalls,
		)
	}

	if requests.JWKSCalls != 0 {
		t.Errorf(
			"JWKS calls = %d, want 0 before token verification",
			requests.JWKSCalls,
		)
	}
}

func TestNewVerifierRejectsNilContext(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOIDCServer(t)

	verifier, err := NewVerifier(
		nil, //nolint:staticcheck // Intentionally verifies nil context rejection.
		VerifierConfig{
			Issuer:   server.issuer(),
			ClientID: testClientID,
		},
	)
	if err == nil {
		t.Fatal(
			"NewVerifier() error = nil, want an error",
		)
	}

	if verifier != nil {
		t.Errorf(
			"NewVerifier() verifier = %v, want nil",
			verifier,
		)
	}

	requests := server.requests()

	if requests.DiscoveryCalls != 0 {
		t.Errorf(
			"OIDC discovery calls = %d, want 0",
			requests.DiscoveryCalls,
		)
	}
}

func TestNewVerifierRejectsInvalidConfig(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*VerifierConfig)
	}{
		{
			name: "empty issuer",
			mutate: func(
				config *VerifierConfig,
			) {
				config.Issuer = ""
			},
		},
		{
			name: "whitespace issuer",
			mutate: func(
				config *VerifierConfig,
			) {
				config.Issuer = "   "
			},
		},
		{
			name: "invalid issuer",
			mutate: func(
				config *VerifierConfig,
			) {
				config.Issuer =
					"://invalid"
			},
		},
		{
			name: "issuer without scheme",
			mutate: func(
				config *VerifierConfig,
			) {
				config.Issuer =
					"example.com"
			},
		},
		{
			name: "issuer without host",
			mutate: func(
				config *VerifierConfig,
			) {
				config.Issuer =
					"https://"
			},
		},
		{
			name: "issuer with unsupported scheme",
			mutate: func(
				config *VerifierConfig,
			) {
				config.Issuer =
					"ftp://example.com"
			},
		},
		{
			name: "empty client ID",
			mutate: func(
				config *VerifierConfig,
			) {
				config.ClientID = ""
			},
		},
		{
			name: "whitespace client ID",
			mutate: func(
				config *VerifierConfig,
			) {
				config.ClientID = "   "
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				server := newTestOIDCServer(t)

				config := VerifierConfig{
					Issuer:   server.issuer(),
					ClientID: testClientID,
				}

				test.mutate(&config)

				verifier, err := NewVerifier(
					context.Background(),
					config,
				)
				if err == nil {
					t.Fatal(
						"NewVerifier() error = nil, want an error",
					)
				}

				if verifier != nil {
					t.Errorf(
						"NewVerifier() verifier = %v, want nil",
						verifier,
					)
				}

				requests := server.requests()

				if requests.DiscoveryCalls != 0 {
					t.Errorf(
						"OIDC discovery calls = %d, want 0",
						requests.DiscoveryCalls,
					)
				}
			},
		)
	}
}

func TestNewVerifierReturnsDiscoveryEndpointError(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOIDCServer(t)

	server.setDiscoveryResponse(
		http.StatusInternalServerError,
		`{
			"error": "internal_server_error"
		}`,
	)

	verifier, err := NewVerifier(
		context.Background(),
		VerifierConfig{
			Issuer:   server.issuer(),
			ClientID: testClientID,
		},
	)
	if err == nil {
		t.Fatal(
			"NewVerifier() error = nil, want an error",
		)
	}

	if verifier != nil {
		t.Errorf(
			"NewVerifier() verifier = %v, want nil",
			verifier,
		)
	}

	requests := server.requests()

	if requests.DiscoveryCalls != 1 {
		t.Errorf(
			"OIDC discovery calls = %d, want 1",
			requests.DiscoveryCalls,
		)
	}
}

func TestNewVerifierRejectsMalformedDiscoveryResponse(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOIDCServer(t)

	server.setDiscoveryResponse(
		http.StatusOK,
		`{"issuer":`,
	)

	verifier, err := NewVerifier(
		context.Background(),
		VerifierConfig{
			Issuer:   server.issuer(),
			ClientID: testClientID,
		},
	)
	if err == nil {
		t.Fatal(
			"NewVerifier() error = nil, want an error",
		)
	}

	if verifier != nil {
		t.Errorf(
			"NewVerifier() verifier = %v, want nil",
			verifier,
		)
	}
}

func TestNewVerifierRejectsDiscoveryIssuerMismatch(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOIDCServer(t)

	server.setDiscoveryResponse(
		http.StatusOK,
		fmt.Sprintf(
			`{
				"issuer": "https://unexpected.example.com",
				"authorization_endpoint": "%s/oauth2/authorize",
				"token_endpoint": "%s/oauth2/token",
				"jwks_uri": "%s",
				"response_types_supported": ["code"],
				"subject_types_supported": ["public"],
				"id_token_signing_alg_values_supported": ["RS256"]
			}`,
			server.issuer(),
			server.issuer(),
			server.jwksEndpoint(),
		),
	)

	verifier, err := NewVerifier(
		context.Background(),
		VerifierConfig{
			Issuer:   server.issuer(),
			ClientID: testClientID,
		},
	)
	if err == nil {
		t.Fatal(
			"NewVerifier() error = nil, want an error",
		)
	}

	if verifier != nil {
		t.Errorf(
			"NewVerifier() verifier = %v, want nil",
			verifier,
		)
	}
}
