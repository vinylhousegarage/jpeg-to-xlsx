package config

import "testing"

func TestLoadProcessor(
	t *testing.T,
) {
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
		testInputBucketName,
	)
	t.Setenv(
		"OUTPUT_BUCKET_NAME",
		testOutputBucketName,
	)
	t.Setenv(
		"BEDROCK_MODEL_ID",
		testBedrockModelID,
	)
	t.Setenv(
		"PROMPT_FILE_NAME",
		defaultPromptFileName,
	)
	t.Setenv(
		"SLACK_TOKEN_TABLE_NAME",
		testSlackTokenTableName,
	)

	config, err := LoadProcessor()
	if err != nil {
		t.Fatalf(
			"LoadProcessor() error = %v",
			err,
		)
	}

	if config == nil {
		t.Fatal(
			"LoadProcessor() config = nil",
		)
	}

	if config.App.Env != appEnvLocal {
		t.Errorf(
			"App.Env = %q, want %q",
			config.App.Env,
			appEnvLocal,
		)
	}

	if config.AWS.Region !=
		defaultAWSRegion {
		t.Errorf(
			"AWS.Region = %q, want %q",
			config.AWS.Region,
			defaultAWSRegion,
		)
	}

	if config.Storage.InputBucketName !=
		testInputBucketName {
		t.Errorf(
			"Storage.InputBucketName = %q, want %q",
			config.Storage.InputBucketName,
			testInputBucketName,
		)
	}

	if config.Storage.OutputBucketName !=
		testOutputBucketName {
		t.Errorf(
			"Storage.OutputBucketName = %q, want %q",
			config.Storage.OutputBucketName,
			testOutputBucketName,
		)
	}

	if config.Bedrock.ModelID !=
		testBedrockModelID {
		t.Errorf(
			"Bedrock.ModelID = %q, want %q",
			config.Bedrock.ModelID,
			testBedrockModelID,
		)
	}

	if config.Bedrock.PromptFileName !=
		defaultPromptFileName {
		t.Errorf(
			"Bedrock.PromptFileName = %q, want %q",
			config.Bedrock.PromptFileName,
			defaultPromptFileName,
		)
	}

	if config.SlackToken.TokenTableName !=
		testSlackTokenTableName {
		t.Errorf(
			"SlackToken.TokenTableName = %q, want %q",
			config.SlackToken.TokenTableName,
			testSlackTokenTableName,
		)
	}
}
