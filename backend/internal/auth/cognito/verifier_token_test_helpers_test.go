package cognito

import (
	"crypto/rsa"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const (
	testSubject = "test-cognito-subject"
	testEmail   = "user@example.com"
)

func (
	server *testOIDCServer,
) newIDTokenClaims(
	now time.Time,
) map[string]any {
	return map[string]any{
		"iss":       server.issuer(),
		"sub":       testSubject,
		"aud":       testClientID,
		"exp":       now.Add(time.Hour).Unix(),
		"iat":       now.Add(-time.Minute).Unix(),
		"nonce":     testNonce,
		"token_use": "id",
		"email":     testEmail,
	}
}

func (
	server *testOIDCServer,
) signIDToken(
	t *testing.T,
	claims map[string]any,
) string {
	t.Helper()

	return signTestIDToken(
		t,
		server.privateKey,
		testOIDCKeyID,
		claims,
	)
}

func signTestIDToken(
	t *testing.T,
	privateKey *rsa.PrivateKey,
	keyID string,
	claims map[string]any,
) string {
	t.Helper()

	options := (&jose.SignerOptions{}).
		WithType("JWT").
		WithHeader(
			"kid",
			keyID,
		)

	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.RS256,
			Key:       privateKey,
		},
		options,
	)
	if err != nil {
		t.Fatalf(
			"create test JWT signer: %v",
			err,
		)
	}

	rawToken, err := jwt.Signed(
		signer,
	).Claims(
		claims,
	).Serialize()
	if err != nil {
		t.Fatalf(
			"sign test ID token: %v",
			err,
		)
	}

	return rawToken
}

func assertEmptyIdentity(
	t *testing.T,
	identity Identity,
) {
	t.Helper()

	if identity.Subject != "" {
		t.Errorf(
			"identity subject = %q, want empty",
			identity.Subject,
		)
	}

	if identity.Email != "" {
		t.Errorf(
			"identity email = %q, want empty",
			identity.Email,
		)
	}
}
