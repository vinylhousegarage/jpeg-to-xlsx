package api

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClient_PostMessage_DecodeError(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						strings.NewReader(`{invalid`),
					),
					Header: make(http.Header),
				}, nil
			},
		),
	}

	client := NewClient(httpClient)
	client.postMessageURL = "https://example.com/chat.postMessage"

	err := client.PostMessage(
		context.Background(),
		"xoxb-test",
		"C123",
		Message{
			Text: "test message",
		},
	)
	if err == nil {
		t.Fatal("PostMessage() error = nil, want an error")
	}

	const wantPrefix = "decode slack post message response:"
	if !strings.HasPrefix(err.Error(), wantPrefix) {
		t.Errorf(
			"PostMessage() error = %q, want prefix %q",
			err.Error(),
			wantPrefix,
		)
	}
}

func TestClient_PostMessage_SlackError(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						strings.NewReader(
							`{"ok":false,"error":"channel_not_found"}`,
						),
					),
					Header: make(http.Header),
				}, nil
			},
		),
	}

	client := NewClient(httpClient)
	client.postMessageURL = "https://example.com/chat.postMessage"

	err := client.PostMessage(
		context.Background(),
		"xoxb-test",
		"C123",
		Message{
			Text: "test message",
		},
	)
	if err == nil {
		t.Fatal("PostMessage() error = nil, want an error")
	}

	const wantError = "slack post message failed: channel_not_found"
	if err.Error() != wantError {
		t.Errorf(
			"PostMessage() error = %q, want %q",
			err.Error(),
			wantError,
		)
	}
}
