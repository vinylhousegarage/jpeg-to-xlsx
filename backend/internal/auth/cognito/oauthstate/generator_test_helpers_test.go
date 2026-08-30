package oauthstate

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"
)

const testGeneratorRandomValueByteLength = 32

var errGeneratorRandom = errors.New(
	"random source failed",
)

type generatorErrorReader struct {
	err error
}

func (r generatorErrorReader) Read(
	_ []byte,
) (int, error) {
	return 0, r.err
}

func assertRandomValue(
	t *testing.T,
	name string,
	value string,
	wantBytes []byte,
) {
	t.Helper()

	if value == "" {
		t.Fatalf(
			"%s is empty",
			name,
		)
	}

	decoded, err :=
		base64.RawURLEncoding.DecodeString(
			value,
		)
	if err != nil {
		t.Fatalf(
			"decode %s: %v",
			name,
			err,
		)
	}

	if len(decoded) !=
		testGeneratorRandomValueByteLength {
		t.Errorf(
			"%s decoded length = %d, want %d",
			name,
			len(decoded),
			testGeneratorRandomValueByteLength,
		)
	}

	if !bytes.Equal(
		decoded,
		wantBytes,
	) {
		t.Errorf(
			"%s decoded value = %q, want %q",
			name,
			decoded,
			wantBytes,
		)
	}
}
