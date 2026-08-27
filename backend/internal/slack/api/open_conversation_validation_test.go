package api

import (
	"context"
	"testing"
)

func TestClient_OpenConversation_ValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accessToken string
		userID      string
		wantError   string
	}{
		{
			name:      "missing access token",
			userID:    "U123",
			wantError: "open slack conversation: access token is empty",
		},
		{
			name:        "missing user ID",
			accessToken: "xoxb-test",
			wantError:   "open slack conversation: user ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewClient(nil)

			channelID, err := client.OpenConversation(
				context.Background(),
				tt.accessToken,
				tt.userID,
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

			if err.Error() != tt.wantError {
				t.Errorf(
					"OpenConversation() error = %q, want %q",
					err.Error(),
					tt.wantError,
				)
			}
		})
	}
}
