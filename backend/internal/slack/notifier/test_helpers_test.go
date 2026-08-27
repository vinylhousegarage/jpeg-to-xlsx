package notifier

import (
	"context"

	slackapi "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/api"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

type stubTokenStore struct {
	token *oauth.Token
	err   error

	called bool
	ctx    context.Context
}

func (s *stubTokenStore) Get(
	ctx context.Context,
) (*oauth.Token, error) {
	s.called = true
	s.ctx = ctx

	return s.token, s.err
}

type stubMessageClient struct {
	err error

	called      bool
	ctx         context.Context
	accessToken string
	channelID   string
	message     slackapi.Message
}

func (s *stubMessageClient) PostMessage(
	ctx context.Context,
	accessToken string,
	channelID string,
	message slackapi.Message,
) error {
	s.called = true
	s.ctx = ctx
	s.accessToken = accessToken
	s.channelID = channelID
	s.message = message

	return s.err
}
