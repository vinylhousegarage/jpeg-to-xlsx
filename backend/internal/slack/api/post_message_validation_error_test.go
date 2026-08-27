package api

import (
	"context"
	"testing"
)

func TestClient_PostMessage_ValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accessToken string
		channelID   string
		message     Message
		wantError   string
	}{
		{
			name:      "missing access token",
			channelID: "C123",
			message: Message{
				Text: "test message",
			},
			wantError: "post slack message: access token is empty",
		},
		{
			name:        "missing channel ID",
			accessToken: "xoxb-test",
			message: Message{
				Text: "test message",
			},
			wantError: "post slack message: channel ID is empty",
		},
		{
			name:        "missing text",
			accessToken: "xoxb-test",
			channelID:   "C123",
			message:     Message{},
			wantError:   "post slack message: text is empty",
		},
		{
			name:        "missing button text",
			accessToken: "xoxb-test",
			channelID:   "C123",
			message: Message{
				Text: "撮影番号：001",
				Button: &Button{
					URL: "https://example.com/test.json",
				},
			},
			wantError: "post slack message: button text is empty",
		},
		{
			name:        "missing button URL",
			accessToken: "xoxb-test",
			channelID:   "C123",
			message: Message{
				Text: "撮影番号：001",
				Button: &Button{
					Text: "ダウンロード",
				},
			},
			wantError: "post slack message: button URL is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewClient(nil)

			err := client.PostMessage(
				context.Background(),
				tt.accessToken,
				tt.channelID,
				tt.message,
			)
			if err == nil {
				t.Fatal("PostMessage() error = nil, want an error")
			}

			if err.Error() != tt.wantError {
				t.Errorf(
					"PostMessage() error = %q, want %q",
					err.Error(),
					tt.wantError,
				)
			}
		})
	}
}
