package config

import "testing"

func TestLoadInputBucketName(
	t *testing.T,
) {
	t.Run(
		"success",
		func(t *testing.T) {
			t.Setenv(
				"INPUT_BUCKET_NAME",
				testInputBucketName,
			)

			value, err :=
				loadInputBucketName()
			if err != nil {
				t.Fatalf(
					"loadInputBucketName() error = %v",
					err,
				)
			}

			if value !=
				testInputBucketName {
				t.Errorf(
					"loadInputBucketName() value = %q, want %q",
					value,
					testInputBucketName,
				)
			}
		},
	)

	t.Run(
		"missing input bucket",
		func(t *testing.T) {
			t.Setenv(
				"INPUT_BUCKET_NAME",
				"",
			)

			value, err :=
				loadInputBucketName()
			if err == nil {
				t.Fatal(
					"loadInputBucketName() error = nil, want an error",
				)
			}

			if value != "" {
				t.Errorf(
					"loadInputBucketName() value = %q, want empty",
					value,
				)
			}

			const wantError = "INPUT_BUCKET_NAME is required"

			if err.Error() != wantError {
				t.Errorf(
					"loadInputBucketName() error = %q, want %q",
					err,
					wantError,
				)
			}
		},
	)
}
