package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_ExchangeCode_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf(
					"method = %s, want %s",
					r.Method,
					http.MethodPost,
				)
				return
			}

			user, pass, ok := r.BasicAuth()
			if !ok {
				t.Error("missing basic auth")
				return
			}

			if user != testClientID {
				t.Errorf(
					"clientID = %q, want %q",
					user,
					testClientID,
				)
			}

			if pass != testClientSecret {
				t.Errorf(
					"clientSecret = %q, want %q",
					pass,
					testClientSecret,
				)
			}

			contentType := r.Header.Get("Content-Type")
			if contentType != "application/x-www-form-urlencoded" {
				t.Errorf(
					"Content-Type = %q, want %q",
					contentType,
					"application/x-www-form-urlencoded",
				)
			}

			if err := r.ParseForm(); err != nil {
				t.Errorf(
					"failed to parse form: %v",
					err,
				)
				return
			}

			if got := r.Form.Get("code"); got != testCode {
				t.Errorf(
					"code = %q, want %q",
					got,
					testCode,
				)
			}

			if got := r.Form.Get("redirect_uri"); got != testRedirectURI {
				t.Errorf(
					"redirect_uri = %q, want %q",
					got,
					testRedirectURI,
				)
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			response := tokenResponse{
				OK:          true,
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
			}
			response.Team.ID = "T123"
			response.AuthedUser.ID = "U123"

			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf(
					"failed to encode response: %v",
					err,
				)
			}
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	token, err := client.ExchangeCode(
		context.Background(),
		testCode,
		testRedirectURI,
	)
	if err != nil {
		t.Fatalf(
			"ExchangeCode() error = %v",
			err,
		)
	}

	if token.AccessToken != "xoxb-test" {
		t.Errorf(
			"AccessToken = %q, want %q",
			token.AccessToken,
			"xoxb-test",
		)
	}

	if token.BotUserID != "B123" {
		t.Errorf(
			"BotUserID = %q, want %q",
			token.BotUserID,
			"B123",
		)
	}

	if token.TeamID != "T123" {
		t.Errorf(
			"TeamID = %q, want %q",
			token.TeamID,
			"T123",
		)
	}

	if token.UserID != "U123" {
		t.Errorf(
			"UserID = %q, want %q",
			token.UserID,
			"U123",
		)
	}
}
