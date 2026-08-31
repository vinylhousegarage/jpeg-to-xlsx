package session

import (
	"net/http"
	"strings"
	"testing"
)

func TestCookieManagerDelete(
	t *testing.T,
) {
	t.Parallel()

	manager := newTestCookieManager(t)
	response := newTestCookieResponseRecorder()

	manager.Delete(response)

	cookie := responseCookie(
		t,
		response,
	)

	if cookie.Name != sessionCookieName {
		t.Errorf(
			"cookie name = %q, want %q",
			cookie.Name,
			sessionCookieName,
		)
	}

	if cookie.Value != "" {
		t.Errorf(
			"cookie value = %q, want empty",
			cookie.Value,
		)
	}

	if cookie.Path != "/" {
		t.Errorf(
			"cookie path = %q, want %q",
			cookie.Path,
			"/",
		)
	}

	if cookie.Domain != "" {
		t.Errorf(
			"cookie domain = %q, want empty",
			cookie.Domain,
		)
	}

	if !cookie.Secure {
		t.Error(
			"cookie Secure = false, want true",
		)
	}

	if !cookie.HttpOnly {
		t.Error(
			"cookie HttpOnly = false, want true",
		)
	}

	if cookie.SameSite !=
		http.SameSiteLaxMode {
		t.Errorf(
			"cookie SameSite = %v, want %v",
			cookie.SameSite,
			http.SameSiteLaxMode,
		)
	}

	if cookie.MaxAge >= 0 {
		t.Errorf(
			"cookie MaxAge = %d, want a negative value",
			cookie.MaxAge,
		)
	}

	if !cookie.Expires.Before(
		testCookieNow(),
	) {
		t.Errorf(
			"cookie Expires = %v, want before %v",
			cookie.Expires,
			testCookieNow(),
		)
	}
}

func TestCookieManagerDeleteDoesNotSetDomain(
	t *testing.T,
) {
	t.Parallel()

	manager := newTestCookieManager(t)
	response := newTestCookieResponseRecorder()

	manager.Delete(response)

	setCookieHeader := response.Header().
		Get("Set-Cookie")

	if setCookieHeader == "" {
		t.Fatal(
			"Set-Cookie header is empty",
		)
	}

	if strings.Contains(
		strings.ToLower(setCookieHeader),
		"domain=",
	) {
		t.Errorf(
			"Set-Cookie header contains Domain: %q",
			setCookieHeader,
		)
	}
}
