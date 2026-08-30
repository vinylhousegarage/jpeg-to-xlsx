package oauthstate

import (
	"bytes"
	"testing"
	"time"
)

func TestGeneratorGenerate(t *testing.T) {
	t.Parallel()

	const ttl = 10 * time.Minute

	now := time.Date(
		2026,
		time.August,
		30,
		12,
		34,
		56,
		0,
		time.UTC,
	)

	stateBytes := bytes.Repeat(
		[]byte{'a'},
		testGeneratorRandomValueByteLength,
	)
	verifierBytes := bytes.Repeat(
		[]byte{'b'},
		testGeneratorRandomValueByteLength,
	)
	nonceBytes := bytes.Repeat(
		[]byte{'c'},
		testGeneratorRandomValueByteLength,
	)

	randomBytes := make(
		[]byte,
		0,
		testGeneratorRandomValueByteLength*3,
	)
	randomBytes = append(
		randomBytes,
		stateBytes...,
	)
	randomBytes = append(
		randomBytes,
		verifierBytes...,
	)
	randomBytes = append(
		randomBytes,
		nonceBytes...,
	)

	generator, err := newGenerator(
		bytes.NewReader(randomBytes),
		func() time.Time {
			return now
		},
		ttl,
	)
	if err != nil {
		t.Fatalf(
			"newGenerator() error = %v",
			err,
		)
	}

	generated, err := generator.Generate()
	if err != nil {
		t.Fatalf(
			"Generate() error = %v",
			err,
		)
	}

	assertRandomValue(
		t,
		"State.Value",
		generated.Value,
		stateBytes,
	)
	assertRandomValue(
		t,
		"State.CodeVerifier",
		generated.CodeVerifier,
		verifierBytes,
	)
	assertRandomValue(
		t,
		"State.Nonce",
		generated.Nonce,
		nonceBytes,
	)

	if generated.Value ==
		generated.CodeVerifier {
		t.Error(
			"State.Value and State.CodeVerifier are equal",
		)
	}

	if generated.Value == generated.Nonce {
		t.Error(
			"State.Value and State.Nonce are equal",
		)
	}

	if generated.CodeVerifier ==
		generated.Nonce {
		t.Error(
			"State.CodeVerifier and State.Nonce are equal",
		)
	}

	if !generated.CreatedAt.Equal(now) {
		t.Errorf(
			"State.CreatedAt = %v, want %v",
			generated.CreatedAt,
			now,
		)
	}

	wantExpiresAt := now.Add(ttl)

	if !generated.ExpiresAt.Equal(
		wantExpiresAt,
	) {
		t.Errorf(
			"State.ExpiresAt = %v, want %v",
			generated.ExpiresAt,
			wantExpiresAt,
		)
	}
}

func TestGeneratorGenerateReturnsUniqueValues(
	t *testing.T,
) {
	t.Parallel()

	generator, err := NewGenerator(
		10 * time.Minute,
	)
	if err != nil {
		t.Fatalf(
			"NewGenerator() error = %v",
			err,
		)
	}

	first, err := generator.Generate()
	if err != nil {
		t.Fatalf(
			"first Generate() error = %v",
			err,
		)
	}

	second, err := generator.Generate()
	if err != nil {
		t.Fatalf(
			"second Generate() error = %v",
			err,
		)
	}

	if first.Value == second.Value {
		t.Errorf(
			"Generate() returned duplicate state values: %q",
			first.Value,
		)
	}

	if first.CodeVerifier ==
		second.CodeVerifier {
		t.Errorf(
			"Generate() returned duplicate code verifiers: %q",
			first.CodeVerifier,
		)
	}

	if first.Nonce == second.Nonce {
		t.Errorf(
			"Generate() returned duplicate nonces: %q",
			first.Nonce,
		)
	}
}
