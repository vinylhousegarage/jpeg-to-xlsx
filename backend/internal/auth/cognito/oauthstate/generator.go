package oauthstate

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

const randomValueByteLength = 32

type Generator struct {
	random io.Reader
	now    func() time.Time
	ttl    time.Duration
}

func NewGenerator(
	ttl time.Duration,
) (
	*Generator,
	error,
) {
	return newGenerator(
		cryptorand.Reader,
		time.Now,
		ttl,
	)
}

func newGenerator(
	random io.Reader,
	now func() time.Time,
	ttl time.Duration,
) (
	*Generator,
	error,
) {
	if random == nil {
		return nil, fmt.Errorf(
			"create OAuth state generator: random source is nil",
		)
	}

	if now == nil {
		return nil, fmt.Errorf(
			"create OAuth state generator: clock is nil",
		)
	}

	if ttl <= 0 {
		return nil, fmt.Errorf(
			"create OAuth state generator: TTL must be positive",
		)
	}

	return &Generator{
		random: random,
		now:    now,
		ttl:    ttl,
	}, nil
}

func (g *Generator) Generate() (
	State,
	error,
) {
	if g == nil {
		return State{}, fmt.Errorf(
			"generate OAuth state: generator is nil",
		)
	}

	stateValue, err :=
		g.generateRandomValue()
	if err != nil {
		return State{}, fmt.Errorf(
			"generate OAuth state value: %w",
			err,
		)
	}

	codeVerifier, err :=
		g.generateRandomValue()
	if err != nil {
		return State{}, fmt.Errorf(
			"generate PKCE code verifier: %w",
			err,
		)
	}

	nonce, err := g.generateRandomValue()
	if err != nil {
		return State{}, fmt.Errorf(
			"generate OIDC nonce: %w",
			err,
		)
	}

	createdAt := g.now().UTC()

	return State{
		Value:        stateValue,
		CodeVerifier: codeVerifier,
		Nonce:        nonce,
		CreatedAt:    createdAt,
		ExpiresAt:    createdAt.Add(g.ttl),
	}, nil
}

func (g *Generator) generateRandomValue() (
	string,
	error,
) {
	randomBytes := make(
		[]byte,
		randomValueByteLength,
	)

	if _, err := io.ReadFull(
		g.random,
		randomBytes,
	); err != nil {
		return "", fmt.Errorf(
			"read random bytes: %w",
			err,
		)
	}

	value :=
		base64.RawURLEncoding.EncodeToString(
			randomBytes,
		)

	return value, nil
}
