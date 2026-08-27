package put

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/storage"
)

type mockPresignerForHandler struct {
	mockResult *v4.PresignedHTTPRequest
	mockErr    error
}

func (m *mockPresignerForHandler) PresignPutObject(
	ctx context.Context,
	params *s3.PutObjectInput,
	optFns ...func(*s3.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {
	return m.mockResult, m.mockErr
}

func TestHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	tests := []struct {
		name           string
		method         string
		requestBody    any
		mockResult     *v4.PresignedHTTPRequest
		mockErr        error
		expectedStatus int
	}{
		{
			name:   "Success: returns presigned URL with valid request",
			method: http.MethodPost,
			requestBody: storage.PutPresignRequest{
				ShotNumber: "test-shot",
			},
			mockResult: &v4.PresignedHTTPRequest{
				URL: "https://example.com/presigned-put-url",
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: returns method not allowed for GET request",
			method:         http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Error: returns bad request for invalid JSON format",
			method:         http.MethodPost,
			requestBody:    "invalid-json-string",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error: returns bad request when shot number is empty",
			method: http.MethodPost,
			requestBody: storage.PutPresignRequest{
				ShotNumber: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error: returns internal server error when AWS service fails",
			method: http.MethodPost,
			requestBody: storage.PutPresignRequest{
				ShotNumber: "test-shot",
			},
			mockResult:     nil,
			mockErr:        errors.New("aws internal error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var bodyReader *bytes.Reader
			if tt.requestBody != nil {
				if str, ok := tt.requestBody.(string); ok {
					bodyReader = bytes.NewReader([]byte(str))
				} else {
					jsonBytes, _ := json.Marshal(tt.requestBody)
					bodyReader = bytes.NewReader(jsonBytes)
				}
			} else {
				bodyReader = bytes.NewReader([]byte{})
			}

			req := httptest.NewRequest(tt.method, "/presign/put", bodyReader)
			rec := httptest.NewRecorder()

			mockPresigner := &mockPresignerForHandler{
				mockResult: tt.mockResult,
				mockErr:    tt.mockErr,
			}
			svc := NewService(mockPresigner)
			handler := NewHandler(svc, logger)

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d (response body: %s)", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}
