package prompts

import (
	"strings"
	"testing"
)

func TestLoadPrompt(t *testing.T) {
	content, err := LoadPrompt("extractor.txt")
	if err != nil {
		t.Fatalf("Failed to load prompt: %v", err)
	}

	if content == "" {
		t.Error("Loaded prompt is empty")
	}

	if !strings.Contains(content, "<system_role>") {
		t.Error("Prompt does not contain expected system_role tag")
	}
}
