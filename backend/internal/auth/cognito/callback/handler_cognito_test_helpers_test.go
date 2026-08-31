package callback

import (
	"context"
	"time"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
)

type fakeOAuthStateStore struct {
	state oauthstate.State
	err   error

	consumeCalls  int
	consumedValue string
	contextWasNil bool
}

func (
	store *fakeOAuthStateStore,
) Save(
	_ context.Context,
	_ oauthstate.State,
) error {
	return nil
}

func (
	store *fakeOAuthStateStore,
) Consume(
	ctx context.Context,
	value string,
) (
	oauthstate.State,
	error,
) {
	store.consumeCalls++
	store.consumedValue = value
	store.contextWasNil = ctx == nil

	if store.err != nil {
		return oauthstate.State{},
			store.err
	}

	return store.state, nil
}

type fakeCognitoClient struct {
	idToken string
	err     error

	exchangeCalls int
	code          string
	codeVerifier  string
	contextWasNil bool
}

func (
	client *fakeCognitoClient,
) Exchange(
	ctx context.Context,
	code string,
	codeVerifier string,
) (
	string,
	error,
) {
	client.exchangeCalls++
	client.code = code
	client.codeVerifier = codeVerifier
	client.contextWasNil = ctx == nil

	if client.err != nil {
		return "", client.err
	}

	return client.idToken, nil
}

type fakeIDTokenVerifier struct {
	identity cognito.Identity
	err      error

	verifyCalls   int
	rawIDToken    string
	expectedNonce string
	contextWasNil bool
}

func (
	verifier *fakeIDTokenVerifier,
) Verify(
	ctx context.Context,
	rawIDToken string,
	expectedNonce string,
) (
	cognito.Identity,
	error,
) {
	verifier.verifyCalls++
	verifier.rawIDToken = rawIDToken
	verifier.expectedNonce = expectedNonce
	verifier.contextWasNil = ctx == nil

	if verifier.err != nil {
		return cognito.Identity{},
			verifier.err
	}

	return verifier.identity, nil
}

func newTestOAuthState() oauthstate.State {
	createdAt := testCallbackNow().
		Add(-time.Minute)

	return oauthstate.State{
		Value:        testOAuthStateValue,
		CodeVerifier: testCodeVerifier,
		Nonce:        testNonce,
		CreatedAt:    createdAt,
		ExpiresAt: createdAt.Add(
			10 * time.Minute,
		),
	}
}
