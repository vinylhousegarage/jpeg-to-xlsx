package session

import (
	"errors"
	"net/http"
	"testing"
)

func TestCookieManagerRead(
	t *testing.T,
) {
	t.Parallel()

	manager := newTestCookieManager(t)
	request := newTestCookieRequestWithSession(
		testCookieRawSessionID,
	)

	rawSessionID, err := manager.Read(request)
	if err != nil {
		t.Fatalf(
			"Read() error = %v",
			err,
		)
	}

	if rawSessionID != testCookieRawSessionID {
		t.Errorf(
			"Read() session ID = %q, want %q",
			rawSessionID,
			testCookieRawSessionID,
		)
	}
}

func TestCookieManagerReadIgnoresOtherCookies(
	t *testing.T,
) {
	t.Parallel()

	manager := newTestCookieManager(t)
	request := newTestCookieRequest()

	request.AddCookie(
		&http.Cookie{
			Name:  "other-cookie",
			Value: "other-value",
		},
	)
	request.AddCookie(
		&http.Cookie{
			Name:  sessionCookieName,
			Value: testCookieRawSessionID,
		},
	)

	rawSessionID, err := manager.Read(request)
	if err != nil {
		t.Fatalf(
			"Read() error = %v",
			err,
		)
	}

	if rawSessionID != testCookieRawSessionID {
		t.Errorf(
			"Read() session ID = %q, want %q",
			rawSessionID,
			testCookieRawSessionID,
		)
	}
}

func TestCookieManagerReadReturnsNotFoundWhenCookieIsMissing(
	t *testing.T,
) {
	t.Parallel()

	manager := newTestCookieManager(t)
	request := newTestCookieRequest()

	rawSessionID, err := manager.Read(request)
	if err == nil {
		t.Fatal(
			"Read() error = nil, want an error",
		)
	}

	if rawSessionID != "" {
		t.Errorf(
			"Read() session ID = %q, want empty",
			rawSessionID,
		)
	}

	if !errors.Is(err, http.ErrNoCookie) {
		t.Errorf(
			"Read() error = %v, want wrapped %v",
			err,
			http.ErrNoCookie,
		)
	}
}

func TestCookieManagerReadRejectsEmptySessionID(
	t *testing.T,
) {
	t.Parallel()

	manager := newTestCookieManager(t)
	request := newTestCookieRequest()

	request.AddCookie(
		&http.Cookie{
			Name:  sessionCookieName,
			Value: "",
		},
	)

	rawSessionID, err := manager.Read(request)
	if err == nil {
		t.Fatal(
			"Read() error = nil, want an error",
		)
	}

	if rawSessionID != "" {
		t.Errorf(
			"Read() session ID = %q, want empty",
			rawSessionID,
		)
	}
}

func TestCookieManagerReadRejectsWhitespaceSessionID(
	t *testing.T,
) {
	t.Parallel()

	manager := newTestCookieManager(t)
	request := newTestCookieRequest()

	request.Header.Set(
		"Cookie",
		sessionCookieName+`="   "`,
	)

	rawSessionID, err := manager.Read(request)
	if err == nil {
		t.Fatal(
			"Read() error = nil, want an error",
		)
	}

	if rawSessionID != "" {
		t.Errorf(
			"Read() session ID = %q, want empty",
			rawSessionID,
		)
	}
}

func TestCookieManagerReadRejectsNilRequest(
	t *testing.T,
) {
	t.Parallel()

	manager := newTestCookieManager(t)

	rawSessionID, err := manager.Read(nil)
	if err == nil {
		t.Fatal(
			"Read() error = nil, want an error",
		)
	}

	if rawSessionID != "" {
		t.Errorf(
			"Read() session ID = %q, want empty",
			rawSessionID,
		)
	}
}
