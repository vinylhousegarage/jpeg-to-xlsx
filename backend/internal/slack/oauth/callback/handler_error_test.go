package callback

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/apierror"
	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

func TestHandler_ServeHTTP_ExchangeCodeError(t *testing.T) {
	t.Parallel()

	exchangeErr := errors.New("slack token exchange failed")

	resolver := &stubSessionResolver{
		session: authsession.Session{
			CognitoSub: testCognitoSub,
		},
	}
	exchanger := &stubCodeExchanger{
		err: exchangeErr,
	}
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

	req := newValidCallbackRequest()
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertErrorResponse(
		t,
		rec,
		http.StatusInternalServerError,
		apierror.ErrorCodeInternal,
	)

	if !resolver.called {
		t.Fatal("Resolve() was not called")
	}

	if !exchanger.called {
		t.Fatal("ExchangeCode() was not called")
	}

	if opener.called {
		t.Error("OpenConversation() was called after ExchangeCode() failed")
	}

	if store.called {
		t.Error("Save() was called after ExchangeCode() failed")
	}

	assertDeleteStateCookie(t, rec, true)
}

func TestHandler_ServeHTTP_OpenConversationError(t *testing.T) {
	t.Parallel()

	token := &oauth.Token{
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
		TeamID:      "T123",
		UserID:      "U123",
	}

	resolver := &stubSessionResolver{
		session: authsession.Session{
			CognitoSub: testCognitoSub,
		},
	}
	exchanger := &stubCodeExchanger{
		token: token,
	}
	opener := &stubConversationOpener{
		err: errors.New("open conversation failed"),
	}
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

	req := newValidCallbackRequest()
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertErrorResponse(
		t,
		rec,
		http.StatusInternalServerError,
		apierror.ErrorCodeInternal,
	)

	if !resolver.called {
		t.Fatal("Resolve() was not called")
	}

	if !exchanger.called {
		t.Fatal("ExchangeCode() was not called")
	}

	if !opener.called {
		t.Fatal("OpenConversation() was not called")
	}

	if opener.gotAccessToken != token.AccessToken {
		t.Errorf(
			"OpenConversation() accessToken = %q, want %q",
			opener.gotAccessToken,
			token.AccessToken,
		)
	}

	if opener.gotUserID != token.UserID {
		t.Errorf(
			"OpenConversation() userID = %q, want %q",
			opener.gotUserID,
			token.UserID,
		)
	}

	if store.called {
		t.Error("Save() was called after OpenConversation() failed")
	}

	assertDeleteStateCookie(t, rec, true)
}

func TestHandler_ServeHTTP_SaveError(t *testing.T) {
	t.Parallel()

	token := &oauth.Token{
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
		TeamID:      "T123",
		UserID:      "U123",
	}

	resolver := &stubSessionResolver{
		session: authsession.Session{
			CognitoSub: testCognitoSub,
		},
	}
	exchanger := &stubCodeExchanger{
		token: token,
	}
	opener := &stubConversationOpener{
		channelID: "D123",
	}
	store := &stubTokenStore{
		err: errors.New("dynamodb save failed"),
	}

	handler := NewHandler(
		testRedirectURI,
		true,
		resolver,
		exchanger,
		opener,
		store,
		zap.NewNop(),
	)

	req := newValidCallbackRequest()
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertErrorResponse(
		t,
		rec,
		http.StatusInternalServerError,
		apierror.ErrorCodeInternal,
	)

	if !resolver.called {
		t.Fatal("Resolve() was not called")
	}

	if !exchanger.called {
		t.Fatal("ExchangeCode() was not called")
	}

	if !opener.called {
		t.Fatal("OpenConversation() was not called")
	}

	if opener.gotAccessToken != token.AccessToken {
		t.Errorf(
			"OpenConversation() accessToken = %q, want %q",
			opener.gotAccessToken,
			token.AccessToken,
		)
	}

	if opener.gotUserID != token.UserID {
		t.Errorf(
			"OpenConversation() userID = %q, want %q",
			opener.gotUserID,
			token.UserID,
		)
	}

	if !store.called {
		t.Fatal("Save() was not called")
	}

	if store.gotCognitoSub != testCognitoSub {
		t.Errorf(
			"Save() cognitoSub = %q, want %q",
			store.gotCognitoSub,
			testCognitoSub,
		)
	}

	if store.gotToken != token {
		t.Errorf(
			"Save() token = %+v, want %+v",
			store.gotToken,
			token,
		)
	}

	if store.gotToken.ChannelID != "D123" {
		t.Errorf(
			"Save() token ChannelID = %q, want %q",
			store.gotToken.ChannelID,
			"D123",
		)
	}

	assertDeleteStateCookie(t, rec, true)
}
