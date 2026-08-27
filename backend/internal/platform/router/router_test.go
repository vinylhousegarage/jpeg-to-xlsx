package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRoutes(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		response string
	}{
		{
			name:     "routes storage upload to presign handler",
			path:     "/api/storage/upload",
			response: "presign",
		},
		{
			name:     "routes slack login to login handler",
			path:     "/api/oauth/slack/login",
			response: "slack login",
		},
		{
			name:     "routes slack callback to callback handler",
			path:     "/api/oauth/slack/callback",
			response: "slack callback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()

			presignHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("presign"))
			})

			slackLoginHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("slack login"))
			})

			slackCallbackHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("slack callback"))
			})

			SetupRoutes(
				mux,
				presignHandler,
				slackLoginHandler,
				slackCallbackHandler,
			)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
			}

			if rec.Body.String() != tt.response {
				t.Errorf("expected body %q, got %q", tt.response, rec.Body.String())
			}
		})
	}
}
