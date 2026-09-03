package status

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
)

func TestHandler_ServeHTTP_Unauthorized(
	t *testing.T,
) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "unauthenticated session",
			err:  authsession.ErrUnauthenticated,
		},
		{
			name: "wrapped unauthenticated session",
			err: fmt.Errorf(
				"resolve session: %w",
				authsession.ErrUnauthenticated,
			),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				dependencies :=
					newTestDependencies()

				dependencies.resolver.err =
					test.err

				handler := newTestHandler(
					t,
					dependencies,
				)

				request :=
					newTestRequest(t)

				response :=
					httptest.NewRecorder()

				handler.ServeHTTP(
					response,
					request,
				)

				if response.Code !=
					http.StatusUnauthorized {
					t.Errorf(
						"status code = %d, want %d",
						response.Code,
						http.StatusUnauthorized,
					)
				}

				if got := response.Header().
					Get("Content-Type"); got !=
					"application/json" {
					t.Errorf(
						"Content-Type = %q, want %q",
						got,
						"application/json",
					)
				}

				if got := response.Header().
					Get("Cache-Control"); got !=
					"no-store" {
					t.Errorf(
						"Cache-Control = %q, want %q",
						got,
						"no-store",
					)
				}

				var body struct {
					Authenticated bool `json:"authenticated"`
				}

				if err := json.NewDecoder(
					response.Body,
				).Decode(&body); err != nil {
					t.Fatalf(
						"decode response body: %v",
						err,
					)
				}

				if body.Authenticated {
					t.Errorf(
						"authenticated = true, want false",
					)
				}

				if dependencies.resolver.resolveCalls !=
					1 {
					t.Errorf(
						"Resolve() calls = %d, want 1",
						dependencies.resolver.resolveCalls,
					)
				}

				if dependencies.resolver.receivedCtx ==
					nil {
					t.Error(
						"Resolve() context is nil",
					)
				}

				if dependencies.resolver.receivedRequest !=
					request {
					t.Error(
						"Resolve() request does not match",
					)
				}
			},
		)
	}
}
