package callback

import (
	"net/http"
	"testing"
)

func TestHandlerDeletesSessionWhenCookieWriteFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()
	dependencies.cookieWriter.err =
		errWriteSessionCookie

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

	if dependencies.cookieWriter.writeCalls != 1 {
		t.Errorf(
			"Write() calls = %d, want 1",
			dependencies.cookieWriter.writeCalls,
		)
	}

	if dependencies.sessionStore.deleteCalls != 1 {
		t.Fatalf(
			"Delete() calls = %d, want 1",
			dependencies.sessionStore.deleteCalls,
		)
	}

	wantIDHash := testSessionIDHash(t)

	if dependencies.sessionStore.deletedIDHash !=
		wantIDHash {
		t.Errorf(
			"Delete() ID hash = %q, want %q",
			dependencies.sessionStore.deletedIDHash,
			wantIDHash,
		)
	}

	if dependencies.sessionStore.contextWasNil {
		t.Error(
			"Delete() context was nil",
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

func TestHandlerStillReturnsErrorWhenSessionRollbackFails(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newCallbackTestDependencies()
	dependencies.cookieWriter.err =
		errWriteSessionCookie
	dependencies.sessionStore.deleteErr =
		errDeleteSession

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

	if dependencies.cookieWriter.writeCalls != 1 {
		t.Errorf(
			"Write() calls = %d, want 1",
			dependencies.cookieWriter.writeCalls,
		)
	}

	if dependencies.sessionStore.deleteCalls != 1 {
		t.Errorf(
			"Delete() calls = %d, want 1",
			dependencies.sessionStore.deleteCalls,
		)
	}

	wantIDHash := testSessionIDHash(t)

	if dependencies.sessionStore.deletedIDHash !=
		wantIDHash {
		t.Errorf(
			"Delete() ID hash = %q, want %q",
			dependencies.sessionStore.deletedIDHash,
			wantIDHash,
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
