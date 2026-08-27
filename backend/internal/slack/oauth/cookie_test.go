package oauth

import (
	"net/http"
	"testing"
	"time"
)

const (
	expectedOAuthStateCookieName = "oauth_state"
	expectedOAuthStateCookiePath = "/api/oauth/slack"
)

func TestBuildStateCookie(t *testing.T) {
	t.Parallel()

	const testState = "test-random-state-string"

	before := time.Now()

	cookie := BuildStateCookie(testState, true)

	if cookie.Name != expectedOAuthStateCookieName {
		t.Errorf(
			"expected cookie name %q, got %q",
			expectedOAuthStateCookieName,
			cookie.Name,
		)
	}

	if cookie.Value != testState {
		t.Errorf(
			"expected cookie value %q, got %q",
			testState,
			cookie.Value,
		)
	}

	if cookie.Path != expectedOAuthStateCookiePath {
		t.Errorf(
			"expected Path %q, got %q",
			expectedOAuthStateCookiePath,
			cookie.Path,
		)
	}

	if !cookie.HttpOnly {
		t.Error("expected HttpOnly to be true")
	}

	if !cookie.Secure {
		t.Error("expected Secure to be true")
	}

	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf(
			"expected SameSite LaxMode, got %v",
			cookie.SameSite,
		)
	}

	expectedMaxAge := int(OAuthStateTTL.Seconds())
	if cookie.MaxAge != expectedMaxAge {
		t.Errorf(
			"expected MaxAge %d, got %d",
			expectedMaxAge,
			cookie.MaxAge,
		)
	}

	minExpires := before.Add(OAuthStateTTL)
	maxExpires := time.Now().Add(OAuthStateTTL)

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

func TestBuildStateCookie_NotSecure(t *testing.T) {
	t.Parallel()

	cookie := BuildStateCookie("test-state", false)

	if cookie.Secure {
		t.Error("expected Secure to be false")
	}
}

func TestBuildDeleteStateCookie(t *testing.T) {
	t.Parallel()

	cookie := BuildDeleteStateCookie(true)

	if cookie.Name != expectedOAuthStateCookieName {
		t.Errorf(
			"expected cookie name %q, got %q",
			expectedOAuthStateCookieName,
			cookie.Name,
		)
	}

	if cookie.Value != "" {
		t.Errorf(
			"expected empty cookie value, got %q",
			cookie.Value,
		)
	}

	if cookie.Path != expectedOAuthStateCookiePath {
		t.Errorf(
			"expected Path %q, got %q",
			expectedOAuthStateCookiePath,
			cookie.Path,
		)
	}

	if !cookie.HttpOnly {
		t.Error("expected HttpOnly to be true")
	}

	if !cookie.Secure {
		t.Error("expected Secure to be true")
	}

	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf(
			"expected SameSite LaxMode, got %v",
			cookie.SameSite,
		)
	}

	if cookie.MaxAge != -1 {
		t.Errorf(
			"expected MaxAge -1, got %d",
			cookie.MaxAge,
		)
	}

	if !cookie.Expires.Equal(time.Unix(0, 0)) {
		t.Errorf(
			"expected Unix epoch expiration, got %v",
			cookie.Expires,
		)
	}
}
