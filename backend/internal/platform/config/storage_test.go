package config

import (
	"strings"
	"testing"
)

func TestLoadAPIStorageConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv(
			"INPUT_BUCKET_NAME",
			"my-input-bucket",
		)
		t.Setenv(
			"OUTPUT_BUCKET_NAME",
			"",
		)

		cfg, err := loadAPIStorageConfig()
		if err != nil {
			t.Fatalf(
				"loadAPIStorageConfig() error = %v",
				err,
			)
		}

		if cfg.InputBucketName != "my-input-bucket" {
			t.Errorf(
				"expected InputBucketName %q, got %q",
				"my-input-bucket",
				cfg.InputBucketName,
			)
		}

		if cfg.OutputBucketName != "" {
			t.Errorf(
				"expected OutputBucketName to be empty, got %q",
				cfg.OutputBucketName,
			)
		}
	})

	t.Run("missing input bucket", func(t *testing.T) {
		t.Setenv(
			"INPUT_BUCKET_NAME",
			"",
		)

		_, err := loadAPIStorageConfig()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(
			err.Error(),
			"INPUT_BUCKET_NAME is required",
		) {
			t.Errorf(
				"expected error containing %q, got %q",
				"INPUT_BUCKET_NAME is required",
				err.Error(),
			)
		}
	})
}

func TestLoadProcessorStorageConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv(
			"INPUT_BUCKET_NAME",
			"my-input-bucket",
		)
		t.Setenv(
			"OUTPUT_BUCKET_NAME",
			"my-output-bucket",
		)

		cfg, err := loadProcessorStorageConfig()
		if err != nil {
			t.Fatalf(
				"loadProcessorStorageConfig() error = %v",
				err,
			)
		}

		if cfg.InputBucketName != "my-input-bucket" {
			t.Errorf(
				"expected InputBucketName %q, got %q",
				"my-input-bucket",
				cfg.InputBucketName,
			)
		}

		if cfg.OutputBucketName != "my-output-bucket" {
			t.Errorf(
				"expected OutputBucketName %q, got %q",
				"my-output-bucket",
				cfg.OutputBucketName,
			)
		}
	})

	tests := []struct {
		name             string
		inputBucketName  string
		outputBucketName string
		expectedErr      string
	}{
		{
			name:             "missing input bucket",
			inputBucketName:  "",
			outputBucketName: "my-output-bucket",
			expectedErr:      "INPUT_BUCKET_NAME is required",
		},
		{
			name:             "missing output bucket",
			inputBucketName:  "my-input-bucket",
			outputBucketName: "",
			expectedErr:      "OUTPUT_BUCKET_NAME is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(
				"INPUT_BUCKET_NAME",
				tt.inputBucketName,
			)
			t.Setenv(
				"OUTPUT_BUCKET_NAME",
				tt.outputBucketName,
			)

			_, err := loadProcessorStorageConfig()
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(
				err.Error(),
				tt.expectedErr,
			) {
				t.Errorf(
					"expected error containing %q, got %q",
					tt.expectedErr,
					err.Error(),
				)
			}
		})
	}
}

func TestLoadInputBucketName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv(
			"INPUT_BUCKET_NAME",
			"my-input-bucket",
		)

		got, err := loadInputBucketName()
		if err != nil {
			t.Fatalf(
				"loadInputBucketName() error = %v",
				err,
			)
		}

		if got != "my-input-bucket" {
			t.Errorf(
				"expected %q, got %q",
				"my-input-bucket",
				got,
			)
		}
	})

	t.Run("missing input bucket", func(t *testing.T) {
		t.Setenv(
			"INPUT_BUCKET_NAME",
			"",
		)

		_, err := loadInputBucketName()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(
			err.Error(),
			"INPUT_BUCKET_NAME is required",
		) {
			t.Errorf(
				"expected error containing %q, got %q",
				"INPUT_BUCKET_NAME is required",
				err.Error(),
			)
		}
	})
}
