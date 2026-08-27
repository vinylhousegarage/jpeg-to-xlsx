package put

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// モック用の構造体
type mockPresigner struct {
	mockResult *v4.PresignedHTTPRequest
	mockErr    error
	gotInput   *s3.PutObjectInput
}

func (m *mockPresigner) PresignPutObject(
	_ context.Context,
	params *s3.PutObjectInput,
	_ ...func(*s3.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {
	m.gotInput = params

	return m.mockResult, m.mockErr
}

// 正常系
func TestService_GeneratePresignURL_Success(t *testing.T) {
	t.Parallel()

	const (
		shotNumber  = "001"
		expectedURL = "https://example.com/presigned-url"
	)

	mock := &mockPresigner{
		mockResult: &v4.PresignedHTTPRequest{
			URL: expectedURL,
		},
	}

	svc := NewService(mock)

	url, expiresAt, err := svc.GeneratePresignURL(
		context.Background(),
		shotNumber,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if url != expectedURL {
		t.Errorf(
			"expected URL %q, got %q",
			expectedURL,
			url,
		)
	}

	if mock.gotInput == nil {
		t.Fatal("PresignPutObject() input is nil")
	}

	if aws.ToString(mock.gotInput.Key) != "001.jpg" {
		t.Errorf(
			"Key = %q, want %q",
			aws.ToString(mock.gotInput.Key),
			"001.jpg",
		)
	}

	if aws.ToString(mock.gotInput.ContentType) != "image/jpeg" {
		t.Errorf(
			"ContentType = %q, want %q",
			aws.ToString(mock.gotInput.ContentType),
			"image/jpeg",
		)
	}

	if mock.gotInput.Metadata["shot-number"] != shotNumber {
		t.Errorf(
			"Metadata[shot-number] = %q, want %q",
			mock.gotInput.Metadata["shot-number"],
			shotNumber,
		)
	}

	remaining := time.Until(expiresAt)

	if remaining > 16*time.Minute ||
		remaining < 14*time.Minute {
		t.Errorf(
			"unexpected expiration time: %v",
			expiresAt,
		)
	}
}

// 異常系
func TestService_GeneratePresignURL_Error(t *testing.T) {
	t.Parallel()

	mock := &mockPresigner{
		mockErr: errors.New("aws error"),
	}

	svc := NewService(mock)

	_, _, err := svc.GeneratePresignURL(
		context.Background(),
		"001",
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
