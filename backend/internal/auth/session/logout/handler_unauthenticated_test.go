package logout

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
)

func TestHandler_ServeHTTP_Unauthenticated(
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
					http.StatusNoContent {
					t.Errorf(
						"status code = %d, want %d",
						response.Code,
						http.StatusNoContent,
					)
				}

				if response.Body.Len() != 0 {
					t.Errorf(
						"response body = %q, want empty",
						response.Body.String(),
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

				if dependencies.store.deleteCalls !=
					0 {
					t.Errorf(
						"Delete() store calls = %d, want 0",
						dependencies.store.deleteCalls,
					)
				}

				if dependencies.cookieDeleter.deleteCalls !=
					1 {
					t.Errorf(
						"Delete() cookie calls = %d, want 1",
						dependencies.cookieDeleter.deleteCalls,
					)
				}

				if dependencies.cookieDeleter.receivedWriter !=
					response {
					t.Error(
						"Delete() cookie writer does not match",
					)
				}
			},
		)
	}
}
