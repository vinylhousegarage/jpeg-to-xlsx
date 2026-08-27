package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClient_OpenConversation_Success(t *testing.T) {
	t.Parallel()

	var gotRequest openConversationRequest

	httpClient := &http.Client{
		Transport: roundTripFunc(
			func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodPost {
					t.Errorf(
						"request method = %q, want %q",
						req.Method,
						http.MethodPost,
					)
				}

				if got := req.Header.Get("Authorization"); got != "Bearer xoxb-test" {
					t.Errorf(
						"Authorization = %q, want %q",
						got,
						"Bearer xoxb-test",
					)
				}

				if got := req.Header.Get("Content-Type"); got != "application/json" {
					t.Errorf(
						"Content-Type = %q, want %q",
						got,
						"application/json",
					)
				}

				if err := json.NewDecoder(req.Body).Decode(&gotRequest); err != nil {
					t.Fatalf(
						"decode request body: %v",
						err,
					)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						strings.NewReader(
							`{"ok":true,"channel":{"id":"C123"}}`,
						),
					),
					Header: make(http.Header),
				}, nil
			},
		),
	}

	client := NewClient(httpClient)
	client.openConversationURL = "https://example.com/conversations.open"

	channelID, err := client.OpenConversation(
		context.Background(),
		"xoxb-test",
		"U123",
	)
	if err != nil {
		t.Fatalf("OpenConversation() error = %v", err)
	}

	if channelID != "C123" {
		t.Errorf(
			"OpenConversation() channelID = %q, want %q",
			channelID,
			"C123",
		)
	}

	if gotRequest.Users != "U123" {
		t.Errorf(
			"request users = %q, want %q",
			gotRequest.Users,
			"U123",
		)
	}
}
