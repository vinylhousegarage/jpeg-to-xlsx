package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	testSessionTableName = "jpeg-to-xlsx-staging-auth-sessions"

	testSessionIDHash = "0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef"

	testCognitoSub = "11111111-2222-3333-4444-555555555555"
)

var (
	errPutSession = errors.New(
		"put session failed",
	)
	errGetSession = errors.New(
		"get session failed",
	)
	errDeleteSession = errors.New(
		"delete session failed",
	)
)

type fakeDynamoDBClient struct {
	putItemCalls  int
	putItemInput  *dynamodb.PutItemInput
	putItemOutput *dynamodb.PutItemOutput
	putItemErr    error

	getItemCalls  int
	getItemInput  *dynamodb.GetItemInput
	getItemOutput *dynamodb.GetItemOutput
	getItemErr    error

	deleteItemCalls  int
	deleteItemInput  *dynamodb.DeleteItemInput
	deleteItemOutput *dynamodb.DeleteItemOutput
	deleteItemErr    error
}

func (
	client *fakeDynamoDBClient,
) PutItem(
	_ context.Context,
	input *dynamodb.PutItemInput,
	_ ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	client.putItemCalls++
	client.putItemInput = input

	if client.putItemErr != nil {
		return nil, client.putItemErr
	}

	if client.putItemOutput != nil {
		return client.putItemOutput, nil
	}

	return &dynamodb.PutItemOutput{}, nil
}

func (
	client *fakeDynamoDBClient,
) GetItem(
	_ context.Context,
	input *dynamodb.GetItemInput,
	_ ...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	client.getItemCalls++
	client.getItemInput = input

	if client.getItemErr != nil {
		return nil, client.getItemErr
	}

	if client.getItemOutput != nil {
		return client.getItemOutput, nil
	}

	return &dynamodb.GetItemOutput{}, nil
}

func (
	client *fakeDynamoDBClient,
) DeleteItem(
	_ context.Context,
	input *dynamodb.DeleteItemInput,
	_ ...func(*dynamodb.Options),
) (*dynamodb.DeleteItemOutput, error) {
	client.deleteItemCalls++
	client.deleteItemInput = input

	if client.deleteItemErr != nil {
		return nil, client.deleteItemErr
	}

	if client.deleteItemOutput != nil {
		return client.deleteItemOutput, nil
	}

	return &dynamodb.DeleteItemOutput{}, nil
}

func newTestDynamoDBStore(
	t *testing.T,
	client *fakeDynamoDBClient,
) *DynamoDBStore {
	t.Helper()

	store, err := NewDynamoDBStore(
		client,
		testSessionTableName,
	)
	if err != nil {
		t.Fatalf(
			"NewDynamoDBStore() error = %v",
			err,
		)
	}

	store.now = func() time.Time {
		return testSessionNow()
	}

	return store
}

func newTestSession() Session {
	now := testSessionNow()

	return Session{
		IDHash:     testSessionIDHash,
		CognitoSub: testCognitoSub,
		CreatedAt:  now.Add(-time.Hour),
		ExpiresAt:  now.Add(7 * time.Hour),
	}
}

func testSessionNow() time.Time {
	return time.Date(
		2026,
		time.August,
		30,
		12,
		34,
		56,
		0,
		time.UTC,
	)
}
