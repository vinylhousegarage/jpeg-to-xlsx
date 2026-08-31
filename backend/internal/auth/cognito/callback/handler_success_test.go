package callback

import (
	"net/http"
	"testing"
)

func TestHandlerEstablishesSessionAndRedirects(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()

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

	request := newCallbackRequest()
	response := newCallbackResponseRecorder()

	handler.ServeHTTP(
		response,
		request,
	)

	if dependencies.oauthStateStore.consumeCalls != 1 {
		t.Errorf(
			"Consume() calls = %d, want 1",
			dependencies.oauthStateStore.consumeCalls,
		)
	}

	if dependencies.oauthStateStore.consumedValue !=
		testOAuthStateValue {
		t.Errorf(
			"Consume() state = %q, want %q",
			dependencies.oauthStateStore.consumedValue,
			testOAuthStateValue,
		)
	}

	if dependencies.oauthStateStore.contextWasNil {
		t.Error(
			"Consume() context was nil",
		)
	}

	if dependencies.cognitoClient.exchangeCalls != 1 {
		t.Errorf(
			"Exchange() calls = %d, want 1",
			dependencies.cognitoClient.exchangeCalls,
		)
	}

	if dependencies.cognitoClient.code !=
		testAuthorizationCode {
		t.Errorf(
			"Exchange() code = %q, want %q",
			dependencies.cognitoClient.code,
			testAuthorizationCode,
		)
	}

	if dependencies.cognitoClient.codeVerifier !=
		testCodeVerifier {
		t.Errorf(
			"Exchange() code verifier = %q, want %q",
			dependencies.cognitoClient.codeVerifier,
			testCodeVerifier,
		)
	}

	if dependencies.cognitoClient.contextWasNil {
		t.Error(
			"Exchange() context was nil",
		)
	}

	if dependencies.verifier.verifyCalls != 1 {
		t.Errorf(
			"Verify() calls = %d, want 1",
			dependencies.verifier.verifyCalls,
		)
	}

	if dependencies.verifier.rawIDToken !=
		testRawIDToken {
		t.Errorf(
			"Verify() ID token = %q, want %q",
			dependencies.verifier.rawIDToken,
			testRawIDToken,
		)
	}

	if dependencies.verifier.expectedNonce !=
		testNonce {
		t.Errorf(
			"Verify() nonce = %q, want %q",
			dependencies.verifier.expectedNonce,
			testNonce,
		)
	}

	if dependencies.verifier.contextWasNil {
		t.Error(
			"Verify() context was nil",
		)
	}

	if dependencies.idGenerator.generateCalls != 1 {
		t.Errorf(
			"Generate() calls = %d, want 1",
			dependencies.idGenerator.generateCalls,
		)
	}

	if dependencies.sessionStore.saveCalls != 1 {
		t.Fatalf(
			"Save() calls = %d, want 1",
			dependencies.sessionStore.saveCalls,
		)
	}

	savedSession :=
		dependencies.sessionStore.savedSession

	wantIDHash := testSessionIDHash(t)

	if savedSession.IDHash != wantIDHash {
		t.Errorf(
			"saved session ID hash = %q, want %q",
			savedSession.IDHash,
			wantIDHash,
		)
	}

	if savedSession.CognitoSub !=
		testCognitoSubject {
		t.Errorf(
			"saved Cognito subject = %q, want %q",
			savedSession.CognitoSub,
			testCognitoSubject,
		)
	}

	wantCreatedAt := testCallbackNow()

	if !savedSession.CreatedAt.Equal(
		wantCreatedAt,
	) {
		t.Errorf(
			"saved session CreatedAt = %v, want %v",
			savedSession.CreatedAt,
			wantCreatedAt,
		)
	}

	wantExpiresAt := wantCreatedAt.Add(
		testSessionLifetime,
	)

	if !savedSession.ExpiresAt.Equal(
		wantExpiresAt,
	) {
		t.Errorf(
			"saved session ExpiresAt = %v, want %v",
			savedSession.ExpiresAt,
			wantExpiresAt,
		)
	}

	if dependencies.sessionStore.contextWasNil {
		t.Error(
			"Save() context was nil",
		)
	}

	if dependencies.cookieWriter.writeCalls != 1 {
		t.Fatalf(
			"Write() calls = %d, want 1",
			dependencies.cookieWriter.writeCalls,
		)
	}

	if dependencies.cookieWriter.rawSessionID !=
		testRawSessionID {
		t.Errorf(
			"Write() session ID = %q, want %q",
			dependencies.cookieWriter.rawSessionID,
			testRawSessionID,
		)
	}

	if dependencies.cookieWriter.responseWasNil {
		t.Error(
			"Write() response writer was nil",
		)
	}

	if dependencies.sessionStore.deleteCalls != 0 {
		t.Errorf(
			"Delete() calls = %d, want 0",
			dependencies.sessionStore.deleteCalls,
		)
	}

	if response.Code != http.StatusSeeOther {
		t.Errorf(
			"response status = %d, want %d",
			response.Code,
			http.StatusSeeOther,
		)
	}

	if location := response.Header().
		Get("Location"); location !=
		testFrontendRedirectURL {
		t.Errorf(
			"response Location = %q, want %q",
			location,
			testFrontendRedirectURL,
		)
	}

	assertEmptyResponseBody(
		t,
		response,
	)
}
