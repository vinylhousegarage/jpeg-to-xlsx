package notifier

import (
	"context"
	"errors"
	"fmt"

	slackapi "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/api"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

type tokenStore interface {
	Get(
		ctx context.Context,
		cognitoSub string,
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
	cognitoSub string,
	message Message,
) error {
	if cognitoSub == "" {
		return errors.New("notify slack: cognito sub is empty")
	}

	if message.ShotNumber == "" {
		return errors.New("notify slack: shot number is empty")
	}

	if message.DownloadURL == "" {
		return errors.New("notify slack: download URL is empty")
	}

	token, err := n.tokenStore.Get(ctx, cognitoSub)
	if err != nil {
		return fmt.Errorf("get slack token: %w", err)
	}

	if token == nil {
		return errors.New("get slack token: token is nil")
	}

	if token.AccessToken == "" {
		return errors.New("get slack token: access token is empty")
	}

	if token.ChannelID == "" {
		return errors.New("get slack token: channel ID is empty")
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
