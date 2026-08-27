package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClient_PostMessage_TransportError(t *testing.T) {
	t.Parallel()

	transportErr := errors.New("network unavailable")

	httpClient := &http.Client{
		Transport: roundTripFunc(
			func(_ *http.Request) (*http.Response, error) {
				return nil, transportErr
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

	if !errors.Is(err, transportErr) {
		t.Errorf(
			"PostMessage() error = %v, want wrapped error %v",
			err,
			transportErr,
		)
	}

	const wantPrefix = "send slack post message request:"
	if !strings.HasPrefix(err.Error(), wantPrefix) {
		t.Errorf(
			"PostMessage() error = %q, want prefix %q",
			err.Error(),
			wantPrefix,
		)
	}
}

func TestClient_PostMessage_HTTPStatusError(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body: io.NopCloser(
						strings.NewReader(""),
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

	const wantError = "slack post message request returned status 500"
	if err.Error() != wantError {
		t.Errorf(
			"PostMessage() error = %q, want %q",
			err.Error(),
			wantError,
		)
	}
}
