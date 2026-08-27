package oauth

import "testing"

func TestGenerateState(t *testing.T) {
	t.Parallel()

	state1, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState() returned error: %v", err)
	}

	state2, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState() returned error: %v", err)
	}

	if state1 == "" {
		t.Error("GenerateState() returned an empty string")
	}

	if state1 == state2 {
		t.Error("GenerateState() generated duplicate values")
	}
}
