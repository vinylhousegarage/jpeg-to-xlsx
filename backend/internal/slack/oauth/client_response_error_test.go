package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_ExchangeCode_MissingAccessToken(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			response := tokenResponse{
				OK: true,
			}
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

	if err == nil {
		t.Fatal(
			"ExchangeCode() error = nil, want an error",
		)
	}

	if token != nil {
		t.Errorf(
			"ExchangeCode() token = %+v, want nil",
			token,
		)
	}

	if !strings.Contains(
		err.Error(),
		"missing access_token",
	) {
		t.Errorf(
			"ExchangeCode() error = %q, want missing access_token error",
			err,
		)
	}
}

func TestClient_ExchangeCode_MissingAuthedUserID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
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

	if err == nil {
		t.Fatal(
			"ExchangeCode() error = nil, want an error",
		)
	}

	if token != nil {
		t.Errorf(
			"ExchangeCode() token = %+v, want nil",
			token,
		)
	}

	if !strings.Contains(
		err.Error(),
		"missing authed_user.id",
	) {
		t.Errorf(
			"ExchangeCode() error = %q, want missing authed_user.id error",
			err,
		)
	}
}
