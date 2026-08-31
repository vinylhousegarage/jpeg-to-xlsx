package session

import (
	"net/http"
	"strings"
	"testing"
)

func TestCookieManagerWrite(
	t *testing.T,
) {
	t.Parallel()

	manager := newTestCookieManager(t)
	response := newTestCookieResponseRecorder()

	err := manager.Write(
		response,
		testCookieRawSessionID,
	)
	if err != nil {
		t.Fatalf(
			"Write() error = %v",
			err,
		)
	}

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

	if cookie.Value != testCookieRawSessionID {
		t.Errorf(
			"cookie value = %q, want %q",
			cookie.Value,
			testCookieRawSessionID,
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

	wantMaxAge := int(
		testCookieLifetime.Seconds(),
	)

	if cookie.MaxAge != wantMaxAge {
		t.Errorf(
			"cookie MaxAge = %d, want %d",
			cookie.MaxAge,
			wantMaxAge,
		)
	}

	wantExpires := testCookieNow().
		Add(testCookieLifetime)

	if !cookie.Expires.Equal(wantExpires) {
		t.Errorf(
			"cookie Expires = %v, want %v",
			cookie.Expires,
			wantExpires,
		)
	}

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

func TestCookieManagerWriteRejectsInvalidSessionID(
	t *testing.T,
) {
	t.Parallel()

	for _, value := range whitespaceCookieValues() {
		value := value

		t.Run(
			"session ID "+value,
			func(t *testing.T) {
				t.Parallel()

				manager :=
					newTestCookieManager(t)
				response :=
					newTestCookieResponseRecorder()

				err := manager.Write(
					response,
					value,
				)
				if err == nil {
					t.Fatal(
						"Write() error = nil, want an error",
					)
				}

				if got := response.Header().
					Get("Set-Cookie"); got != "" {
					t.Errorf(
						"Set-Cookie header = %q, want empty",
						got,
					)
				}
			},
		)
	}
}
