package config

import "testing"

func TestLoadAPIStorageConfig(
	t *testing.T,
) {
	t.Run(
		"success",
		func(t *testing.T) {
			t.Setenv(
				"INPUT_BUCKET_NAME",
				testInputBucketName,
			)
			t.Setenv(
				"OUTPUT_BUCKET_NAME",
				"",
			)

			config, err :=
				loadAPIStorageConfig()
			if err != nil {
				t.Fatalf(
					"loadAPIStorageConfig() error = %v",
					err,
				)
			}

			if config.InputBucketName !=
				testInputBucketName {
				t.Errorf(
					"InputBucketName = %q, want %q",
					config.InputBucketName,
					testInputBucketName,
				)
			}

			if config.OutputBucketName != "" {
				t.Errorf(
					"OutputBucketName = %q, want empty",
					config.OutputBucketName,
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

			config, err :=
				loadAPIStorageConfig()
			if err == nil {
				t.Fatal(
					"loadAPIStorageConfig() error = nil, want an error",
				)
			}

			if config !=
				(StorageConfig{}) {
				t.Errorf(
					"loadAPIStorageConfig() config = %+v, want zero value",
					config,
				)
			}

			const wantError = "INPUT_BUCKET_NAME is required"

			if err.Error() != wantError {
				t.Errorf(
					"loadAPIStorageConfig() error = %q, want %q",
					err,
					wantError,
				)
			}
		},
	)
}
