package notifier

import (
	"context"
	"fmt"

	slackapi "github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/api"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

type tokenStore interface {
	Get(
		ctx context.Context,
	) (*oauth.Token, error)
}

type messageClient interface {
	PostMessage(
		ctx context.Context,
		accessToken string,
		channelID string,
		message slackapi.Message,
	) error
}

// Notifier sends Slack messages using stored OAuth tokens.
type Notifier struct {
	tokenStore tokenStore
	client     messageClient
}

func NewNotifier(
	tokenStore tokenStore,
	client messageClient,
) *Notifier {
	return &Notifier{
		tokenStore: tokenStore,
		client:     client,
	}
}

func (n *Notifier) Notify(
	ctx context.Context,
	message Message,
) error {
	if message.ShotNumber == "" {
		return fmt.Errorf("notify slack: shot number is empty")
	}

	if message.DownloadURL == "" {
		return fmt.Errorf("notify slack: download URL is empty")
	}

	token, err := n.tokenStore.Get(ctx)
	if err != nil {
		return fmt.Errorf("get slack token: %w", err)
	}

	if token == nil {
		return fmt.Errorf("get slack token: token is nil")
	}

	if token.AccessToken == "" {
		return fmt.Errorf("get slack token: access token is empty")
	}

	if token.ChannelID == "" {
		return fmt.Errorf("get slack token: channel ID is empty")
	}

	slackMessage := slackapi.Message{
		Text: fmt.Sprintf(
			"撮影番号：%s",
			message.ShotNumber,
		),
		Button: &slackapi.Button{
			Text: "ダウンロード",
			URL:  message.DownloadURL,
		},
	}

	if err := n.client.PostMessage(
		ctx,
		token.AccessToken,
		token.ChannelID,
		slackMessage,
	); err != nil {
		return fmt.Errorf("post slack message: %w", err)
	}

	return nil
}
