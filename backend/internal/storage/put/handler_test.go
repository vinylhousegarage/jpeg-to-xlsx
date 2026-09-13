package put

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"

	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/storage"
)

const testHandlerCognitoSub = "test-cognito-sub"

type mockSessionResolverForHandler struct {
	session authsession.Session
	err     error
	called  bool
}

func (m *mockSessionResolverForHandler) Resolve(
	_ context.Context,
	_ *http.Request,
) (
	authsession.Session,
	error,
) {
	m.called = true

	return m.session, m.err
}

type mockPresignerForHandler struct {
	mockResult *v4.PresignedHTTPRequest
	mockErr    error
	called     bool
	gotInput   *s3.PutObjectInput
}

func (m *mockPresignerForHandler) PresignPutObject(
	_ context.Context,
	params *s3.PutObjectInput,
	_ ...func(*s3.PresignOptions),
) (
	*v4.PresignedHTTPRequest,
	error,
) {
	m.called = true
	m.gotInput = params

	return m.mockResult, m.mockErr
}

func TestHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	tests := []struct {
		name                  string
		method                string
		requestBody           any
		resolverSession       authsession.Session
		resolverErr           error
		mockResult            *v4.PresignedHTTPRequest
		mockErr               error
		expectedStatus        int
		expectedObjectKey     string
		expectResolverCalled  bool
		expectPresignerCalled bool
	}{
		{
			name:   "Success: returns presigned URL with valid request",
			method: http.MethodPost,
			requestBody: storage.PutPresignRequest{
				ShotNumber: "test-shot",
			},
			resolverSession: authsession.Session{
				CognitoSub: testHandlerCognitoSub,
			},
			mockResult: &v4.PresignedHTTPRequest{
				URL: "https://example.com/presigned-put-url",
			},
			expectedStatus:        http.StatusOK,
			expectedObjectKey:     "test-cognito-sub/test-shot.jpg",
			expectResolverCalled:  true,
			expectPresignerCalled: true,
		},
		{
			name:                  "Error: returns method not allowed for GET request",
			method:                http.MethodGet,
			requestBody:           nil,
			expectedStatus:        http.StatusMethodNotAllowed,
			expectResolverCalled:  false,
			expectPresignerCalled: false,
		},
		{
			name:                  "Error: returns bad request for invalid JSON format",
			method:                http.MethodPost,
			requestBody:           "invalid-json-string",
			expectedStatus:        http.StatusBadRequest,
			expectResolverCalled:  false,
			expectPresignerCalled: false,
		},
		{
			name:   "Error: returns bad request when shot number is empty",
			method: http.MethodPost,
			requestBody: storage.PutPresignRequest{
				ShotNumber: "",
			},
			expectedStatus:        http.StatusBadRequest,
			expectResolverCalled:  false,
			expectPresignerCalled: false,
		},
		{
			name:   "Error: returns unauthorized when session is unauthenticated",
			method: http.MethodPost,
			requestBody: storage.PutPresignRequest{
				ShotNumber: "test-shot",
			},
			resolverErr:           authsession.ErrUnauthenticated,
			expectedStatus:        http.StatusUnauthorized,
			expectResolverCalled:  true,
			expectPresignerCalled: false,
		},
		{
			name:   "Error: returns internal server error when session resolution fails",
			method: http.MethodPost,
			requestBody: storage.PutPresignRequest{
				ShotNumber: "test-shot",
			},
			resolverErr:           errors.New("dynamodb unavailable"),
			expectedStatus:        http.StatusInternalServerError,
			expectResolverCalled:  true,
			expectPresignerCalled: false,
		},
		{
			name:   "Error: returns internal server error when AWS service fails",
			method: http.MethodPost,
			requestBody: storage.PutPresignRequest{
				ShotNumber: "test-shot",
			},
			resolverSession: authsession.Session{
				CognitoSub: testHandlerCognitoSub,
			},
			mockErr:               errors.New("aws internal error"),
			expectedStatus:        http.StatusInternalServerError,
			expectResolverCalled:  true,
			expectPresignerCalled: true,
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
					jsonBytes, err := json.Marshal(tt.requestBody)
					if err != nil {
						t.Fatalf("failed to marshal request body: %v", err)
					}

					bodyReader = bytes.NewReader(jsonBytes)
				}
			} else {
				bodyReader = bytes.NewReader(nil)
			}

			req := httptest.NewRequest(tt.method, "/presign/put", bodyReader)
			rec := httptest.NewRecorder()

			resolver := &mockSessionResolverForHandler{
				session: tt.resolverSession,
				err:     tt.resolverErr,
			}

			mockPresigner := &mockPresignerForHandler{
				mockResult: tt.mockResult,
				mockErr:    tt.mockErr,
			}

			svc := NewService(mockPresigner)
			handler := NewHandler(svc, resolver, logger)
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf(
					"expected status %d, got %d (response body: %s)",
					tt.expectedStatus,
					rec.Code,
					rec.Body.String(),
				)
			}

			if resolver.called != tt.expectResolverCalled {
				t.Errorf(
					"Resolve() called = %t, want %t",
					resolver.called,
					tt.expectResolverCalled,
				)
			}

			if mockPresigner.called != tt.expectPresignerCalled {
				t.Errorf(
					"PresignPutObject() called = %t, want %t",
					mockPresigner.called,
					tt.expectPresignerCalled,
				)
			}

			if tt.expectedObjectKey != "" {
				if mockPresigner.gotInput == nil {
					t.Fatal("PresignPutObject() input is nil")
				}

				gotObjectKey := aws.ToString(mockPresigner.gotInput.Key)

				if gotObjectKey != tt.expectedObjectKey {
					t.Errorf(
						"object key = %q, want %q",
						gotObjectKey,
						tt.expectedObjectKey,
					)
				}
			}
		})
	}
}
