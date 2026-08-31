package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	testCookieRawSessionID = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQ"

	testCookieRequestURL = "https://example.com/api/auth/session"

	testCookieLifetime = 8 * time.Hour
)

func newTestCookieManager(
	t *testing.T,
) *CookieManager {
	t.Helper()

	manager, err := NewCookieManager(
		testCookieLifetime,
	)
	if err != nil {
		t.Fatalf(
			"NewCookieManager() error = %v",
			err,
		)
	}

	manager.now = testCookieNow

	return manager
}

func newTestCookieRequest() *http.Request {
	return httptest.NewRequest(
		http.MethodGet,
		testCookieRequestURL,
		nil,
	)
}

func newTestCookieRequestWithSession(
	rawSessionID string,
) *http.Request {
	request := newTestCookieRequest()

	request.AddCookie(
		&http.Cookie{
			Name:  sessionCookieName,
			Value: rawSessionID,
		},
	)

	return request
}

func newTestCookieResponseRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func responseCookie(
	t *testing.T,
	response *httptest.ResponseRecorder,
) *http.Cookie {
	t.Helper()

	result := response.Result()

	defer func() {
		if err := result.Body.Close(); err != nil {
			t.Errorf(
				"close response body: %v",
				err,
			)
		}
	}()

	cookies := result.Cookies()

	if len(cookies) != 1 {
		t.Fatalf(
			"response cookie count = %d, want 1",
			len(cookies),
		)
	}

	return cookies[0]
}

func testCookieNow() time.Time {
	return time.Date(
		2026,
		time.August,
		31,
		12,
		34,
		56,
		0,
		time.UTC,
	)
}

func whitespaceCookieValues() []string {
	return []string{
		"",
		" ",
		"   ",
		"\t",
		"\n",
		strings.Repeat(
			" ",
			10,
		),
	}
}
