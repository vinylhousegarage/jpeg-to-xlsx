package oauthstate

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const testOAuthStateTableName = "test-cognito-oauth-states"

var errTestDynamoDB = errors.New(
	"test DynamoDB error",
)

type fakeDynamoDBClient struct {
	putItemInput  *dynamodb.PutItemInput
	putItemOutput *dynamodb.PutItemOutput
	putItemErr    error
	putItemCalls  int

	deleteItemInput  *dynamodb.DeleteItemInput
	deleteItemOutput *dynamodb.DeleteItemOutput
	deleteItemErr    error
	deleteItemCalls  int
}

func (
	client *fakeDynamoDBClient,
) PutItem(
	_ context.Context,
	input *dynamodb.PutItemInput,
	_ ...func(*dynamodb.Options),
) (
	*dynamodb.PutItemOutput,
	error,
) {
	client.putItemCalls++
	client.putItemInput = input

	if client.putItemErr != nil {
		return nil, client.putItemErr
	}

	if client.putItemOutput == nil {
		return &dynamodb.PutItemOutput{}, nil
	}

	return client.putItemOutput, nil
}

func (
	client *fakeDynamoDBClient,
) DeleteItem(
	_ context.Context,
	input *dynamodb.DeleteItemInput,
	_ ...func(*dynamodb.Options),
) (
	*dynamodb.DeleteItemOutput,
	error,
) {
	client.deleteItemCalls++
	client.deleteItemInput = input

	if client.deleteItemErr != nil {
		return nil, client.deleteItemErr
	}

	if client.deleteItemOutput == nil {
		return &dynamodb.DeleteItemOutput{}, nil
	}

	return client.deleteItemOutput, nil
}

func newTestState() State {
	createdAt := time.Date(
		2026,
		time.August,
		30,
		12,
		34,
		56,
		0,
		time.UTC,
	)

	return State{
		Value:        "test-state",
		CodeVerifier: "test-code-verifier",
		Nonce:        "test-nonce",
		CreatedAt:    createdAt,
		ExpiresAt: createdAt.Add(
			10 * time.Minute,
		),
	}
}
