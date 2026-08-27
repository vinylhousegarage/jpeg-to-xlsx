package config

import "testing"

func TestLoadBedrockPlayground(t *testing.T) {
	t.Setenv(
		"APP_ENV",
		appEnvLocal,
	)
	t.Setenv(
		"AWS_LAMBDA_FUNCTION_NAME",
		"",
	)
	t.Setenv(
		"AWS_REGION",
		defaultAWSRegion,
	)
	t.Setenv(
		"BEDROCK_MODEL_ID",
		"jp.anthropic.claude-sonnet-4-6",
	)
	t.Setenv(
		"PROMPT_FILE_NAME",
		defaultPromptFileName,
	)

	cfg, err := LoadBedrockPlayground()
	if err != nil {
		t.Fatalf(
			"LoadBedrockPlayground() error = %v",
			err,
		)
	}

	if cfg.App.Env != appEnvLocal {
		t.Errorf(
			"expected App.Env %q, got %q",
			appEnvLocal,
			cfg.App.Env,
		)
	}

	if cfg.AWS.Region != defaultAWSRegion {
		t.Errorf(
			"expected AWS.Region %q, got %q",
			defaultAWSRegion,
			cfg.AWS.Region,
		)
	}

	if cfg.Bedrock.ModelID != "jp.anthropic.claude-sonnet-4-6" {
		t.Errorf(
			"expected Bedrock.ModelID %q, got %q",
			"jp.anthropic.claude-sonnet-4-6",
			cfg.Bedrock.ModelID,
		)
	}

	if cfg.Bedrock.PromptFileName != defaultPromptFileName {
		t.Errorf(
			"expected Bedrock.PromptFileName %q, got %q",
			defaultPromptFileName,
			cfg.Bedrock.PromptFileName,
		)
	}
}
