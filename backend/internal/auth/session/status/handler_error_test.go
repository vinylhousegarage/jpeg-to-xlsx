package status

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ServeHTTP_InternalError(
	t *testing.T,
) {
	dependencies :=
		newTestDependencies()

	const resolverErrorMessage = "resolve session failed"

	dependencies.resolver.err =
		errors.New(
			resolverErrorMessage,
		)

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
		http.StatusInternalServerError {
		t.Errorf(
			"status code = %d, want %d",
			response.Code,
			http.StatusInternalServerError,
		)
	}

	if got := response.Header().
		Get("Content-Type"); got !=
		"application/json" {
		t.Errorf(
			"Content-Type = %q, want %q",
			got,
			"application/json",
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

	var body struct {
		Authenticated bool `json:"authenticated"`
	}

	if err := json.NewDecoder(
		response.Body,
	).Decode(&body); err != nil {
		t.Fatalf(
			"decode response body: %v",
			err,
		)
	}

	if body.Authenticated {
		t.Error(
			"authenticated = true, want false",
		)
	}

	if strings.Contains(
		response.Body.String(),
		resolverErrorMessage,
	) {
		t.Error(
			"response body exposes resolver error",
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
}
