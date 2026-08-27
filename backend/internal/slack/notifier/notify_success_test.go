package notifier

import (
	"context"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

func TestNotifier_Notify_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tokenStore := &stubTokenStore{
		token: &oauth.Token{
			TeamID:      "T123",
			AccessToken: "xoxb-test",
			BotUserID:   "B123",
			ChannelID:   "C123",
		},
	}
	client := &stubMessageClient{}

	notifier := NewNotifier(tokenStore, client)

	err := notifier.Notify(
		ctx,
		Message{
			ShotNumber:  "001",
			DownloadURL: "https://example.com/test.json",
		},
	)
	if err != nil {
		t.Fatalf("Notify() error = %v", err)
	}

	if !tokenStore.called {
		t.Fatal("Get() was not called")
	}

	if tokenStore.ctx != ctx {
		t.Error("Get() received an unexpected context")
	}

	if !client.called {
		t.Fatal("PostMessage() was not called")
	}

	if client.ctx != ctx {
		t.Error("PostMessage() received an unexpected context")
	}

	if client.accessToken != "xoxb-test" {
		t.Errorf(
			"PostMessage() accessToken = %q, want %q",
			client.accessToken,
			"xoxb-test",
		)
	}

	if client.channelID != "C123" {
		t.Errorf(
			"PostMessage() channelID = %q, want %q",
			client.channelID,
			"C123",
		)
	}

	const wantText = "撮影番号：001"
	if client.message.Text != wantText {
		t.Errorf(
			"PostMessage() message.Text = %q, want %q",
			client.message.Text,
			wantText,
		)
	}

	if client.message.Button == nil {
		t.Fatal("PostMessage() message.Button = nil")
	}

	const wantButtonText = "ダウンロード"
	if client.message.Button.Text != wantButtonText {
		t.Errorf(
			"PostMessage() button.Text = %q, want %q",
			client.message.Button.Text,
			wantButtonText,
		)
	}

	const wantButtonURL = "https://example.com/test.json"
	if client.message.Button.URL != wantButtonURL {
		t.Errorf(
			"PostMessage() button.URL = %q, want %q",
			client.message.Button.URL,
			wantButtonURL,
		)
	}
}
