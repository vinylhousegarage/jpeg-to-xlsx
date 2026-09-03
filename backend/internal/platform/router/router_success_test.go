package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRoutes(
	t *testing.T,
) {
	tests := []struct {
		name     string
		method   string
		path     string
		response string
	}{
		{
			name:     "authentication login",
			method:   http.MethodGet,
			path:     "/api/auth/login",
			response: testAuthLoginResponse,
		},
		{
			name:     "authentication callback",
			method:   http.MethodGet,
			path:     "/api/auth/callback",
			response: testAuthCallbackResponse,
		},
		{
			name:     "authentication session",
			method:   http.MethodGet,
			path:     "/api/auth/session",
			response: testAuthSessionResponse,
		},
		{
			name:     "authentication logout",
			method:   http.MethodPost,
			path:     "/api/auth/logout",
			response: testAuthLogoutResponse,
		},
		{
			name:     "Slack login",
			method:   http.MethodGet,
			path:     "/api/oauth/slack/login",
			response: testSlackLoginResponse,
		},
		{
			name:     "Slack callback",
			method:   http.MethodGet,
			path:     "/api/oauth/slack/callback",
			response: testSlackCallbackResponse,
		},
		{
			name:     "storage upload",
			method:   http.MethodPost,
			path:     "/api/storage/upload",
			response: testStorageResponse,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				mux := newTestRouter()

				request := httptest.NewRequest(
					test.method,
					test.path,
					nil,
				)

				response := httptest.NewRecorder()

				mux.ServeHTTP(
					response,
					request,
				)

				if response.Code !=
					http.StatusOK {
					t.Fatalf(
						"response status = %d, want %d",
						response.Code,
						http.StatusOK,
					)
				}

				if response.Body.String() !=
					test.response {
					t.Errorf(
						"response body = %q, want %q",
						response.Body.String(),
						test.response,
					)
				}
			},
		)
	}
}
