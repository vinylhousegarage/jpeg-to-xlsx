package notifier

import (
	"context"
	"errors"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

func TestNotifier_Notify_GetTokenError(t *testing.T) {
	t.Parallel()

	getErr := errors.New("dynamodb unavailable")

	tokenStore := &stubTokenStore{
		err: getErr,
	}
	client := &stubMessageClient{}

	notifier := NewNotifier(tokenStore, client)

	err := notifier.Notify(
		context.Background(),
		Message{
			ShotNumber:  "001",
			DownloadURL: "https://example.com/test.json",
		},
	)
	if err == nil {
		t.Fatal("Notify() error = nil, want an error")
	}

	if !errors.Is(err, getErr) {
		t.Errorf(
			"Notify() error = %v, want wrapped error %v",
			err,
			getErr,
		)
	}

	const wantError = "get slack token: dynamodb unavailable"
	if err.Error() != wantError {
		t.Errorf(
			"Notify() error = %q, want %q",
			err.Error(),
			wantError,
		)
	}

	if !tokenStore.called {
		t.Fatal("Get() was not called")
	}

	if client.called {
		t.Error("PostMessage() was called after Get() failed")
	}
}

func TestNotifier_Notify_PostMessageError(t *testing.T) {
	t.Parallel()

	postErr := errors.New("slack unavailable")

	tokenStore := &stubTokenStore{
		token: &oauth.Token{
			TeamID:      "T123",
			AccessToken: "xoxb-test",
			BotUserID:   "B123",
			ChannelID:   "C123",
		},
	}
	client := &stubMessageClient{
		err: postErr,
	}

	notifier := NewNotifier(tokenStore, client)

	err := notifier.Notify(
		context.Background(),
		Message{
			ShotNumber:  "001",
			DownloadURL: "https://example.com/test.json",
		},
	)
	if err == nil {
		t.Fatal("Notify() error = nil, want an error")
	}

	if !errors.Is(err, postErr) {
		t.Errorf(
			"Notify() error = %v, want wrapped error %v",
			err,
			postErr,
		)
	}

	const wantError = "post slack message: slack unavailable"
	if err.Error() != wantError {
		t.Errorf(
			"Notify() error = %q, want %q",
			err.Error(),
			wantError,
		)
	}

	if !tokenStore.called {
		t.Fatal("Get() was not called")
	}

	if !client.called {
		t.Fatal("PostMessage() was not called")
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
