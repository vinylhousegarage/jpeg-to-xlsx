package session

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dynamodbClient interface {
	PutItem(
		ctx context.Context,
		input *dynamodb.PutItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)

	GetItem(
		ctx context.Context,
		input *dynamodb.GetItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)

	DeleteItem(
		ctx context.Context,
		input *dynamodb.DeleteItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.DeleteItemOutput, error)
}

type DynamoDBStore struct {
	client    dynamodbClient
	tableName string
	now       func() time.Time
}

type sessionItem struct {
	IDHash     string `dynamodbav:"id_hash"`
	CognitoSub string `dynamodbav:"cognito_sub"`
	CreatedAt  int64  `dynamodbav:"created_at"`
	ExpiresAt  int64  `dynamodbav:"expires_at"`
}

var _ Store = (*DynamoDBStore)(nil)

func NewDynamoDBStore(
	client dynamodbClient,
	tableName string,
) (*DynamoDBStore, error) {
	if client == nil {
		return nil, fmt.Errorf(
			"create session store: DynamoDB client is nil",
		)
	}

	if strings.TrimSpace(tableName) == "" {
		return nil, fmt.Errorf(
			"create session store: table name is empty",
		)
	}

	return &DynamoDBStore{
		client:    client,
		tableName: tableName,
		now:       time.Now,
	}, nil
}
