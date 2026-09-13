package callback

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

func TestHandler_ServeHTTP_Success(t *testing.T) {
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

	const expectedLocation = "/?slack=connected"

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			rec.Code,
			http.StatusSeeOther,
			rec.Body.String(),
		)
	}

	if location := rec.Header().Get("Location"); location != expectedLocation {
		t.Errorf(
			"Location = %q, want %q",
			location,
			expectedLocation,
		)
	}

	if !resolver.called {
		t.Fatal("Resolve() was not called")
	}

	if resolver.request != req {
		t.Error("Resolve() received an unexpected request")
	}

	if resolver.ctx != req.Context() {
		t.Error("Resolve() received an unexpected context")
	}

	if !exchanger.called {
		t.Fatal("ExchangeCode() was not called")
	}

	if exchanger.gotCode != testCode {
		t.Errorf("ExchangeCode() code = %q, want %q", exchanger.gotCode, testCode)
	}

	if exchanger.gotRedirect != testRedirectURI {
		t.Errorf(
			"ExchangeCode() redirectURI = %q, want %q",
			exchanger.gotRedirect,
			testRedirectURI,
		)
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
		t.Errorf("Save() token = %+v, want %+v", store.gotToken, token)
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
