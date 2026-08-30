package session

import (
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"testing"
)

const testSessionIDByteLength = 32

var errTestRandom = errors.New(
	"random source failed",
)

type errorReader struct {
	err error
}

func (r errorReader) Read(
	_ []byte,
) (int, error) {
	return 0, r.err
}

func TestIDGeneratorGenerate(t *testing.T) {
	t.Parallel()

	randomBytes := strings.Repeat(
		"a",
		testSessionIDByteLength,
	)

	generator := newIDGenerator(
		strings.NewReader(randomBytes),
	)

	sessionID, err := generator.Generate()
	if err != nil {
		t.Fatalf(
			"Generate() error = %v",
			err,
		)
	}

	if sessionID == "" {
		t.Fatal(
			"Generate() returned an empty session ID",
		)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(
		sessionID,
	)
	if err != nil {
		t.Fatalf(
			"decode generated session ID: %v",
			err,
		)
	}

	if len(decoded) != testSessionIDByteLength {
		t.Errorf(
			"decoded session ID length = %d, want %d",
			len(decoded),
			testSessionIDByteLength,
		)
	}

	if string(decoded) != randomBytes {
		t.Errorf(
			"decoded session ID = %q, want %q",
			string(decoded),
			randomBytes,
		)
	}

	if strings.Contains(sessionID, "=") {
		t.Errorf(
			"session ID %q contains Base64 padding",
			sessionID,
		)
	}
}

func TestIDGeneratorGenerateReturnsUniqueIDs(
	t *testing.T,
) {
	t.Parallel()

	generator := NewIDGenerator()

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

	if first == second {
		t.Errorf(
			"Generate() returned duplicate session IDs: %q",
			first,
		)
	}
}

func TestIDGeneratorGenerateReturnsRandomSourceError(
	t *testing.T,
) {
	t.Parallel()

	generator := newIDGenerator(
		errorReader{
			err: errTestRandom,
		},
	)

	sessionID, err := generator.Generate()
	if err == nil {
		t.Fatal(
			"Generate() error = nil, want an error",
		)
	}

	if sessionID != "" {
		t.Errorf(
			"Generate() session ID = %q, want empty",
			sessionID,
		)
	}

	if !errors.Is(err, errTestRandom) {
		t.Errorf(
			"Generate() error = %v, want wrapped %v",
			err,
			errTestRandom,
		)
	}
}

func TestIDGeneratorGenerateReturnsShortReadError(
	t *testing.T,
) {
	t.Parallel()

	generator := newIDGenerator(
		strings.NewReader("too-short"),
	)

	sessionID, err := generator.Generate()
	if err == nil {
		t.Fatal(
			"Generate() error = nil, want an error",
		)
	}

	if sessionID != "" {
		t.Errorf(
			"Generate() session ID = %q, want empty",
			sessionID,
		)
	}

	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf(
			"Generate() error = %v, want wrapped %v",
			err,
			io.ErrUnexpectedEOF,
		)
	}
}

func TestHashID(t *testing.T) {
	t.Parallel()

	const sessionID = "test-session-id"

	first, err := HashID(sessionID)
	if err != nil {
		t.Fatalf(
			"first HashID() error = %v",
			err,
		)
	}

	second, err := HashID(sessionID)
	if err != nil {
		t.Fatalf(
			"second HashID() error = %v",
			err,
		)
	}

	if first == "" {
		t.Fatal(
			"HashID() returned an empty hash",
		)
	}

	if first != second {
		t.Errorf(
			"HashID() is not deterministic: %q != %q",
			first,
			second,
		)
	}

	if len(first) != 64 {
		t.Errorf(
			"HashID() length = %d, want 64",
			len(first),
		)
	}
}

func TestHashIDReturnsDifferentHashes(
	t *testing.T,
) {
	t.Parallel()

	first, err := HashID("first-session-id")
	if err != nil {
		t.Fatalf(
			"HashID(first) error = %v",
			err,
		)
	}

	second, err := HashID("second-session-id")
	if err != nil {
		t.Fatalf(
			"HashID(second) error = %v",
			err,
		)
	}

	if first == second {
		t.Errorf(
			"HashID() returned the same hash for different IDs: %q",
			first,
		)
	}
}

func TestHashIDRejectsEmptyID(t *testing.T) {
	t.Parallel()

	hash, err := HashID("")
	if err == nil {
		t.Fatal(
			"HashID() error = nil, want an error",
		)
	}

	if hash != "" {
		t.Errorf(
			"HashID() hash = %q, want empty",
			hash,
		)
	}
}

func TestHashIDRejectsWhitespaceOnlyID(
	t *testing.T,
) {
	t.Parallel()

	hash, err := HashID(" \t\n")
	if err == nil {
		t.Fatal(
			"HashID() error = nil, want an error",
		)
	}

	if hash != "" {
		t.Errorf(
			"HashID() hash = %q, want empty",
			hash,
		)
	}
}
