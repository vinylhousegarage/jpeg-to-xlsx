package callback

import (
	"net/http"
	"testing"
)

func TestHandlerReturnsInternalServerErrorWhenSessionSaveFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()
	dependencies.sessionStore.saveErr =
		errSaveSession

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

	if dependencies.sessionStore.saveCalls != 1 {
		t.Errorf(
			"Save() calls = %d, want 1",
			dependencies.sessionStore.saveCalls,
		)
	}

	if dependencies.cookieWriter.writeCalls != 0 {
		t.Errorf(
			"Write() calls = %d, want 0",
			dependencies.cookieWriter.writeCalls,
		)
	}

	if dependencies.sessionStore.deleteCalls != 0 {
		t.Errorf(
			"Delete() calls = %d, want 0",
			dependencies.sessionStore.deleteCalls,
		)
	}

	if location := response.Header().
		Get("Location"); location != "" {
		t.Errorf(
			"response Location = %q, want empty",
			location,
		)
	}
}
