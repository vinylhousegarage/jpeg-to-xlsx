package callback

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/apierror"
	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
)

func TestHandler_ServeHTTP_RejectsUnsupportedMethod(t *testing.T) {
	t.Parallel()

	resolver := &stubSessionResolver{}
	exchanger := &stubCodeExchanger{}
	opener := &stubConversationOpener{}
	store := &stubTokenStore{}

	handler := NewHandler(
		testRedirectURI,
		true,
		resolver,
		exchanger,
		opener,
		store,
		zap.NewNop(),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		testCallbackPath,
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertErrorResponse(
		t,
		rec,
		http.StatusMethodNotAllowed,
		apierror.ErrorCodeInvalidMethod,
	)

	if resolver.called {
		t.Error("Resolve() was called")
	}

	if exchanger.called {
		t.Error("ExchangeCode() was called")
	}

	if opener.called {
		t.Error("OpenConversation() was called")
	}

	if store.called {
		t.Error("Save() was called")
	}

	if len(rec.Result().Cookies()) != 0 {
		t.Error("unexpected response cookie")
	}
}

func TestHandler_ServeHTTP_MissingStateCookie(t *testing.T) {
	t.Parallel()

	resolver := &stubSessionResolver{
		session: authsession.Session{
			CognitoSub: testCognitoSub,
		},
	}
	exchanger := &stubCodeExchanger{}
	opener := &stubConversationOpener{}
	store := &stubTokenStore{}

	handler := NewHandler(
		testRedirectURI,
		true,
		resolver,
		exchanger,
		opener,
		store,
		zap.NewNop(),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		testCallbackPath+
			"?state="+testState+
			"&code="+testCode,
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertErrorResponse(
		t,
		rec,
		http.StatusBadRequest,
		apierror.ErrorCodeMissingState,
	)

	if !resolver.called {
		t.Fatal("Resolve() was not called")
	}

	if exchanger.called {
		t.Error("ExchangeCode() was called")
	}

	if opener.called {
		t.Error("OpenConversation() was called")
	}

	if store.called {
		t.Error("Save() was called")
	}

	if len(rec.Result().Cookies()) != 0 {
		t.Error("unexpected response cookie")
	}
}

func TestHandler_ServeHTTP_InvalidState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		queryState  string
		cookieState string
	}{
		{
			name:        "missing query state",
			queryState:  "",
			cookieState: testState,
		},
		{
			name:        "state mismatch",
			queryState:  "different-state",
			cookieState: testState,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resolver := &stubSessionResolver{
				session: authsession.Session{
					CognitoSub: testCognitoSub,
				},
			}
			exchanger := &stubCodeExchanger{}
			opener := &stubConversationOpener{}
			store := &stubTokenStore{}

			handler := NewHandler(
				testRedirectURI,
				true,
				resolver,
				exchanger,
				opener,
				store,
				zap.NewNop(),
			)

			req := httptest.NewRequest(
				http.MethodGet,
				testCallbackPath+
					"?state="+tt.queryState+
					"&code="+testCode,
				nil,
			)

			req.AddCookie(
				&http.Cookie{
					Name:  oauthStateCookieName,
					Value: tt.cookieState,
				},
			)

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assertErrorResponse(
				t,
				rec,
				http.StatusBadRequest,
				apierror.ErrorCodeInvalidState,
			)

			if !resolver.called {
				t.Fatal("Resolve() was not called")
			}

			if exchanger.called {
				t.Error("ExchangeCode() was called")
			}

			if opener.called {
				t.Error("OpenConversation() was called")
			}

			if store.called {
				t.Error("Save() was called")
			}

			if len(rec.Result().Cookies()) != 0 {
				t.Error("state cookie was deleted before successful validation")
			}
		})
	}
}

func TestHandler_ServeHTTP_MissingCode(t *testing.T) {
	t.Parallel()

	resolver := &stubSessionResolver{
		session: authsession.Session{
			CognitoSub: testCognitoSub,
		},
	}
	exchanger := &stubCodeExchanger{}
	opener := &stubConversationOpener{}
	store := &stubTokenStore{}

	handler := NewHandler(
		testRedirectURI,
		false,
		resolver,
		exchanger,
		opener,
		store,
		zap.NewNop(),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		testCallbackPath+
			"?state="+testState,
		nil,
	)

	req.AddCookie(
		&http.Cookie{
			Name:  oauthStateCookieName,
			Value: testState,
		},
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertErrorResponse(
		t,
		rec,
		http.StatusBadRequest,
		apierror.ErrorCodeMissingCode,
	)

	if !resolver.called {
		t.Fatal("Resolve() was not called")
	}

	if exchanger.called {
		t.Error("ExchangeCode() was called")
	}

	if opener.called {
		t.Error("OpenConversation() was called")
	}

	if store.called {
		t.Error("Save() was called")
	}

	assertDeleteStateCookie(t, rec, false)
}
