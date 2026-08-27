package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestWriteError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantCode    ErrorCode
		wantLevel   zapcore.Level
		wantMessage string
	}{
		{
			name: "APIError returns specified status and warn log",
			err: New(
				ErrorCodeMissingShotNumber,
				http.StatusBadRequest,
				nil,
			),
			wantStatus:  http.StatusBadRequest,
			wantCode:    ErrorCodeMissingShotNumber,
			wantLevel:   zapcore.WarnLevel,
			wantMessage: "request rejected",
		},
		{
			name:        "unknown error returns 500 and error log",
			err:         errors.New("unknown database error"),
			wantStatus:  http.StatusInternalServerError,
			wantCode:    ErrorCodeInternal,
			wantLevel:   zapcore.ErrorLevel,
			wantMessage: "request failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			core, observedLogs := observer.New(zapcore.DebugLevel)
			logger := zap.New(core)
			w := httptest.NewRecorder()

			WriteError(w, tt.err, logger)

			if w.Code != tt.wantStatus {
				t.Errorf(
					"status = %d, want %d",
					w.Code,
					tt.wantStatus,
				)
			}

			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf(
					"Content-Type = %q, want %q",
					contentType,
					"application/json",
				)
			}

			var response ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if response.Error != string(tt.wantCode) {
				t.Errorf(
					"error code = %q, want %q",
					response.Error,
					tt.wantCode,
				)
			}

			logs := observedLogs.All()
			if len(logs) != 1 {
				t.Fatalf("log count = %d, want 1", len(logs))
			}

			entry := logs[0]

			if entry.Level != tt.wantLevel {
				t.Errorf(
					"log level = %v, want %v",
					entry.Level,
					tt.wantLevel,
				)
			}

			if entry.Message != tt.wantMessage {
				t.Errorf(
					"log message = %q, want %q",
					entry.Message,
					tt.wantMessage,
				)
			}

			fields := entry.ContextMap()

			if fields["error_code"] != string(tt.wantCode) {
				t.Errorf(
					"log error_code = %v, want %q",
					fields["error_code"],
					tt.wantCode,
				)
			}

			if fields["http_status"] != int64(tt.wantStatus) {
				t.Errorf(
					"log http_status = %v, want %d",
					fields["http_status"],
					tt.wantStatus,
				)
			}

			if fields["error"] == nil {
				t.Error("expected error field in log")
			}
		})
	}
}
