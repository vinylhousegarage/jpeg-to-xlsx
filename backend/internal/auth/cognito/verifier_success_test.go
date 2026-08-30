package cognito

import (
	"context"
	"testing"
	"time"
)

func TestVerifierVerify(
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
		context.Background(),
		rawIDToken,
		testNonce,
	)
	if err != nil {
		t.Fatalf(
			"Verify() error = %v",
			err,
		)
	}

	if identity.Subject != testSubject {
		t.Errorf(
			"Verify() subject = %q, want %q",
			identity.Subject,
			testSubject,
		)
	}

	if identity.Email != testEmail {
		t.Errorf(
			"Verify() email = %q, want %q",
			identity.Email,
			testEmail,
		)
	}

	requests := server.requests()

	if requests.DiscoveryCalls != 1 {
		t.Errorf(
			"OIDC discovery calls = %d, want 1",
			requests.DiscoveryCalls,
		)
	}

	if requests.JWKSCalls != 1 {
		t.Errorf(
			"JWKS calls = %d, want 1",
			requests.JWKSCalls,
		)
	}
}

func TestVerifierVerifyUsesCachedJWKS(
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

	now := time.Now().UTC()

	firstToken := server.signIDToken(
		t,
		server.newIDTokenClaims(now),
	)

	firstIdentity, err := verifier.Verify(
		context.Background(),
		firstToken,
		testNonce,
	)
	if err != nil {
		t.Fatalf(
			"first Verify() error = %v",
			err,
		)
	}

	if firstIdentity.Subject != testSubject {
		t.Errorf(
			"first Verify() subject = %q, want %q",
			firstIdentity.Subject,
			testSubject,
		)
	}

	secondClaims :=
		server.newIDTokenClaims(now)

	secondClaims["sub"] =
		"second-test-cognito-subject"
	secondClaims["email"] =
		"second-user@example.com"

	secondToken := server.signIDToken(
		t,
		secondClaims,
	)

	secondIdentity, err := verifier.Verify(
		context.Background(),
		secondToken,
		testNonce,
	)
	if err != nil {
		t.Fatalf(
			"second Verify() error = %v",
			err,
		)
	}

	if secondIdentity.Subject !=
		"second-test-cognito-subject" {
		t.Errorf(
			"second Verify() subject = %q, want %q",
			secondIdentity.Subject,
			"second-test-cognito-subject",
		)
	}

	if secondIdentity.Email !=
		"second-user@example.com" {
		t.Errorf(
			"second Verify() email = %q, want %q",
			secondIdentity.Email,
			"second-user@example.com",
		)
	}

	requests := server.requests()

	if requests.DiscoveryCalls != 1 {
		t.Errorf(
			"OIDC discovery calls = %d, want 1",
			requests.DiscoveryCalls,
		)
	}

	if requests.JWKSCalls != 1 {
		t.Errorf(
			"JWKS calls = %d, want 1 because the signing key should be cached",
			requests.JWKSCalls,
		)
	}
}
