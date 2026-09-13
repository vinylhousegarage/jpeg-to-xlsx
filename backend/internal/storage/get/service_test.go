package get

import (
	"context"
	"errors"
	"testing"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockPresigner struct {
	mockResult *v4.PresignedHTTPRequest
	mockErr    error
	gotKey     string
}

func (m *mockPresigner) PresignGetObject(
	_ context.Context,
	params *s3.GetObjectInput,
	_ ...func(*s3.PresignOptions),
) (
	*v4.PresignedHTTPRequest,
	error,
) {
	if params.Key != nil {
		m.gotKey = *params.Key
	}

	return m.mockResult, m.mockErr
}

func TestService_GeneratePresignURL_Success(
	t *testing.T,
) {
	t.Parallel()

	const (
		cognitoSub  = "test-cognito-sub"
		shotNumber  = "001"
		expectedURL = "https://example.com/presigned-url"
		expectedKey = "test-cognito-sub/001.jpg"
	)

	mock := &mockPresigner{
		mockResult: &v4.PresignedHTTPRequest{
			URL: expectedURL,
		},
	}

	svc := NewService(mock)

	url, expiresAt, err := svc.GeneratePresignURL(
		context.Background(),
		cognitoSub,
		shotNumber,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if url != expectedURL {
		t.Errorf("expected URL %q, got %q", expectedURL, url)
	}

	if mock.gotKey != expectedKey {
		t.Errorf("expected key %q, got %q", expectedKey, mock.gotKey)
	}

	remaining := time.Until(expiresAt)
	if remaining > 16*time.Minute || remaining < 14*time.Minute {
		t.Errorf("unexpected expiration time: %v", expiresAt)
	}
}

func TestService_GeneratePresignURL_Error(t *testing.T) {
	t.Parallel()

	mock := &mockPresigner{
		mockErr: errors.New("aws error"),
	}

	svc := NewService(mock)

	_, _, err := svc.GeneratePresignURL(
		context.Background(),
		"test-cognito-sub",
		"001",
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
