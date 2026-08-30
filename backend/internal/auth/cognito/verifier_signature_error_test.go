package cognito

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
)

func TestVerifierVerifyRejectsInvalidSignature(
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

	otherPrivateKey, err := rsa.GenerateKey(
		rand.Reader,
		2048,
	)
	if err != nil {
		t.Fatalf(
			"generate alternate RSA key: %v",
			err,
		)
	}

	claims := server.newIDTokenClaims(
		time.Now().UTC(),
	)

	rawIDToken := signTestIDToken(
		t,
		otherPrivateKey,
		testOIDCKeyID,
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

func TestVerifierVerifyRejectsUnknownKeyID(
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

	rawIDToken := signTestIDToken(
		t,
		server.privateKey,
		"unknown-key-id",
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
