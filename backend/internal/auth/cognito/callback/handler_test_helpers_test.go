package callback

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
)

const (
	testAuthorizationCode = "test-authorization-code"

	testOAuthStateValue = "test-oauth-state"

	testCodeVerifier = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQ"

	testNonce = "test-oidc-nonce"

	testRawIDToken = "test-cognito-id-token"

	testCognitoSubject = "test-cognito-subject"

	testCognitoEmail = "user@example.com"

	testRawSessionID = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFG"

	testFrontendRedirectURL = "https://example.com/"

	testSessionLifetime = 8 * time.Hour
)

var (
	errConsumeOAuthState = errors.New(
		"consume OAuth state failed",
	)

	errExchangeAuthorizationCode = errors.New(
		"exchange authorization code failed",
	)

	errVerifyIDToken = errors.New(
		"verify ID token failed",
	)

	errGenerateSessionID = errors.New(
		"generate session ID failed",
	)

	errSaveSession = errors.New(
		"save session failed",
	)

	errDeleteSession = errors.New(
		"delete session failed",
	)

	errWriteSessionCookie = errors.New(
		"write session cookie failed",
	)
)

type callbackTestDependencies struct {
	oauthStateStore *fakeOAuthStateStore
	cognitoClient   *fakeCognitoClient
	verifier        *fakeIDTokenVerifier
	idGenerator     *fakeSessionIDGenerator
	sessionStore    *fakeSessionStore
	cookieWriter    *fakeSessionCookieWriter
}

func newCallbackTestDependencies() callbackTestDependencies {
	return callbackTestDependencies{
		oauthStateStore: &fakeOAuthStateStore{
			state: newTestOAuthState(),
		},
		cognitoClient: &fakeCognitoClient{
			idToken: testRawIDToken,
		},
		verifier: &fakeIDTokenVerifier{
			identity: cognito.Identity{
				Subject: testCognitoSubject,
				Email:   testCognitoEmail,
			},
		},
		idGenerator: &fakeSessionIDGenerator{
			sessionID: testRawSessionID,
		},
		sessionStore: &fakeSessionStore{},
		cookieWriter: &fakeSessionCookieWriter{},
	}
}

func newCallbackRequest() *http.Request {
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/auth/callback",
		nil,
	)

	query := request.URL.Query()

	query.Set(
		"code",
		testAuthorizationCode,
	)
	query.Set(
		"state",
		testOAuthStateValue,
	)

	request.URL.RawQuery = query.Encode()

	return request
}

func newCallbackRequestWithQuery(
	code string,
	stateValue string,
) *http.Request {
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/auth/callback",
		nil,
	)

	query := request.URL.Query()

	if code != "" {
		query.Set(
			"code",
			code,
		)
	}

	if stateValue != "" {
		query.Set(
			"state",
			stateValue,
		)
	}

	request.URL.RawQuery = query.Encode()

	return request
}

func newCallbackResponseRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func testCallbackNow() time.Time {
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

func testSessionIDHash(
	t *testing.T,
) string {
	t.Helper()

	idHash, err := session.HashID(
		testRawSessionID,
	)
	if err != nil {
		t.Fatalf(
			"HashID() error = %v",
			err,
		)
	}

	return idHash
}

func assertNoSessionWasEstablished(
	t *testing.T,
	dependencies callbackTestDependencies,
) {
	t.Helper()

	if dependencies.sessionStore.saveCalls != 0 {
		t.Errorf(
			"Save() calls = %d, want 0",
			dependencies.sessionStore.saveCalls,
		)
	}

	if dependencies.cookieWriter.writeCalls != 0 {
		t.Errorf(
			"Write() calls = %d, want 0",
			dependencies.cookieWriter.writeCalls,
		)
	}
}

func assertEmptyResponseBody(
	t *testing.T,
	response *httptest.ResponseRecorder,
) {
	t.Helper()

	body := strings.TrimSpace(
		response.Body.String(),
	)

	if body != "" {
		t.Errorf(
			"response body = %q, want empty",
			body,
		)
	}
}

func validTestHandlerConfig() Config {
	return Config{
		RedirectURL:     testFrontendRedirectURL,
		SessionLifetime: testSessionLifetime,
	}
}

func newTestHandler(
	t *testing.T,
	dependencies callbackTestDependencies,
) *Handler {
	t.Helper()

	handler, err := NewHandler(
		dependencies.oauthStateStore,
		dependencies.cognitoClient,
		dependencies.verifier,
		dependencies.idGenerator,
		dependencies.sessionStore,
		dependencies.cookieWriter,
		validTestHandlerConfig(),
	)
	if err != nil {
		t.Fatalf(
			"NewHandler() error = %v",
			err,
		)
	}

	handler.now = testCallbackNow

	return handler
}
