package login

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerRejectsUnsupportedMethod(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newTestHandlerDependencies()

	handler := newTestHandler(
		t,
		dependencies,
	)

	request := newTestLoginRequest()
	request.Method = http.MethodPost

	response := newTestResponseRecorder()

	handler.ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusMethodNotAllowed {
		t.Errorf(
			"response status = %d, want %d",
			response.Code,
			http.StatusMethodNotAllowed,
		)
	}

	if got := response.Header().Get(
		"Allow",
	); got != http.MethodGet {
		t.Errorf(
			"Allow header = %q, want %q",
			got,
			http.MethodGet,
		)
	}

	if dependencies.
		stateGenerator.
		generateCalls != 0 {
		t.Errorf(
			"Generate() calls = %d, want 0",
			dependencies.
				stateGenerator.
				generateCalls,
		)
	}

	if dependencies.
		authorizationClient.
		authorizationURLCalls != 0 {
		t.Errorf(
			"AuthorizationURL() calls = %d, want 0",
			dependencies.
				authorizationClient.
				authorizationURLCalls,
		)
	}

	if dependencies.
		stateStore.
		saveCalls != 0 {
		t.Errorf(
			"Save() calls = %d, want 0",
			dependencies.
				stateStore.
				saveCalls,
		)
	}
}

func TestHandlerReturnsErrorWhenStateGenerationFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newTestHandlerDependencies()

	dependencies.stateGenerator.err =
		errGenerateOAuthState

	handler := newTestHandler(
		t,
		dependencies,
	)

	response := newTestResponseRecorder()

	handler.ServeHTTP(
		response,
		newTestLoginRequest(),
	)

	assertInternalServerError(
		t,
		response,
		errGenerateOAuthState.Error(),
	)

	if dependencies.
		stateGenerator.
		generateCalls != 1 {
		t.Errorf(
			"Generate() calls = %d, want 1",
			dependencies.
				stateGenerator.
				generateCalls,
		)
	}

	if dependencies.
		authorizationClient.
		authorizationURLCalls != 0 {
		t.Errorf(
			"AuthorizationURL() calls = %d, want 0",
			dependencies.
				authorizationClient.
				authorizationURLCalls,
		)
	}

	if dependencies.
		stateStore.
		saveCalls != 0 {
		t.Errorf(
			"Save() calls = %d, want 0",
			dependencies.
				stateStore.
				saveCalls,
		)
	}
}

func TestHandlerReturnsErrorWhenAuthorizationURLCreationFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newTestHandlerDependencies()

	dependencies.authorizationClient.err =
		errCreateAuthorizationURL

	handler := newTestHandler(
		t,
		dependencies,
	)

	response := newTestResponseRecorder()

	handler.ServeHTTP(
		response,
		newTestLoginRequest(),
	)

	assertInternalServerError(
		t,
		response,
		errCreateAuthorizationURL.Error(),
	)

	if dependencies.
		stateGenerator.
		generateCalls != 1 {
		t.Errorf(
			"Generate() calls = %d, want 1",
			dependencies.
				stateGenerator.
				generateCalls,
		)
	}

	if dependencies.
		authorizationClient.
		authorizationURLCalls != 1 {
		t.Errorf(
			"AuthorizationURL() calls = %d, want 1",
			dependencies.
				authorizationClient.
				authorizationURLCalls,
		)
	}

	if dependencies.
		stateStore.
		saveCalls != 0 {
		t.Errorf(
			"Save() calls = %d, want 0",
			dependencies.
				stateStore.
				saveCalls,
		)
	}
}

func TestHandlerReturnsErrorWhenStateSaveFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newTestHandlerDependencies()

	dependencies.stateStore.err =
		errSaveOAuthState

	handler := newTestHandler(
		t,
		dependencies,
	)

	response := newTestResponseRecorder()

	handler.ServeHTTP(
		response,
		newTestLoginRequest(),
	)

	assertInternalServerError(
		t,
		response,
		errSaveOAuthState.Error(),
	)

	if dependencies.
		stateGenerator.
		generateCalls != 1 {
		t.Errorf(
			"Generate() calls = %d, want 1",
			dependencies.
				stateGenerator.
				generateCalls,
		)
	}

	if dependencies.
		authorizationClient.
		authorizationURLCalls != 1 {
		t.Errorf(
			"AuthorizationURL() calls = %d, want 1",
			dependencies.
				authorizationClient.
				authorizationURLCalls,
		)
	}

	if dependencies.
		stateStore.
		saveCalls != 1 {
		t.Errorf(
			"Save() calls = %d, want 1",
			dependencies.
				stateStore.
				saveCalls,
		)
	}
}

func assertInternalServerError(
	t *testing.T,
	response *httptest.ResponseRecorder,
	sensitiveMessage string,
) {
	t.Helper()

	if response.Code !=
		http.StatusInternalServerError {
		t.Errorf(
			"response status = %d, want %d",
			response.Code,
			http.StatusInternalServerError,
		)
	}

	if location := response.Header().Get(
		"Location",
	); location != "" {
		t.Errorf(
			"Location header = %q, want empty",
			location,
		)
	}

	if strings.Contains(
		response.Body.String(),
		sensitiveMessage,
	) {
		t.Errorf(
			"response body contains internal error %q",
			sensitiveMessage,
		)
	}
}
