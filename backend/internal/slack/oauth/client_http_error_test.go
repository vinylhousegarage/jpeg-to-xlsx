package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_ExchangeCode_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusInternalServerError)
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
		t.Fatal("ExchangeCode() error = nil, want an error")
	}

	if token != nil {
		t.Errorf(
			"ExchangeCode() token = %+v, want nil",
			token,
		)
	}

	if !strings.Contains(
		err.Error(),
		"returned status 500",
	) {
		t.Errorf(
			"ExchangeCode() error = %q, want status error",
			err,
		)
	}
}

func TestClient_ExchangeCode_InvalidJSON(t *testing.T) {
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

			if _, err := w.Write([]byte("{")); err != nil {
				t.Errorf(
					"failed to write response: %v",
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
		t.Fatal("ExchangeCode() error = nil, want an error")
	}

	if token != nil {
		t.Errorf(
			"ExchangeCode() token = %+v, want nil",
			token,
		)
	}

	if !strings.Contains(
		err.Error(),
		"decode Slack OAuth token response",
	) {
		t.Errorf(
			"ExchangeCode() error = %q, want decode error",
			err,
		)
	}
}

func TestClient_ExchangeCode_SlackError(t *testing.T) {
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
				OK:    false,
				Error: "invalid_code",
			}

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
		t.Fatal("ExchangeCode() error = nil, want an error")
	}

	if token != nil {
		t.Errorf(
			"ExchangeCode() token = %+v, want nil",
			token,
		)
	}

	if !strings.Contains(
		err.Error(),
		"invalid_code",
	) {
		t.Errorf(
			"ExchangeCode() error = %q, want invalid_code",
			err,
		)
	}
}
