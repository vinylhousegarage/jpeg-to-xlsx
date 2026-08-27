package callback

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/apierror"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

const (
	testState        = "test-state"
	testCode         = "test-code"
	testCallbackPath = "/api/oauth/slack/callback"
	testRedirectURI  = "https://example.com/api/oauth/slack/callback"

	expectedStateCookiePath = "/api/oauth/slack"
)

type stubCodeExchanger struct {
	token *oauth.Token
	err   error

	called      bool
	gotCode     string
	gotRedirect string
}

func (s *stubCodeExchanger) ExchangeCode(
	_ context.Context,
	code string,
	redirectURI string,
) (*oauth.Token, error) {
	s.called = true
	s.gotCode = code
	s.gotRedirect = redirectURI

	return s.token, s.err
}

type stubConversationOpener struct {
	channelID string
	err       error

	called         bool
	gotAccessToken string
	gotUserID      string
}

func (s *stubConversationOpener) OpenConversation(
	_ context.Context,
	accessToken string,
	userID string,
) (string, error) {
	s.called = true
	s.gotAccessToken = accessToken
	s.gotUserID = userID

	return s.channelID, s.err
}

type stubTokenStore struct {
	err error

	called   bool
	gotToken *oauth.Token
}

func (s *stubTokenStore) Save(
	_ context.Context,
	token *oauth.Token,
) error {
	s.called = true
	s.gotToken = token

	return s.err
}

func newValidCallbackRequest() *http.Request {
	req := httptest.NewRequest(
		http.MethodGet,
		testCallbackPath+
			"?state="+testState+
			"&code="+testCode,
		nil,
	)

	req.AddCookie(
		&http.Cookie{
			Name:  oauthStateCookieName,
			Value: testState,
		},
	)

	return req
}

func assertErrorResponse(
	t *testing.T,
	rec *httptest.ResponseRecorder,
	wantStatus int,
	wantCode apierror.ErrorCode,
) {
	t.Helper()

	if rec.Code != wantStatus {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			rec.Code,
			wantStatus,
			rec.Body.String(),
		)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf(
			"Content-Type = %q, want %q",
			got,
			"application/json",
		)
	}

	var response apierror.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf(
			"failed to decode error response: %v",
			err,
		)
	}

	if response.Error != string(wantCode) {
		t.Errorf(
			"error code = %q, want %q",
			response.Error,
			wantCode,
		)
	}
}

func assertDeleteStateCookie(
	t *testing.T,
	rec *httptest.ResponseRecorder,
	wantSecure bool,
) {
	t.Helper()

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf(
			"response cookie count = %d, want 1",
			len(cookies),
		)
	}

	cookie := cookies[0]

	if cookie.Name != oauthStateCookieName {
		t.Errorf(
			"cookie name = %q, want %q",
			cookie.Name,
			oauthStateCookieName,
		)
	}

	if cookie.Value != "" {
		t.Errorf(
			"cookie value = %q, want empty",
			cookie.Value,
		)
	}

	if cookie.Path != expectedStateCookiePath {
		t.Errorf(
			"cookie Path = %q, want %q",
			cookie.Path,
			expectedStateCookiePath,
		)
	}

	if !cookie.HttpOnly {
		t.Error("cookie HttpOnly = false, want true")
	}

	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf(
			"cookie SameSite = %v, want %v",
			cookie.SameSite,
			http.SameSiteLaxMode,
		)
	}

	if cookie.MaxAge != -1 {
		t.Errorf(
			"cookie MaxAge = %d, want -1",
			cookie.MaxAge,
		)
	}

	if cookie.Secure != wantSecure {
		t.Errorf(
			"cookie Secure = %t, want %t",
			cookie.Secure,
			wantSecure,
		)
	}
}
