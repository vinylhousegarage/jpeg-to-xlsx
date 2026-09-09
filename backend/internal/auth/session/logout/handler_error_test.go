package logout

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ServeHTTP_InternalError(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name               string
		configure          func(*testDependencies)
		errorMessage       string
		wantLogoutURLCalls int
		wantResolveCalls   int
		wantStoreCalls     int
		wantCookieDeletes  int
	}{
		{
			name: "logout URL provider error",
			configure: func(dependencies *testDependencies) {
				dependencies.logoutURLProvider.err = errors.New("create logout URL failed")
			},
			errorMessage:       "create logout URL failed",
			wantLogoutURLCalls: 1,
			wantResolveCalls:   0,
			wantStoreCalls:     0,
			wantCookieDeletes:  0,
		},
		{
			name: "session resolver error",
			configure: func(dependencies *testDependencies) {
				dependencies.resolver.err = errors.New("resolve session failed")
			},
			errorMessage:       "resolve session failed",
			wantLogoutURLCalls: 1,
			wantResolveCalls:   1,
			wantStoreCalls:     0,
			wantCookieDeletes:  0,
		},
		{
			name: "session store delete error",
			configure: func(dependencies *testDependencies) {
				dependencies.store.deleteErr = errors.New("delete session failed")
			},
			errorMessage:       "delete session failed",
			wantLogoutURLCalls: 1,
			wantResolveCalls:   1,
			wantStoreCalls:     1,
			wantCookieDeletes:  0,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				dependencies := newTestDependencies()
				test.configure(dependencies)
				handler := newTestHandler(t, dependencies)
				request := newTestRequest(t)
				response := httptest.NewRecorder()
				handler.ServeHTTP(response, request)

				if response.Code != http.StatusInternalServerError {
					t.Errorf(
						"status code = %d, want %d",
						response.Code,
						http.StatusInternalServerError,
					)
				}

				if got := response.Header().Get("Cache-Control"); got != "no-store" {
					t.Errorf(
						"Cache-Control = %q, want %q",
						got,
						"no-store",
					)
				}

				if strings.Contains(response.Body.String(), test.errorMessage) {
					t.Error("response body exposes internal error")
				}

				if dependencies.logoutURLProvider.calls != test.wantLogoutURLCalls {
					t.Errorf(
						"LogoutURL() calls = %d, want %d",
						dependencies.logoutURLProvider.calls,
						test.wantLogoutURLCalls,
					)
				}

				if dependencies.resolver.resolveCalls != test.wantResolveCalls {
					t.Errorf(
						"Resolve() calls = %d, want %d",
						dependencies.resolver.resolveCalls,
						test.wantResolveCalls,
					)
				}

				if dependencies.store.deleteCalls != test.wantStoreCalls {
					t.Errorf(
						"Delete() store calls = %d, want %d",
						dependencies.store.deleteCalls,
						test.wantStoreCalls,
					)
				}

				if dependencies.cookieDeleter.deleteCalls != test.wantCookieDeletes {
					t.Errorf(
						"Delete() cookie calls = %d, want %d",
						dependencies.cookieDeleter.deleteCalls,
						test.wantCookieDeletes,
					)
				}
			},
		)
	}
}
