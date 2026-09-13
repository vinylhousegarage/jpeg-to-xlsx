package tokenstore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	testTableName  = "slack-tokens"
	testCognitoSub = "cognito-user-123"
)

type stubDynamoDBClient struct {
	putItemOutput *dynamodb.PutItemOutput
	putItemErr    error
	putItemCalled bool
	putItemInput  *dynamodb.PutItemInput

	getItemOutput *dynamodb.GetItemOutput
	getItemErr    error
	getItemCalled bool
	getItemInput  *dynamodb.GetItemInput
}

func (s *stubDynamoDBClient) PutItem(
	_ context.Context,
	input *dynamodb.PutItemInput,
	_ ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	s.putItemCalled = true
	s.putItemInput = input

	return s.putItemOutput, s.putItemErr
}

func (s *stubDynamoDBClient) GetItem(
	_ context.Context,
	input *dynamodb.GetItemInput,
	_ ...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	s.getItemCalled = true
	s.getItemInput = input

	return s.getItemOutput, s.getItemErr
}
