package session

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

const sessionIDByteLength = 32

type IDGenerator struct {
	random io.Reader
}

func NewIDGenerator() *IDGenerator {
	return newIDGenerator(
		cryptorand.Reader,
	)
}

func newIDGenerator(
	random io.Reader,
) *IDGenerator {
	return &IDGenerator{
		random: random,
	}
}

func (g *IDGenerator) Generate() (
	string,
	error,
) {
	if g == nil || g.random == nil {
		return "", fmt.Errorf(
			"generate session ID: random source is nil",
		)
	}

	randomBytes := make(
		[]byte,
		sessionIDByteLength,
	)

	if _, err := io.ReadFull(
		g.random,
		randomBytes,
	); err != nil {
		return "", fmt.Errorf(
			"generate session ID: read random bytes: %w",
			err,
		)
	}

	sessionID :=
		base64.RawURLEncoding.EncodeToString(
			randomBytes,
		)

	return sessionID, nil
}

func HashID(
	sessionID string,
) (
	string,
	error,
) {
	if strings.TrimSpace(sessionID) == "" {
		return "", fmt.Errorf(
			"hash session ID: session ID is empty",
		)
	}

	hash := sha256.Sum256(
		[]byte(sessionID),
	)

	return hex.EncodeToString(
		hash[:],
	), nil
}
