package status

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

	testSessionExpiresAt = testSessionCreatedAt.Add(
		time.Hour,
	)
)

type fakeSessionResolver struct {
	session         authsession.Session
	err             error
	resolveCalls    int
	receivedCtx     context.Context
	receivedRequest *http.Request
}

func (
	resolver *fakeSessionResolver,
) Resolve(
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

type testDependencies struct {
	resolver *fakeSessionResolver
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
	}
}

func newTestHandler(
	t *testing.T,
	dependencies *testDependencies,
) *Handler {
	t.Helper()

	handler, err := NewHandler(
		dependencies.resolver,
	)
	if err != nil {
		t.Fatalf(
			"NewHandler() error = %v",
			err,
		)
	}

	return handler
}

func newTestRequest(
	t *testing.T,
) *http.Request {
	t.Helper()

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/auth/session",
		nil,
	)

	return request
}
