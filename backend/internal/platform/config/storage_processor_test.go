package config

import "testing"

func TestLoadProcessorStorageConfig(
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
				testOutputBucketName,
			)

			config, err :=
				loadProcessorStorageConfig()
			if err != nil {
				t.Fatalf(
					"loadProcessorStorageConfig() error = %v",
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

			if config.OutputBucketName !=
				testOutputBucketName {
				t.Errorf(
					"OutputBucketName = %q, want %q",
					config.OutputBucketName,
					testOutputBucketName,
				)
			}
		},
	)

	tests := []struct {
		name             string
		inputBucketName  string
		outputBucketName string
		wantError        string
	}{
		{
			name:             "missing input bucket",
			inputBucketName:  "",
			outputBucketName: testOutputBucketName,
			wantError:        "INPUT_BUCKET_NAME is required",
		},
		{
			name:             "missing output bucket",
			inputBucketName:  testInputBucketName,
			outputBucketName: "",
			wantError:        "OUTPUT_BUCKET_NAME is required",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Setenv(
					"INPUT_BUCKET_NAME",
					test.inputBucketName,
				)
				t.Setenv(
					"OUTPUT_BUCKET_NAME",
					test.outputBucketName,
				)

				config, err :=
					loadProcessorStorageConfig()
				if err == nil {
					t.Fatal(
						"loadProcessorStorageConfig() error = nil, want an error",
					)
				}

				if config !=
					(StorageConfig{}) {
					t.Errorf(
						"loadProcessorStorageConfig() config = %+v, want zero value",
						config,
					)
				}

				if err.Error() !=
					test.wantError {
					t.Errorf(
						"loadProcessorStorageConfig() error = %q, want %q",
						err,
						test.wantError,
					)
				}
			},
		)
	}
}
