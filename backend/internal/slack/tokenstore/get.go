package tokenstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

var ErrTokenNotFound = errors.New("slack token not found")

func (s *Store) Get(
	ctx context.Context,
) (*oauth.Token, error) {
	output, err := s.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: &s.tableName,
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{
					Value: defaultTokenID,
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get slack token item: %w",
			err,
		)
	}

	if output == nil || len(output.Item) == 0 {
		return nil, fmt.Errorf(
			"get slack token: %w",
			ErrTokenNotFound,
		)
	}

	var item tokenItem
	if err := attributevalue.UnmarshalMap(
		output.Item,
		&item,
	); err != nil {
		return nil, fmt.Errorf(
			"unmarshal slack token item: %w",
			err,
		)
	}

	if item.TeamID == "" {
		return nil, fmt.Errorf(
			"get slack token: stored team ID is empty",
		)
	}

	if item.AccessToken == "" {
		return nil, fmt.Errorf(
			"get slack token: stored access token is empty",
		)
	}

	if item.BotUserID == "" {
		return nil, fmt.Errorf(
			"get slack token: stored bot user ID is empty",
		)
	}

	if item.UserID == "" {
		return nil, fmt.Errorf(
			"get slack token: stored user ID is empty",
		)
	}

	if item.ChannelID == "" {
		return nil, fmt.Errorf(
			"get slack token: stored channel ID is empty",
		)
	}

	return &oauth.Token{
		TeamID:      item.TeamID,
		AccessToken: item.AccessToken,
		BotUserID:   item.BotUserID,
		UserID:      item.UserID,
		ChannelID:   item.ChannelID,
	}, nil
}
