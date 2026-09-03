package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ServeHTTP_Success(
	t *testing.T,
) {
	dependencies :=
		newTestDependencies()

	handler := newTestHandler(
		t,
		dependencies,
	)

	request :=
		newTestRequest(t)

	response :=
		httptest.NewRecorder()

	handler.ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusNoContent {
		t.Errorf(
			"status code = %d, want %d",
			response.Code,
			http.StatusNoContent,
		)
	}

	if response.Body.Len() != 0 {
		t.Errorf(
			"response body = %q, want empty",
			response.Body.String(),
		)
	}

	if got := response.Header().
		Get("Cache-Control"); got !=
		"no-store" {
		t.Errorf(
			"Cache-Control = %q, want %q",
			got,
			"no-store",
		)
	}

	if dependencies.resolver.resolveCalls !=
		1 {
		t.Errorf(
			"Resolve() calls = %d, want 1",
			dependencies.resolver.resolveCalls,
		)
	}

	if dependencies.resolver.receivedCtx ==
		nil {
		t.Error(
			"Resolve() context is nil",
		)
	}

	if dependencies.resolver.receivedRequest !=
		request {
		t.Error(
			"Resolve() request does not match",
		)
	}

	if dependencies.store.deleteCalls !=
		1 {
		t.Errorf(
			"Delete() store calls = %d, want 1",
			dependencies.store.deleteCalls,
		)
	}

	if dependencies.store.receivedCtx ==
		nil {
		t.Error(
			"Delete() store context is nil",
		)
	}

	if dependencies.store.receivedSessionHash !=
		testSessionIDHash {
		t.Errorf(
			"Delete() session ID hash = %q, want %q",
			dependencies.store.receivedSessionHash,
			testSessionIDHash,
		)
	}

	if dependencies.cookieDeleter.deleteCalls !=
		1 {
		t.Errorf(
			"Delete() cookie calls = %d, want 1",
			dependencies.cookieDeleter.deleteCalls,
		)
	}

	if dependencies.cookieDeleter.receivedWriter !=
		response {
		t.Error(
			"Delete() cookie writer does not match",
		)
	}
}
