package login

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
)

const (
	testLoginPath = "/api/auth/login"

	testAuthorizationURL = "https://example.auth.ap-northeast-1.amazoncognito.com/oauth2/authorize" +
		"?client_id=test-client-id" +
		"&response_type=code" +
		"&state=test-oauth-state"
)

var (
	errGenerateOAuthState = errors.New(
		"generate OAuth state failed",
	)
	errSaveOAuthState = errors.New(
		"save OAuth state failed",
	)
	errCreateAuthorizationURL = errors.New(
		"create authorization URL failed",
	)
)

type fakeStateGenerator struct {
	generateCalls int
	state         oauthstate.State
	err           error
}

func (
	generator *fakeStateGenerator,
) Generate() (
	oauthstate.State,
	error,
) {
	generator.generateCalls++

	if generator.err != nil {
		return oauthstate.State{},
			generator.err
	}

	return generator.state, nil
}

type fakeStateStore struct {
	saveCalls  int
	savedState oauthstate.State

	contextWasNil bool
	err           error
}

func (
	store *fakeStateStore,
) Save(
	ctx context.Context,
	state oauthstate.State,
) error {
	store.saveCalls++
	store.savedState = state
	store.contextWasNil = ctx == nil

	return store.err
}

type fakeAuthorizationClient struct {
	authorizationURLCalls int
	receivedState         oauthstate.State

	authorizationURL string
	err              error
}

func (
	client *fakeAuthorizationClient,
) AuthorizationURL(
	state oauthstate.State,
) (string, error) {
	client.authorizationURLCalls++
	client.receivedState = state

	if client.err != nil {
		return "", client.err
	}

	return client.authorizationURL, nil
}

type testHandlerDependencies struct {
	stateGenerator      *fakeStateGenerator
	stateStore          *fakeStateStore
	authorizationClient *fakeAuthorizationClient
}

func newTestHandlerDependencies() testHandlerDependencies {
	state := newTestLoginOAuthState()

	return testHandlerDependencies{
		stateGenerator: &fakeStateGenerator{
			state: state,
		},
		stateStore: &fakeStateStore{},
		authorizationClient: &fakeAuthorizationClient{
			authorizationURL: testAuthorizationURL,
		},
	}
}

func newTestHandler(
	t *testing.T,
	dependencies testHandlerDependencies,
) *Handler {
	t.Helper()

	handler, err := NewHandler(
		dependencies.stateGenerator,
		dependencies.stateStore,
		dependencies.authorizationClient,
	)
	if err != nil {
		t.Fatalf(
			"NewHandler() error = %v",
			err,
		)
	}

	return handler
}

func newTestLoginOAuthState() oauthstate.State {
	createdAt := time.Date(
		2026,
		time.August,
		30,
		12,
		34,
		56,
		0,
		time.UTC,
	)

	return oauthstate.State{
		Value:        "test-oauth-state",
		CodeVerifier: "vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv",
		Nonce:        "test-oidc-nonce",
		CreatedAt:    createdAt,
		ExpiresAt: createdAt.Add(
			10 * time.Minute,
		),
	}
}

func newTestLoginRequest() *http.Request {
	return httptest.NewRequest(
		http.MethodGet,
		testLoginPath,
		nil,
	)
}

func newTestResponseRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}
