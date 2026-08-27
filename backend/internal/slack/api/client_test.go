package api

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("uses provided HTTP client", func(t *testing.T) {
		t.Parallel()

		httpClient := &http.Client{}

		client := NewClient(httpClient)

		if client.httpClient != httpClient {
			t.Error("NewClient() did not use provided HTTP client")
		}

		if client.postMessageURL != slackPostMessageURL {
			t.Errorf(
				"postMessageURL = %q, want %q",
				client.postMessageURL,
				slackPostMessageURL,
			)
		}

		if client.openConversationURL != slackOpenConversationURL {
			t.Errorf(
				"openConversationURL = %q, want %q",
				client.openConversationURL,
				slackOpenConversationURL,
			)
		}
	})

	t.Run("uses default HTTP client when nil", func(t *testing.T) {
		t.Parallel()

		client := NewClient(nil)

		if client.httpClient != http.DefaultClient {
			t.Error("NewClient() did not use http.DefaultClient")
		}

		if client.postMessageURL != slackPostMessageURL {
			t.Errorf(
				"postMessageURL = %q, want %q",
				client.postMessageURL,
				slackPostMessageURL,
			)
		}

		if client.openConversationURL != slackOpenConversationURL {
			t.Errorf(
				"openConversationURL = %q, want %q",
				client.openConversationURL,
				slackOpenConversationURL,
			)
		}
	})
}
