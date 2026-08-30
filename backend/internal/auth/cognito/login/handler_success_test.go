package login

import (
	"net/http"
	"strings"
	"testing"
)

func TestHandlerRedirectsToCognitoAuthorizationURL(
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
	response := newTestResponseRecorder()

	handler.ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusFound {
		t.Errorf(
			"response status = %d, want %d",
			response.Code,
			http.StatusFound,
		)
	}

	if location := response.Header().Get(
		"Location",
	); location != testAuthorizationURL {
		t.Errorf(
			"Location header = %q, want %q",
			location,
			testAuthorizationURL,
		)
	}

	if cacheControl := response.Header().Get(
		"Cache-Control",
	); cacheControl != "no-store" {
		t.Errorf(
			"Cache-Control header = %q, want %q",
			cacheControl,
			"no-store",
		)
	}

	if pragma := response.Header().Get(
		"Pragma",
	); pragma != "no-cache" {
		t.Errorf(
			"Pragma header = %q, want %q",
			pragma,
			"no-cache",
		)
	}
}

func TestHandlerGeneratesAndStoresOAuthState(
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
	response := newTestResponseRecorder()

	handler.ServeHTTP(
		response,
		request,
	)

	wantState :=
		dependencies.stateGenerator.state

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
		authorizationClient.
		receivedState != wantState {
		t.Errorf(
			"AuthorizationURL() state = %+v, want %+v",
			dependencies.
				authorizationClient.
				receivedState,
			wantState,
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

	if dependencies.
		stateStore.
		savedState != wantState {
		t.Errorf(
			"Save() state = %+v, want %+v",
			dependencies.
				stateStore.
				savedState,
			wantState,
		)
	}

	if dependencies.
		stateStore.
		contextWasNil {
		t.Error(
			"Save() context = nil, want request context",
		)
	}
}

func TestHandlerDoesNotExposeCodeVerifier(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newTestHandlerDependencies()

	handler := newTestHandler(
		t,
		dependencies,
	)

	response := newTestResponseRecorder()

	handler.ServeHTTP(
		response,
		newTestLoginRequest(),
	)

	codeVerifier :=
		dependencies.
			stateGenerator.
			state.
			CodeVerifier

	if strings.Contains(
		response.Header().Get("Location"),
		codeVerifier,
	) {
		t.Error(
			"Location header contains the PKCE code verifier",
		)
	}

	if strings.Contains(
		response.Body.String(),
		codeVerifier,
	) {
		t.Error(
			"response body contains the PKCE code verifier",
		)
	}
}
