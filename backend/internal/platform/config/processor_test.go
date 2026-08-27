package config

import "testing"

func TestLoadProcessor(t *testing.T) {
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
		"INPUT_BUCKET_NAME",
		"my-test-bucket",
	)
	t.Setenv(
		"OUTPUT_BUCKET_NAME",
		"my-output-bucket",
	)

	t.Setenv(
		"BEDROCK_MODEL_ID",
		"jp.anthropic.claude-sonnet-4-6",
	)
	t.Setenv(
		"PROMPT_FILE_NAME",
		defaultPromptFileName,
	)

	t.Setenv(
		"SLACK_TOKEN_TABLE_NAME",
		"slack-tokens",
	)

	cfg, err := LoadProcessor()
	if err != nil {
		t.Fatalf(
			"LoadProcessor() error = %v",
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

	if cfg.Storage.InputBucketName != "my-test-bucket" {
		t.Errorf(
			"expected Storage.InputBucketName %q, got %q",
			"my-test-bucket",
			cfg.Storage.InputBucketName,
		)
	}

	if cfg.Storage.OutputBucketName != "my-output-bucket" {
		t.Errorf(
			"expected Storage.OutputBucketName %q, got %q",
			"my-output-bucket",
			cfg.Storage.OutputBucketName,
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

	if cfg.SlackToken.TokenTableName != "slack-tokens" {
		t.Errorf(
			"expected SlackToken.TokenTableName %q, got %q",
			"slack-tokens",
			cfg.SlackToken.TokenTableName,
		)
	}
}
