package callback

import (
	"net/http"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
)

func TestHandlerReturnsBadRequestWhenOAuthStateIsNotFound(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()
	dependencies.oauthStateStore.err =
		oauthstate.ErrNotFound

	handler := newTestHandler(
		t,
		dependencies,
	)

	response :=
		newCallbackResponseRecorder()

	handler.ServeHTTP(
		response,
		newCallbackRequest(),
	)

	if response.Code !=
		http.StatusBadRequest {
		t.Errorf(
			"response status = %d, want %d",
			response.Code,
			http.StatusBadRequest,
		)
	}

	if dependencies.oauthStateStore.consumeCalls != 1 {
		t.Errorf(
			"Consume() calls = %d, want 1",
			dependencies.oauthStateStore.consumeCalls,
		)
	}

	if dependencies.cognitoClient.exchangeCalls != 0 {
		t.Errorf(
			"Exchange() calls = %d, want 0",
			dependencies.cognitoClient.exchangeCalls,
		)
	}

	if dependencies.verifier.verifyCalls != 0 {
		t.Errorf(
			"Verify() calls = %d, want 0",
			dependencies.verifier.verifyCalls,
		)
	}

	assertNoSessionWasEstablished(
		t,
		dependencies,
	)
}

func TestHandlerReturnsInternalServerErrorWhenOAuthStateStoreFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()
	dependencies.oauthStateStore.err =
		errConsumeOAuthState

	handler := newTestHandler(
		t,
		dependencies,
	)

	response :=
		newCallbackResponseRecorder()

	handler.ServeHTTP(
		response,
		newCallbackRequest(),
	)

	if response.Code !=
		http.StatusInternalServerError {
		t.Errorf(
			"response status = %d, want %d",
			response.Code,
			http.StatusInternalServerError,
		)
	}

	if dependencies.oauthStateStore.consumeCalls != 1 {
		t.Errorf(
			"Consume() calls = %d, want 1",
			dependencies.oauthStateStore.consumeCalls,
		)
	}

	if dependencies.cognitoClient.exchangeCalls != 0 {
		t.Errorf(
			"Exchange() calls = %d, want 0",
			dependencies.cognitoClient.exchangeCalls,
		)
	}

	if dependencies.verifier.verifyCalls != 0 {
		t.Errorf(
			"Verify() calls = %d, want 0",
			dependencies.verifier.verifyCalls,
		)
	}

	assertNoSessionWasEstablished(
		t,
		dependencies,
	)
}

func TestHandlerReturnsBadGatewayWhenCodeExchangeFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()
	dependencies.cognitoClient.err =
		errExchangeAuthorizationCode

	handler := newTestHandler(
		t,
		dependencies,
	)

	response :=
		newCallbackResponseRecorder()

	handler.ServeHTTP(
		response,
		newCallbackRequest(),
	)

	if response.Code !=
		http.StatusBadGateway {
		t.Errorf(
			"response status = %d, want %d",
			response.Code,
			http.StatusBadGateway,
		)
	}

	if dependencies.oauthStateStore.consumeCalls != 1 {
		t.Errorf(
			"Consume() calls = %d, want 1",
			dependencies.oauthStateStore.consumeCalls,
		)
	}

	if dependencies.cognitoClient.exchangeCalls != 1 {
		t.Errorf(
			"Exchange() calls = %d, want 1",
			dependencies.cognitoClient.exchangeCalls,
		)
	}

	if dependencies.verifier.verifyCalls != 0 {
		t.Errorf(
			"Verify() calls = %d, want 0",
			dependencies.verifier.verifyCalls,
		)
	}

	if dependencies.idGenerator.generateCalls != 0 {
		t.Errorf(
			"Generate() calls = %d, want 0",
			dependencies.idGenerator.generateCalls,
		)
	}

	assertNoSessionWasEstablished(
		t,
		dependencies,
	)
}

func TestHandlerReturnsUnauthorizedWhenIDTokenVerificationFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()
	dependencies.verifier.err =
		errVerifyIDToken

	handler := newTestHandler(
		t,
		dependencies,
	)

	response :=
		newCallbackResponseRecorder()

	handler.ServeHTTP(
		response,
		newCallbackRequest(),
	)

	if response.Code !=
		http.StatusUnauthorized {
		t.Errorf(
			"response status = %d, want %d",
			response.Code,
			http.StatusUnauthorized,
		)
	}

	if dependencies.oauthStateStore.consumeCalls != 1 {
		t.Errorf(
			"Consume() calls = %d, want 1",
			dependencies.oauthStateStore.consumeCalls,
		)
	}

	if dependencies.cognitoClient.exchangeCalls != 1 {
		t.Errorf(
			"Exchange() calls = %d, want 1",
			dependencies.cognitoClient.exchangeCalls,
		)
	}

	if dependencies.verifier.verifyCalls != 1 {
		t.Errorf(
			"Verify() calls = %d, want 1",
			dependencies.verifier.verifyCalls,
		)
	}

	if dependencies.idGenerator.generateCalls != 0 {
		t.Errorf(
			"Generate() calls = %d, want 0",
			dependencies.idGenerator.generateCalls,
		)
	}

	assertNoSessionWasEstablished(
		t,
		dependencies,
	)
}
