package oauthstate

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"
)

func TestGeneratorGenerateReturnsRandomSourceError(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name   string
		prefix []byte
	}{
		{
			name:   "state",
			prefix: nil,
		},
		{
			name: "code verifier",
			prefix: make(
				[]byte,
				testGeneratorRandomValueByteLength,
			),
		},
		{
			name: "nonce",
			prefix: make(
				[]byte,
				testGeneratorRandomValueByteLength*2,
			),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				random := io.MultiReader(
					bytes.NewReader(
						test.prefix,
					),
					generatorErrorReader{
						err: errGeneratorRandom,
					},
				)

				generator, err := newGenerator(
					random,
					time.Now,
					10*time.Minute,
				)
				if err != nil {
					t.Fatalf(
						"newGenerator() error = %v",
						err,
					)
				}

				generated, err :=
					generator.Generate()
				if err == nil {
					t.Fatal(
						"Generate() error = nil, want an error",
					)
				}

				if generated != (State{}) {
					t.Errorf(
						"Generate() state = %+v, want zero value",
						generated,
					)
				}

				if !errors.Is(
					err,
					errGeneratorRandom,
				) {
					t.Errorf(
						"Generate() error = %v, want wrapped %v",
						err,
						errGeneratorRandom,
					)
				}
			},
		)
	}
}

func TestNewGeneratorRejectsInvalidTTL(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name string
		ttl  time.Duration
	}{
		{
			name: "zero",
			ttl:  0,
		},
		{
			name: "negative",
			ttl:  -time.Minute,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				generator, err := NewGenerator(
					test.ttl,
				)
				if err == nil {
					t.Fatal(
						"NewGenerator() error = nil, want an error",
					)
				}

				if generator != nil {
					t.Errorf(
						"NewGenerator() generator = %v, want nil",
						generator,
					)
				}
			},
		)
	}
}

func TestNewGeneratorRejectsNilRandomSource(
	t *testing.T,
) {
	t.Parallel()

	generator, err := newGenerator(
		nil,
		time.Now,
		10*time.Minute,
	)
	if err == nil {
		t.Fatal(
			"newGenerator() error = nil, want an error",
		)
	}

	if generator != nil {
		t.Errorf(
			"newGenerator() generator = %v, want nil",
			generator,
		)
	}
}

func TestNewGeneratorRejectsNilClock(
	t *testing.T,
) {
	t.Parallel()

	generator, err := newGenerator(
		bytes.NewReader(
			make(
				[]byte,
				testGeneratorRandomValueByteLength*3,
			),
		),
		nil,
		10*time.Minute,
	)
	if err == nil {
		t.Fatal(
			"newGenerator() error = nil, want an error",
		)
	}

	if generator != nil {
		t.Errorf(
			"newGenerator() generator = %v, want nil",
			generator,
		)
	}
}
