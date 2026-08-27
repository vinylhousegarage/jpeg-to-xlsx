package config

import (
	"strings"
	"testing"
)

func TestLoadBedrockConfig(t *testing.T) {
	t.Run("success with defaults", func(t *testing.T) {
		t.Setenv(
			"BEDROCK_MODEL_ID",
			"jp.anthropic.claude-sonnet-4-6",
		)
		t.Setenv(
			"PROMPT_FILE_NAME",
			"",
		)

		cfg, err := loadBedrockConfig()
		if err != nil {
			t.Fatalf(
				"loadBedrockConfig() error = %v",
				err,
			)
		}

		if cfg.ModelID != "jp.anthropic.claude-sonnet-4-6" {
			t.Errorf(
				"expected ModelID %q, got %q",
				"jp.anthropic.claude-sonnet-4-6",
				cfg.ModelID,
			)
		}

		if cfg.PromptFileName != defaultPromptFileName {
			t.Errorf(
				"expected PromptFileName %q, got %q",
				defaultPromptFileName,
				cfg.PromptFileName,
			)
		}
	})

	t.Run("success with custom prompt file", func(t *testing.T) {
		t.Setenv(
			"BEDROCK_MODEL_ID",
			"test-model",
		)
		t.Setenv(
			"PROMPT_FILE_NAME",
			"custom-prompt.txt",
		)

		cfg, err := loadBedrockConfig()
		if err != nil {
			t.Fatalf(
				"loadBedrockConfig() error = %v",
				err,
			)
		}

		if cfg.ModelID != "test-model" {
			t.Errorf(
				"expected ModelID %q, got %q",
				"test-model",
				cfg.ModelID,
			)
		}

		if cfg.PromptFileName != "custom-prompt.txt" {
			t.Errorf(
				"expected PromptFileName %q, got %q",
				"custom-prompt.txt",
				cfg.PromptFileName,
			)
		}
	})

	t.Run("missing model ID", func(t *testing.T) {
		t.Setenv(
			"BEDROCK_MODEL_ID",
			"",
		)
		t.Setenv(
			"PROMPT_FILE_NAME",
			"",
		)

		_, err := loadBedrockConfig()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(
			err.Error(),
			"BEDROCK_MODEL_ID is required",
		) {
			t.Errorf(
				"expected error containing %q, got %q",
				"BEDROCK_MODEL_ID is required",
				err.Error(),
			)
		}
	})
}
