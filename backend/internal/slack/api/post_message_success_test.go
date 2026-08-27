package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClient_PostMessage_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		message    Message
		wantText   string
		wantBlocks bool
		wantButton *Button
	}{
		{
			name: "text only",
			message: Message{
				Text: "test message",
			},
			wantText:   "test message",
			wantBlocks: false,
		},
		{
			name: "with button",
			message: Message{
				Text: "撮影番号：001",
				Button: &Button{
					Text: "ダウンロード",
					URL:  "https://example.com/test.json",
				},
			},
			wantText:   "撮影番号：001",
			wantBlocks: true,
			wantButton: &Button{
				Text: "ダウンロード",
				URL:  "https://example.com/test.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var gotRequest postMessageRequest

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
								strings.NewReader(`{"ok":true}`),
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
				tt.message,
			)
			if err != nil {
				t.Fatalf("PostMessage() error = %v", err)
			}

			if gotRequest.Channel != "C123" {
				t.Errorf(
					"request channel = %q, want %q",
					gotRequest.Channel,
					"C123",
				)
			}

			if gotRequest.Text != tt.wantText {
				t.Errorf(
					"request text = %q, want %q",
					gotRequest.Text,
					tt.wantText,
				)
			}

			if !tt.wantBlocks {
				if len(gotRequest.Blocks) != 0 {
					t.Errorf(
						"request blocks length = %d, want 0",
						len(gotRequest.Blocks),
					)
				}

				return
			}

			if len(gotRequest.Blocks) != 2 {
				t.Fatalf(
					"request blocks length = %d, want 2",
					len(gotRequest.Blocks),
				)
			}

			section := gotRequest.Blocks[0]

			if section.Type != "section" {
				t.Errorf(
					"section type = %q, want %q",
					section.Type,
					"section",
				)
			}

			if section.Text == nil {
				t.Fatal("section text = nil")
			}

			if section.Text.Type != "mrkdwn" {
				t.Errorf(
					"section text type = %q, want %q",
					section.Text.Type,
					"mrkdwn",
				)
			}

			if section.Text.Text != tt.wantText {
				t.Errorf(
					"section text = %q, want %q",
					section.Text.Text,
					tt.wantText,
				)
			}

			actions := gotRequest.Blocks[1]

			if actions.Type != "actions" {
				t.Errorf(
					"actions type = %q, want %q",
					actions.Type,
					"actions",
				)
			}

			if len(actions.Elements) != 1 {
				t.Fatalf(
					"actions elements length = %d, want 1",
					len(actions.Elements),
				)
			}

			button := actions.Elements[0]

			if button.Type != "button" {
				t.Errorf(
					"button type = %q, want %q",
					button.Type,
					"button",
				)
			}

			if button.Text.Type != "plain_text" {
				t.Errorf(
					"button text type = %q, want %q",
					button.Text.Type,
					"plain_text",
				)
			}

			if button.Text.Text != tt.wantButton.Text {
				t.Errorf(
					"button text = %q, want %q",
					button.Text.Text,
					tt.wantButton.Text,
				)
			}

			if button.URL != tt.wantButton.URL {
				t.Errorf(
					"button URL = %q, want %q",
					button.URL,
					tt.wantButton.URL,
				)
			}
		})
	}
}
