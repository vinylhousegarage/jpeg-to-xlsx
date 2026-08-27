package tokenstore

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

func (s *Store) Save(
	ctx context.Context,
	token *oauth.Token,
) error {
	if token == nil {
		return fmt.Errorf("save slack token: token is nil")
	}

	if token.TeamID == "" {
		return fmt.Errorf("save slack token: team ID is empty")
	}

	if token.AccessToken == "" {
		return fmt.Errorf("save slack token: access token is empty")
	}

	if token.BotUserID == "" {
		return fmt.Errorf("save slack token: bot user ID is empty")
	}

	if token.UserID == "" {
		return fmt.Errorf("save slack token: user ID is empty")
	}

	if token.ChannelID == "" {
		return fmt.Errorf("save slack token: channel ID is empty")
	}

	item := tokenItem{
		ID:          defaultTokenID,
		TeamID:      token.TeamID,
		AccessToken: token.AccessToken,
		BotUserID:   token.BotUserID,
		UserID:      token.UserID,
		ChannelID:   token.ChannelID,
		UpdatedAt:   s.now().UTC().Format(time.RFC3339),
	}

	attributes, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf(
			"marshal slack token item: %w",
			err,
		)
	}

	_, err = s.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: &s.tableName,
			Item:      attributes,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"put slack token item: %w",
			err,
		)
	}

	return nil
}
