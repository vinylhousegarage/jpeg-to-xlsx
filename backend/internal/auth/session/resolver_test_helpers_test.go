package session

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	testResolverRawSessionID = "test-raw-session-id"

	testResolverCognitoSub = "test-cognito-subject"
)

var (
	errResolverReadCookie = errors.New(
		"read session cookie failed",
	)

	errResolverGetSession = errors.New(
		"get session failed",
	)
)

type fakeResolverCookieReader struct {
	rawSessionID string
	err          error

	readCalls int
	request   *http.Request
}

func (
	reader *fakeResolverCookieReader,
) Read(
	request *http.Request,
) (
	string,
	error,
) {
	reader.readCalls++
	reader.request = request

	if reader.err != nil {
		return "", reader.err
	}

	return reader.rawSessionID, nil
}

type fakeResolverSessionStore struct {
	session Session
	getErr  error

	getCalls int
	getCtx   context.Context
	idHash   string

	saveCalls   int
	deleteCalls int
}

func (
	store *fakeResolverSessionStore,
) Save(
	_ context.Context,
	_ Session,
) error {
	store.saveCalls++

	return nil
}

func (
	store *fakeResolverSessionStore,
) Get(
	ctx context.Context,
	idHash string,
) (
	Session,
	error,
) {
	store.getCalls++
	store.getCtx = ctx
	store.idHash = idHash

	if store.getErr != nil {
		return Session{}, store.getErr
	}

	return store.session, nil
}

func (
	store *fakeResolverSessionStore,
) Delete(
	_ context.Context,
	_ string,
) error {
	store.deleteCalls++

	return nil
}

func newResolverTestDependencies(
	t *testing.T,
) (
	*fakeResolverCookieReader,
	*fakeResolverSessionStore,
) {
	t.Helper()

	sessionValue := newResolverTestSession(t)

	return &fakeResolverCookieReader{
			rawSessionID: testResolverRawSessionID,
		}, &fakeResolverSessionStore{
			session: sessionValue,
		}
}

func newResolverTestSession(
	t *testing.T,
) Session {
	t.Helper()

	idHash := resolverTestSessionIDHash(t)
	now := resolverTestNow()

	return Session{
		IDHash:     idHash,
		CognitoSub: testResolverCognitoSub,
		CreatedAt:  now.Add(-time.Hour),
		ExpiresAt:  now.Add(time.Hour),
	}
}

func newResolverTestRequest() *http.Request {
	return httptest.NewRequest(
		http.MethodGet,
		"/api/test",
		nil,
	)
}

func resolverTestNow() time.Time {
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

func resolverTestSessionIDHash(
	t *testing.T,
) string {
	t.Helper()

	idHash, err := HashID(
		testResolverRawSessionID,
	)
	if err != nil {
		t.Fatalf(
			"HashID() error = %v",
			err,
		)
	}

	return idHash
}
