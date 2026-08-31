package config

import "testing"

func TestLoadBedrockConfig(
	t *testing.T,
) {
	tests := []struct {
		name               string
		modelID            string
		promptFileName     string
		wantModelID        string
		wantPromptFileName string
	}{
		{
			name:               "defaults prompt file name",
			modelID:            testBedrockModelID,
			promptFileName:     "",
			wantModelID:        testBedrockModelID,
			wantPromptFileName: defaultPromptFileName,
		},
		{
			name:               "custom prompt file name",
			modelID:            "test-model",
			promptFileName:     "custom-prompt.txt",
			wantModelID:        "test-model",
			wantPromptFileName: "custom-prompt.txt",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Setenv(
					"BEDROCK_MODEL_ID",
					test.modelID,
				)
				t.Setenv(
					"PROMPT_FILE_NAME",
					test.promptFileName,
				)

				config, err :=
					loadBedrockConfig()
				if err != nil {
					t.Fatalf(
						"loadBedrockConfig() error = %v",
						err,
					)
				}

				if config.ModelID !=
					test.wantModelID {
					t.Errorf(
						"ModelID = %q, want %q",
						config.ModelID,
						test.wantModelID,
					)
				}

				if config.PromptFileName !=
					test.wantPromptFileName {
					t.Errorf(
						"PromptFileName = %q, want %q",
						config.PromptFileName,
						test.wantPromptFileName,
					)
				}
			},
		)
	}
}

func TestLoadBedrockConfigRequiresModelID(
	t *testing.T,
) {
	t.Setenv(
		"BEDROCK_MODEL_ID",
		"",
	)
	t.Setenv(
		"PROMPT_FILE_NAME",
		"",
	)

	config, err := loadBedrockConfig()
	if err == nil {
		t.Fatal(
			"loadBedrockConfig() error = nil, want an error",
		)
	}

	if config != (BedrockConfig{}) {
		t.Errorf(
			"loadBedrockConfig() config = %+v, want zero value",
			config,
		)
	}

	const wantError = "BEDROCK_MODEL_ID is required"

	if err.Error() != wantError {
		t.Errorf(
			"loadBedrockConfig() error = %q, want %q",
			err,
			wantError,
		)
	}
}
