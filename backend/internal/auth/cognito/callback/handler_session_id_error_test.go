package callback

import (
	"net/http"
	"testing"
)

func TestHandlerReturnsInternalServerErrorWhenSessionIDGenerationFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()
	dependencies.idGenerator.err =
		errGenerateSessionID

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

	if dependencies.idGenerator.generateCalls != 1 {
		t.Errorf(
			"Generate() calls = %d, want 1",
			dependencies.idGenerator.generateCalls,
		)
	}

	if dependencies.sessionStore.saveCalls != 0 {
		t.Errorf(
			"Save() calls = %d, want 0",
			dependencies.sessionStore.saveCalls,
		)
	}

	if dependencies.cookieWriter.writeCalls != 0 {
		t.Errorf(
			"Write() calls = %d, want 0",
			dependencies.cookieWriter.writeCalls,
		)
	}
}

func TestHandlerReturnsInternalServerErrorWhenGeneratedSessionIDIsEmpty(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()
	dependencies.idGenerator.sessionID = ""

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

	if dependencies.idGenerator.generateCalls != 1 {
		t.Errorf(
			"Generate() calls = %d, want 1",
			dependencies.idGenerator.generateCalls,
		)
	}

	if dependencies.sessionStore.saveCalls != 0 {
		t.Errorf(
			"Save() calls = %d, want 0",
			dependencies.sessionStore.saveCalls,
		)
	}

	if dependencies.cookieWriter.writeCalls != 0 {
		t.Errorf(
			"Write() calls = %d, want 0",
			dependencies.cookieWriter.writeCalls,
		)
	}
}
