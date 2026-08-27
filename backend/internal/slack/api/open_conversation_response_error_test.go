package api

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClient_OpenConversation_DecodeError(t *testing.T) {
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
	client.openConversationURL = "https://example.com/conversations.open"

	channelID, err := client.OpenConversation(
		context.Background(),
		"xoxb-test",
		"U123",
	)
	if err == nil {
		t.Fatal("OpenConversation() error = nil, want an error")
	}

	if channelID != "" {
		t.Errorf(
			"OpenConversation() channelID = %q, want empty",
			channelID,
		)
	}

	const wantPrefix = "decode slack open conversation response:"
	if !strings.HasPrefix(err.Error(), wantPrefix) {
		t.Errorf(
			"OpenConversation() error = %q, want prefix %q",
			err.Error(),
			wantPrefix,
		)
	}
}

func TestClient_OpenConversation_SlackError(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						strings.NewReader(
							`{"ok":false,"error":"user_not_found"}`,
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
	if err == nil {
		t.Fatal("OpenConversation() error = nil, want an error")
	}

	if channelID != "" {
		t.Errorf(
			"OpenConversation() channelID = %q, want empty",
			channelID,
		)
	}

	const wantError = "slack open conversation failed: user_not_found"
	if err.Error() != wantError {
		t.Errorf(
			"OpenConversation() error = %q, want %q",
			err.Error(),
			wantError,
		)
	}
}

func TestClient_OpenConversation_MissingChannelID(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(
			func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						strings.NewReader(
							`{"ok":true,"channel":{"id":""}}`,
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
	if err == nil {
		t.Fatal("OpenConversation() error = nil, want an error")
	}

	if channelID != "" {
		t.Errorf(
			"OpenConversation() channelID = %q, want empty",
			channelID,
		)
	}

	const wantError = "slack open conversation response is missing channel.id"
	if err.Error() != wantError {
		t.Errorf(
			"OpenConversation() error = %q, want %q",
			err.Error(),
			wantError,
		)
	}
}
