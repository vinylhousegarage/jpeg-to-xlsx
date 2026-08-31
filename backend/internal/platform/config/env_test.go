package config

import "testing"

func TestLoadRequiredEnv(
	t *testing.T,
) {
	const (
		environmentName = "TEST_REQUIRED_ENV"

		wantValue = "test-value"
	)

	t.Setenv(
		environmentName,
		wantValue,
	)

	value, err := loadRequiredEnv(
		environmentName,
	)
	if err != nil {
		t.Fatalf(
			"loadRequiredEnv() error = %v",
			err,
		)
	}

	if value != wantValue {
		t.Errorf(
			"loadRequiredEnv() value = %q, want %q",
			value,
			wantValue,
		)
	}
}

func TestLoadRequiredEnvRejectsEmptyValue(
	t *testing.T,
) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "empty",
			value: "",
		},
		{
			name:  "spaces",
			value: "   ",
		},
		{
			name:  "tab",
			value: "\t",
		},
		{
			name:  "newline",
			value: "\n",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				const environmentName = "TEST_REQUIRED_ENV"

				t.Setenv(
					environmentName,
					test.value,
				)

				value, err :=
					loadRequiredEnv(
						environmentName,
					)
				if err == nil {
					t.Fatal(
						"loadRequiredEnv() error = nil, want an error",
					)
				}

				if value != "" {
					t.Errorf(
						"loadRequiredEnv() value = %q, want empty",
						value,
					)
				}

				const wantError = "TEST_REQUIRED_ENV is required"

				if err.Error() != wantError {
					t.Errorf(
						"loadRequiredEnv() error = %q, want %q",
						err,
						wantError,
					)
				}
			},
		)
	}
}
