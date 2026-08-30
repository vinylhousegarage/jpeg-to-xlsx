package cognito

import (
	"context"
	"testing"
	"time"
)

func TestVerifierVerifyRejectsInvalidClaims(
	t *testing.T,
) {
	t.Parallel()

	now := time.Now().UTC()

	tests := []struct {
		name          string
		mutate        func(map[string]any)
		expectedNonce string
	}{
		{
			name: "missing issuer",
			mutate: func(
				claims map[string]any,
			) {
				delete(claims, "iss")
			},
			expectedNonce: testNonce,
		},
		{
			name: "issuer mismatch",
			mutate: func(
				claims map[string]any,
			) {
				claims["iss"] =
					"https://unexpected.example.com"
			},
			expectedNonce: testNonce,
		},
		{
			name: "missing audience",
			mutate: func(
				claims map[string]any,
			) {
				delete(claims, "aud")
			},
			expectedNonce: testNonce,
		},
		{
			name: "audience mismatch",
			mutate: func(
				claims map[string]any,
			) {
				claims["aud"] =
					"unexpected-client-id"
			},
			expectedNonce: testNonce,
		},
		{
			name: "missing expiration",
			mutate: func(
				claims map[string]any,
			) {
				delete(claims, "exp")
			},
			expectedNonce: testNonce,
		},
		{
			name: "expired token",
			mutate: func(
				claims map[string]any,
			) {
				claims["exp"] =
					now.Add(-time.Hour).Unix()
			},
			expectedNonce: testNonce,
		},
		{
			name: "missing subject",
			mutate: func(
				claims map[string]any,
			) {
				delete(claims, "sub")
			},
			expectedNonce: testNonce,
		},
		{
			name: "empty subject",
			mutate: func(
				claims map[string]any,
			) {
				claims["sub"] = ""
			},
			expectedNonce: testNonce,
		},
		{
			name: "whitespace subject",
			mutate: func(
				claims map[string]any,
			) {
				claims["sub"] = "   "
			},
			expectedNonce: testNonce,
		},
		{
			name: "missing nonce",
			mutate: func(
				claims map[string]any,
			) {
				delete(claims, "nonce")
			},
			expectedNonce: testNonce,
		},
		{
			name: "empty nonce",
			mutate: func(
				claims map[string]any,
			) {
				claims["nonce"] = ""
			},
			expectedNonce: testNonce,
		},
		{
			name: "whitespace nonce",
			mutate: func(
				claims map[string]any,
			) {
				claims["nonce"] = "   "
			},
			expectedNonce: testNonce,
		},
		{
			name: "nonce mismatch",
			mutate: func(
				_ map[string]any,
			) {
			},
			expectedNonce: "unexpected-nonce",
		},
		{
			name: "missing token use",
			mutate: func(
				claims map[string]any,
			) {
				delete(
					claims,
					"token_use",
				)
			},
			expectedNonce: testNonce,
		},
		{
			name: "empty token use",
			mutate: func(
				claims map[string]any,
			) {
				claims["token_use"] = ""
			},
			expectedNonce: testNonce,
		},
		{
			name: "access token use",
			mutate: func(
				claims map[string]any,
			) {
				claims["token_use"] =
					"access"
			},
			expectedNonce: testNonce,
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

				claims :=
					server.newIDTokenClaims(
						now,
					)

				test.mutate(claims)

				rawIDToken :=
					server.signIDToken(
						t,
						claims,
					)

				identity, err :=
					verifier.Verify(
						context.Background(),
						rawIDToken,
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
			},
		)
	}
}
