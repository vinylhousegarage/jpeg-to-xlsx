package config

import "testing"

func TestLoadBedrockPlayground(
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
		"BEDROCK_MODEL_ID",
		testBedrockModelID,
	)
	t.Setenv(
		"PROMPT_FILE_NAME",
		defaultPromptFileName,
	)

	config, err := LoadBedrockPlayground()
	if err != nil {
		t.Fatalf(
			"LoadBedrockPlayground() error = %v",
			err,
		)
	}

	if config == nil {
		t.Fatal(
			"LoadBedrockPlayground() config = nil",
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
}
