package callback

import "testing"

func TestNewHandler(
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

	if handler == nil {
		t.Fatal(
			"NewHandler() handler = nil",
		)
	}

	if handler.oauthStateStore !=
		dependencies.oauthStateStore {
		t.Error(
			"NewHandler() OAuth state store was not assigned",
		)
	}

	if handler.cognitoClient !=
		dependencies.cognitoClient {
		t.Error(
			"NewHandler() Cognito client was not assigned",
		)
	}

	if handler.verifier !=
		dependencies.verifier {
		t.Error(
			"NewHandler() verifier was not assigned",
		)
	}

	if handler.idGenerator !=
		dependencies.idGenerator {
		t.Error(
			"NewHandler() session ID generator was not assigned",
		)
	}

	if handler.sessionStore !=
		dependencies.sessionStore {
		t.Error(
			"NewHandler() session store was not assigned",
		)
	}

	if handler.cookieWriter !=
		dependencies.cookieWriter {
		t.Error(
			"NewHandler() cookie writer was not assigned",
		)
	}

	if handler.redirectURL !=
		testFrontendRedirectURL {
		t.Errorf(
			"NewHandler() redirect URL = %q, want %q",
			handler.redirectURL,
			testFrontendRedirectURL,
		)
	}

	if handler.sessionLifetime !=
		testSessionLifetime {
		t.Errorf(
			"NewHandler() session lifetime = %v, want %v",
			handler.sessionLifetime,
			testSessionLifetime,
		)
	}

	if handler.now == nil {
		t.Error(
			"NewHandler() clock = nil",
		)
	}
}
