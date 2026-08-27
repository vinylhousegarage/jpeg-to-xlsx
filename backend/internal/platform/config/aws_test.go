package config

import "testing"

func TestLoadAWSConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		t.Setenv("AWS_REGION", "")
		t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "")

		cfg := loadAWSConfig()

		if cfg.Region != defaultAWSRegion {
			t.Errorf(
				"expected Region %q, got %q",
				defaultAWSRegion,
				cfg.Region,
			)
		}

		if cfg.IsLambda {
			t.Error("expected IsLambda false")
		}
	})

	t.Run("custom values", func(t *testing.T) {
		t.Setenv("AWS_REGION", "us-east-1")
		t.Setenv(
			"AWS_LAMBDA_FUNCTION_NAME",
			"my-lambda-function",
		)

		cfg := loadAWSConfig()

		if cfg.Region != "us-east-1" {
			t.Errorf(
				"expected Region %q, got %q",
				"us-east-1",
				cfg.Region,
			)
		}

		if !cfg.IsLambda {
			t.Error("expected IsLambda true")
		}
	})
}
