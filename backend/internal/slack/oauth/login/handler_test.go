package login

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/apierror"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

const (
	testClientID    = "test-client-id"
	testRedirectURI = "https://example.com/oauth/slack/callback"
)

func TestHandler_ServeHTTP_RedirectsToSlack(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(
		http.MethodGet,
		"/oauth/slack/login",
		nil,
	)
	rec := httptest.NewRecorder()

	before := time.Now()

	handler := NewHandler(
		testClientID,
		testRedirectURI,
		true,
		zap.NewNop(),
	)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"expected status %d, got %d (response body: %s)",
			http.StatusFound,
			rec.Code,
			rec.Body.String(),
		)
	}

	location := rec.Header().Get("Location")
	if location == "" {
		t.Fatal("expected Location header to be set")
	}

	authURL, err := url.Parse(location)
	if err != nil {
		t.Fatalf("failed to parse Location header: %v", err)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one response cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "oauth_state" {
		t.Fatalf(
			"expected cookie name %q, got %q",
			"oauth_state",
			cookie.Name,
		)
	}

	if cookie.Value == "" {
		t.Fatal("expected a non-empty OAuth state cookie")
	}

	redirectState := authURL.Query().Get("state")
	if redirectState == "" {
		t.Fatal("expected state query parameter in redirect URL")
	}

	if redirectState != cookie.Value {
		t.Errorf(
			"expected redirect state to match cookie value: redirect=%q cookie=%q",
			redirectState,
			cookie.Value,
		)
	}

	minExpires := before.Add(oauth.OAuthStateTTL - time.Second)
	maxExpires := time.Now().Add(oauth.OAuthStateTTL + time.Second)

	if cookie.Expires.Before(minExpires) ||
		cookie.Expires.After(maxExpires) {
		t.Errorf(
			"expected Expires between %v and %v, got %v",
			minExpires,
			maxExpires,
			cookie.Expires,
		)
	}
}

func TestHandler_ServeHTTP_RejectsUnsupportedMethods(t *testing.T) {
	t.Parallel()

	methods := []string{
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}

	for _, method := range methods {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(
				method,
				"/oauth/slack/login",
				nil,
			)
			rec := httptest.NewRecorder()

			handler := NewHandler(
				testClientID,
				testRedirectURI,
				true,
				zap.NewNop(),
			)
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf(
					"expected status %d, got %d (response body: %s)",
					http.StatusMethodNotAllowed,
					rec.Code,
					rec.Body.String(),
				)
			}

			if location := rec.Header().Get("Location"); location != "" {
				t.Errorf("expected no redirect, got %q", location)
			}

			if cookies := rec.Result().Cookies(); len(cookies) != 0 {
				t.Errorf(
					"expected no response cookies, got %d",
					len(cookies),
				)
			}

			var response apierror.ErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}

			if response.Error != string(apierror.ErrorCodeInvalidMethod) {
				t.Errorf(
					"expected error code %q, got %q",
					apierror.ErrorCodeInvalidMethod,
					response.Error,
				)
			}
		})
	}
}
