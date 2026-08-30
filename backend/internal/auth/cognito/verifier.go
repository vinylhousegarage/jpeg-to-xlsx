package cognito

import (
	"context"
	"crypto/subtle"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type VerifierConfig struct {
	Issuer   string
	ClientID string
}

type Identity struct {
	Subject string
	Email   string
}

type Verifier struct {
	idTokenVerifier *oidc.IDTokenVerifier
}

type idTokenClaims struct {
	Nonce    string `json:"nonce"`
	TokenUse string `json:"token_use"`
	Email    string `json:"email"`
}

func NewVerifier(
	ctx context.Context,
	config VerifierConfig,
) (*Verifier, error) {
	if ctx == nil {
		return nil, fmt.Errorf(
			"create Cognito verifier: context is nil",
		)
	}

	if err := validateEndpoint(
		"issuer",
		config.Issuer,
	); err != nil {
		return nil, fmt.Errorf(
			"create Cognito verifier: %w",
			err,
		)
	}

	if strings.TrimSpace(config.ClientID) == "" {
		return nil, fmt.Errorf(
			"create Cognito verifier: client ID is empty",
		)
	}

	provider, err := oidc.NewProvider(
		ctx,
		config.Issuer,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create Cognito verifier: discover provider: %w",
			err,
		)
	}

	idTokenVerifier := provider.Verifier(
		&oidc.Config{
			ClientID: config.ClientID,
			SupportedSigningAlgs: []string{
				"RS256",
			},
		},
	)

	return &Verifier{
		idTokenVerifier: idTokenVerifier,
	}, nil
}

func (
	verifier *Verifier,
) Verify(
	ctx context.Context,
	rawIDToken string,
	expectedNonce string,
) (Identity, error) {
	if verifier == nil ||
		verifier.idTokenVerifier == nil {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: verifier is nil",
		)
	}

	if ctx == nil {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: context is nil",
		)
	}

	if strings.TrimSpace(rawIDToken) == "" {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: ID token is empty",
		)
	}

	if strings.TrimSpace(expectedNonce) == "" {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: expected nonce is empty",
		)
	}

	idToken, err :=
		verifier.idTokenVerifier.Verify(
			ctx,
			rawIDToken,
		)
	if err != nil {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: verify signature and standard claims: %w",
			err,
		)
	}

	if strings.TrimSpace(idToken.Subject) == "" {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: subject is empty",
		)
	}

	var claims idTokenClaims

	if err := idToken.Claims(&claims); err != nil {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: decode claims: %w",
			err,
		)
	}

	if strings.TrimSpace(claims.Nonce) == "" {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: nonce is empty",
		)
	}

	if !equalConstantTime(
		claims.Nonce,
		expectedNonce,
	) {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: nonce does not match",
		)
	}

	if claims.TokenUse != "id" {
		return Identity{}, fmt.Errorf(
			"verify Cognito ID token: token_use = %q, want %q",
			claims.TokenUse,
			"id",
		)
	}

	return Identity{
		Subject: idToken.Subject,
		Email:   claims.Email,
	}, nil
}

func equalConstantTime(
	left string,
	right string,
) bool {
	return subtle.ConstantTimeCompare(
		[]byte(left),
		[]byte(right),
	) == 1
}
