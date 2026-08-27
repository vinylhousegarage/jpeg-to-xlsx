package tokenstore

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const defaultTokenID = "default"

type dynamodbAPI interface {
	GetItem(
		ctx context.Context,
		params *dynamodb.GetItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)

	PutItem(
		ctx context.Context,
		params *dynamodb.PutItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
}

// Store persists Slack OAuth tokens in DynamoDB.
type Store struct {
	client    dynamodbAPI
	tableName string
	now       func() time.Time
}

func NewStore(
	client dynamodbAPI,
	tableName string,
) *Store {
	return &Store{
		client:    client,
		tableName: tableName,
		now:       time.Now,
	}
}

type tokenItem struct {
	ID          string `dynamodbav:"id"`
	TeamID      string `dynamodbav:"team_id"`
	AccessToken string `dynamodbav:"access_token"`
	BotUserID   string `dynamodbav:"bot_user_id"`
	UserID      string `dynamodbav:"user_id"`
	ChannelID   string `dynamodbav:"channel_id"`
	UpdatedAt   string `dynamodbav:"updated_at"`
}
