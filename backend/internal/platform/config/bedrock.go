package config

import "os"

func loadBedrockConfig() (
	BedrockConfig,
	error,
) {
	modelID, err := loadRequiredEnv(
		"BEDROCK_MODEL_ID",
	)
	if err != nil {
		return BedrockConfig{}, err
	}

	promptFileName := os.Getenv(
		"PROMPT_FILE_NAME",
	)
	if promptFileName == "" {
		promptFileName =
			defaultPromptFileName
	}

	return BedrockConfig{
		ModelID:        modelID,
		PromptFileName: promptFileName,
	}, nil
}
