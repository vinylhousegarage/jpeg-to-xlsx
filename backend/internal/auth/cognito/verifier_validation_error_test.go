package cognito

import (
	"context"
	"testing"
	"time"
)

func TestVerifierVerifyRejectsInvalidInput(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name          string
		rawIDToken    string
		expectedNonce string
	}{
		{
			name:          "empty ID token",
			rawIDToken:    "",
			expectedNonce: testNonce,
		},
		{
			name:          "whitespace ID token",
			rawIDToken:    "   ",
			expectedNonce: testNonce,
		},
		{
			name:          "empty expected nonce",
			rawIDToken:    testIDToken,
			expectedNonce: "",
		},
		{
			name:          "whitespace expected nonce",
			rawIDToken:    testIDToken,
			expectedNonce: "   ",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
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

				identity, err := verifier.Verify(
					context.Background(),
					test.rawIDToken,
					test.expectedNonce,
				)
				if err == nil {
					t.Fatal(
						"Verify() error = nil, want an error",
					)
				}

				assertEmptyIdentity(
					t,
					identity,
				)

				requests := server.requests()

				if requests.JWKSCalls != 0 {
					t.Errorf(
						"JWKS calls = %d, want 0",
						requests.JWKSCalls,
					)
				}
			},
		)
	}
}

func TestVerifierVerifyRejectsNilContext(
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

	claims := server.newIDTokenClaims(
		time.Now().UTC(),
	)

	rawIDToken := server.signIDToken(
		t,
		claims,
	)

	identity, err := verifier.Verify(
		nil, //nolint:staticcheck // Intentionally verifies nil context rejection.
		rawIDToken,
		testNonce,
	)
	if err == nil {
		t.Fatal(
			"Verify() error = nil, want an error",
		)
	}

	assertEmptyIdentity(
		t,
		identity,
	)

	requests := server.requests()

	if requests.JWKSCalls != 0 {
		t.Errorf(
			"JWKS calls = %d, want 0",
			requests.JWKSCalls,
		)
	}
}

func TestVerifierVerifyRejectsMalformedIDToken(
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

	identity, err := verifier.Verify(
		context.Background(),
		"not-a-jwt",
		testNonce,
	)
	if err == nil {
		t.Fatal(
			"Verify() error = nil, want an error",
		)
	}

	assertEmptyIdentity(
		t,
		identity,
	)

	requests := server.requests()

	if requests.JWKSCalls != 0 {
		t.Errorf(
			"JWKS calls = %d, want 0",
			requests.JWKSCalls,
		)
	}
}
