package callback

import (
	"context"
	"net/http"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
)

type fakeSessionIDGenerator struct {
	sessionID string
	err       error

	generateCalls int
}

func (
	generator *fakeSessionIDGenerator,
) Generate() (
	string,
	error,
) {
	generator.generateCalls++

	if generator.err != nil {
		return "", generator.err
	}

	return generator.sessionID, nil
}

type fakeSessionStore struct {
	saveErr   error
	getErr    error
	deleteErr error

	saveCalls   int
	getCalls    int
	deleteCalls int

	savedSession  session.Session
	getIDHash     string
	deletedIDHash string

	contextWasNil bool
}

func (
	store *fakeSessionStore,
) Save(
	ctx context.Context,
	value session.Session,
) error {
	store.saveCalls++
	store.savedSession = value
	store.contextWasNil = ctx == nil

	if store.saveErr != nil {
		return store.saveErr
	}

	return nil
}

func (
	store *fakeSessionStore,
) Get(
	ctx context.Context,
	idHash string,
) (
	session.Session,
	error,
) {
	store.getCalls++
	store.getIDHash = idHash
	store.contextWasNil = ctx == nil

	if store.getErr != nil {
		return session.Session{},
			store.getErr
	}

	return session.Session{},
		nil
}

func (
	store *fakeSessionStore,
) Delete(
	ctx context.Context,
	idHash string,
) error {
	store.deleteCalls++
	store.deletedIDHash = idHash
	store.contextWasNil = ctx == nil

	if store.deleteErr != nil {
		return store.deleteErr
	}

	return nil
}

type fakeSessionCookieWriter struct {
	err error

	writeCalls     int
	rawSessionID   string
	responseWasNil bool
}

func (
	writer *fakeSessionCookieWriter,
) Write(
	response http.ResponseWriter,
	rawSessionID string,
) error {
	writer.writeCalls++
	writer.rawSessionID = rawSessionID
	writer.responseWasNil = response == nil

	if writer.err != nil {
		return writer.err
	}

	return nil
}
