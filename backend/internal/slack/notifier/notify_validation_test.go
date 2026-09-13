package notifier

import (
	"context"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

func TestNotifier_Notify_ValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		cognitoSub string
		message    Message
		wantError  string
	}{
		{
			name:       "missing cognito sub",
			cognitoSub: "",
			message: Message{
				ShotNumber:  "001",
				DownloadURL: "https://example.com/test.xlsx",
			},
			wantError: "notify slack: cognito sub is empty",
		},
		{
			name:       "missing shot number",
			cognitoSub: testNotifierCognitoSub,
			message: Message{
				DownloadURL: "https://example.com/test.xlsx",
			},
			wantError: "notify slack: shot number is empty",
		},
		{
			name:       "missing download URL",
			cognitoSub: testNotifierCognitoSub,
			message: Message{
				ShotNumber: "001",
			},
			wantError: "notify slack: download URL is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tokenStore := &stubTokenStore{}
			client := &stubMessageClient{}

			notifier := NewNotifier(tokenStore, client)

			err := notifier.Notify(context.Background(), tt.cognitoSub, tt.message)
			if err == nil {
				t.Fatal("Notify() error = nil, want an error")
			}

			if err.Error() != tt.wantError {
				t.Errorf("Notify() error = %q, want %q", err.Error(), tt.wantError)
			}

			if tokenStore.called {
				t.Error("Get() was called for invalid input")
			}

			if client.called {
				t.Error("PostMessage() was called for invalid input")
			}
		})
	}
}

func TestNotifier_Notify_InvalidToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		token     *oauth.Token
		wantError string
	}{
		{
			name:      "nil token",
			token:     nil,
			wantError: "get slack token: token is nil",
		},
		{
			name: "missing access token",
			token: &oauth.Token{
				TeamID:    "T123",
				BotUserID: "B123",
				ChannelID: "C123",
			},
			wantError: "get slack token: access token is empty",
		},
		{
			name: "missing channel ID",
			token: &oauth.Token{
				TeamID:      "T123",
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
			},
			wantError: "get slack token: channel ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tokenStore := &stubTokenStore{
				token: tt.token,
			}
			client := &stubMessageClient{}

			notifier := NewNotifier(tokenStore, client)

			err := notifier.Notify(
				context.Background(),
				testNotifierCognitoSub,
				Message{
					ShotNumber:  "001",
					DownloadURL: "https://example.com/test.xlsx",
				},
			)
			if err == nil {
				t.Fatal("Notify() error = nil, want an error")
			}

			if err.Error() != tt.wantError {
				t.Errorf("Notify() error = %q, want %q", err.Error(), tt.wantError)
			}

			if !tokenStore.called {
				t.Fatal("Get() was not called")
			}

			if tokenStore.cognitoSub != testNotifierCognitoSub {
				t.Errorf(
					"Get() cognitoSub = %q, want %q",
					tokenStore.cognitoSub,
					testNotifierCognitoSub,
				)
			}

			if client.called {
				t.Error("PostMessage() was called with an invalid token")
			}
		})
	}
}
