package logout

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
)

const (
	testCognitoSub    = "test-cognito-sub"
	testSessionIDHash = "test-session-id-hash"
	testLogoutURL     = "https://example.auth.ap-northeast-1." +
		"amazoncognito.com/logout?" +
		"client_id=test-client-id&" +
		"logout_uri=https%3A%2F%2Fexample.com%2F"
)

var (
	testSessionCreatedAt = time.Date(
		2026,
		time.September,
		3,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	testSessionExpiresAt = testSessionCreatedAt.Add(time.Hour)
)

type fakeSessionResolver struct {
	session         authsession.Session
	err             error
	resolveCalls    int
	receivedCtx     context.Context
	receivedRequest *http.Request
}

func (resolver *fakeSessionResolver) Resolve(
	ctx context.Context,
	request *http.Request,
) (
	authsession.Session,
	error,
) {
	resolver.resolveCalls++
	resolver.receivedCtx = ctx
	resolver.receivedRequest = request

	if resolver.err != nil {
		return authsession.Session{}, resolver.err
	}

	return resolver.session, nil
}

type fakeSessionStore struct {
	deleteErr           error
	deleteCalls         int
	receivedCtx         context.Context
	receivedSessionHash string
}

func (store *fakeSessionStore) Delete(
	ctx context.Context,
	sessionIDHash string,
) error {
	store.deleteCalls++
	store.receivedCtx = ctx
	store.receivedSessionHash = sessionIDHash

	return store.deleteErr
}

type fakeCookieDeleter struct {
	deleteCalls    int
	receivedWriter http.ResponseWriter
}

func (deleter *fakeCookieDeleter) Delete(
	writer http.ResponseWriter,
) {
	deleter.deleteCalls++
	deleter.receivedWriter = writer
}

type fakeLogoutURLProvider struct {
	logoutURL string
	err       error
	calls     int
}

func (provider *fakeLogoutURLProvider) LogoutURL() (
	string,
	error,
) {
	provider.calls++

	if provider.err != nil {
		return "", provider.err
	}

	return provider.logoutURL, nil
}

type testDependencies struct {
	resolver          *fakeSessionResolver
	store             *fakeSessionStore
	cookieDeleter     *fakeCookieDeleter
	logoutURLProvider *fakeLogoutURLProvider
}

func newTestDependencies() *testDependencies {
	return &testDependencies{
		resolver: &fakeSessionResolver{
			session: authsession.Session{
				IDHash:     testSessionIDHash,
				CognitoSub: testCognitoSub,
				CreatedAt:  testSessionCreatedAt,
				ExpiresAt:  testSessionExpiresAt,
			},
		},
		store:         &fakeSessionStore{},
		cookieDeleter: &fakeCookieDeleter{},
		logoutURLProvider: &fakeLogoutURLProvider{
			logoutURL: testLogoutURL,
		},
	}
}

func newTestHandler(
	t *testing.T,
	dependencies *testDependencies,
) *Handler {
	t.Helper()

	handler, err := NewHandler(
		dependencies.resolver,
		dependencies.store,
		dependencies.cookieDeleter,
		dependencies.logoutURLProvider,
	)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	return handler
}

func newTestRequest(t *testing.T) *http.Request {
	t.Helper()

	return httptest.NewRequest(
		http.MethodPost,
		"/api/auth/logout",
		nil,
	)
}
