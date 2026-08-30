package cognito

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestVerifierVerifyReturnsJWKSEndpointError(
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

	server.setJWKSResponse(
		http.StatusInternalServerError,
		`{
			"error": "internal_server_error"
		}`,
	)

	claims := server.newIDTokenClaims(
		time.Now().UTC(),
	)

	rawIDToken := server.signIDToken(
		t,
		claims,
	)

	identity, err := verifier.Verify(
		context.Background(),
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

	if requests.JWKSCalls == 0 {
		t.Error(
			"JWKS calls = 0, want at least 1",
		)
	}
}

func TestVerifierVerifyRejectsMalformedJWKSResponse(
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

	server.setJWKSResponse(
		http.StatusOK,
		`{"keys":`,
	)

	claims := server.newIDTokenClaims(
		time.Now().UTC(),
	)

	rawIDToken := server.signIDToken(
		t,
		claims,
	)

	identity, err := verifier.Verify(
		context.Background(),
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

	if requests.JWKSCalls == 0 {
		t.Error(
			"JWKS calls = 0, want at least 1",
		)
	}
}
