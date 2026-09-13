package storage

import "testing"

func TestBuildObjectKey(t *testing.T) {
	t.Parallel()

	const (
		cognitoSub = "test-cognito-sub"
		shotNumber = "001"
		want       = "test-cognito-sub/001.jpg"
	)

	got := BuildObjectKey(
		cognitoSub,
		shotNumber,
	)

	if got != want {
		t.Errorf("BuildObjectKey() = %q, want %q", got, want)
	}
}

func TestExtractCognitoSub_Success(t *testing.T) {
	t.Parallel()

	const objectKey = "test-cognito-sub/001.jpg"

	got, err := ExtractCognitoSub(objectKey)
	if err != nil {
		t.Fatalf("ExtractCognitoSub() error = %v", err)
	}

	const want = "test-cognito-sub"

	if got != want {
		t.Errorf("ExtractCognitoSub() = %q, want %q", got, want)
	}
}
