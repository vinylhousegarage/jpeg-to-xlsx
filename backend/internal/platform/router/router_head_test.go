package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRoutesRejectsHEAD(
	t *testing.T,
) {
	tests := []struct {
		name          string
		path          string
		allowedMethod string
	}{
		{
			name:          "authentication login",
			path:          "/api/auth/login",
			allowedMethod: http.MethodGet,
		},
		{
			name:          "authentication callback",
			path:          "/api/auth/callback",
			allowedMethod: http.MethodGet,
		},
		{
			name:          "Slack login",
			path:          "/api/oauth/slack/login",
			allowedMethod: http.MethodGet,
		},
		{
			name:          "Slack callback",
			path:          "/api/oauth/slack/callback",
			allowedMethod: http.MethodGet,
		},
		{
			name:          "storage upload",
			path:          "/api/storage/upload",
			allowedMethod: http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				mux := newTestRouter()

				request := httptest.NewRequest(
					http.MethodHead,
					test.path,
					nil,
				)

				response := httptest.NewRecorder()

				mux.ServeHTTP(
					response,
					request,
				)

				if response.Code !=
					http.StatusMethodNotAllowed {
					t.Errorf(
						"response status = %d, want %d",
						response.Code,
						http.StatusMethodNotAllowed,
					)
				}

				if got := response.Header().Get(
					"Allow",
				); got != test.allowedMethod {
					t.Errorf(
						"Allow header = %q, want %q",
						got,
						test.allowedMethod,
					)
				}

				if response.Body.Len() != 0 {
					t.Errorf(
						"response body = %q, want empty",
						response.Body.String(),
					)
				}
			},
		)
	}
}
