package tokenstore

import (
	"context"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

func TestStore_Save_ValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		token     *oauth.Token
		wantError string
	}{
		{
			name:      "nil token",
			token:     nil,
			wantError: "save slack token: token is nil",
		},
		{
			name: "missing team ID",
			token: &oauth.Token{
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
				UserID:      "U123",
				ChannelID:   "D123",
			},
			wantError: "save slack token: team ID is empty",
		},
		{
			name: "missing access token",
			token: &oauth.Token{
				TeamID:    "T123",
				BotUserID: "B123",
				UserID:    "U123",
				ChannelID: "D123",
			},
			wantError: "save slack token: access token is empty",
		},
		{
			name: "missing bot user ID",
			token: &oauth.Token{
				TeamID:      "T123",
				AccessToken: "xoxb-test",
				UserID:      "U123",
				ChannelID:   "D123",
			},
			wantError: "save slack token: bot user ID is empty",
		},
		{
			name: "missing user ID",
			token: &oauth.Token{
				TeamID:      "T123",
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
				ChannelID:   "D123",
			},
			wantError: "save slack token: user ID is empty",
		},
		{
			name: "missing channel ID",
			token: &oauth.Token{
				TeamID:      "T123",
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
				UserID:      "U123",
			},
			wantError: "save slack token: channel ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &stubDynamoDBClient{}
			store := NewStore(
				client,
				testTableName,
			)

			err := store.Save(
				context.Background(),
				tt.token,
			)
			if err == nil {
				t.Fatal(
					"Save() error = nil, want an error",
				)
			}

			if err.Error() != tt.wantError {
				t.Errorf(
					"Save() error = %q, want %q",
					err.Error(),
					tt.wantError,
				)
			}

			if client.putItemCalled {
				t.Error(
					"PutItem() was called for invalid token",
				)
			}

			if client.getItemCalled {
				t.Error(
					"GetItem() was called by Save()",
				)
			}
		})
	}
}
