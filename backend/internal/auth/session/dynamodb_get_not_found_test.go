package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDynamoDBStoreGetReturnsNotFoundWhenItemDoesNotExist(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{
		getItemOutput: &dynamodb.GetItemOutput{},
	}

	store := newTestDynamoDBStore(
		t,
		client,
	)

	session, err := store.Get(
		context.Background(),
		testSessionIDHash,
	)
	if err == nil {
		t.Fatal(
			"Get() error = nil, want an error",
		)
	}

	if !errors.Is(
		err,
		ErrNotFound,
	) {
		t.Errorf(
			"Get() error = %v, want %v",
			err,
			ErrNotFound,
		)
	}

	if session != (Session{}) {
		t.Errorf(
			"Get() session = %+v, want zero value",
			session,
		)
	}

	if client.getItemCalls != 1 {
		t.Errorf(
			"GetItem() calls = %d, want 1",
			client.getItemCalls,
		)
	}
}

func TestDynamoDBStoreGetReturnsNotFoundWhenSessionExpired(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name      string
		expiresAt time.Time
	}{
		{
			name: "expired before current time",
			expiresAt: testSessionNow().
				Add(-time.Minute),
		},
		{
			name:      "expired at current time",
			expiresAt: testSessionNow(),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				storedSession :=
					newTestSession()

				storedSession.ExpiresAt =
					test.expiresAt

				item, err :=
					marshalTestSessionItem(
						storedSession,
					)
				if err != nil {
					t.Fatalf(
						"marshal GetItem output: %v",
						err,
					)
				}

				client :=
					&fakeDynamoDBClient{
						getItemOutput: &dynamodb.GetItemOutput{
							Item: item,
						},
					}

				store :=
					newTestDynamoDBStore(
						t,
						client,
					)

				session, err := store.Get(
					context.Background(),
					testSessionIDHash,
				)
				if err == nil {
					t.Fatal(
						"Get() error = nil, want an error",
					)
				}

				if !errors.Is(
					err,
					ErrNotFound,
				) {
					t.Errorf(
						"Get() error = %v, want %v",
						err,
						ErrNotFound,
					)
				}

				if session != (Session{}) {
					t.Errorf(
						"Get() session = %+v, want zero value",
						session,
					)
				}

				if client.getItemCalls != 1 {
					t.Errorf(
						"GetItem() calls = %d, want 1",
						client.getItemCalls,
					)
				}
			},
		)
	}
}

func marshalTestSessionItem(
	session Session,
) (
	map[string]types.AttributeValue,
	error,
) {
	return attributevalue.MarshalMap(
		struct {
			SessionIDHash string `dynamodbav:"session_id_hash"`
			CognitoSub    string `dynamodbav:"cognito_sub"`
			CreatedAt     int64  `dynamodbav:"created_at"`
			ExpiresAt     int64  `dynamodbav:"expires_at"`
		}{
			SessionIDHash: session.IDHash,
			CognitoSub:    session.CognitoSub,
			CreatedAt:     session.CreatedAt.Unix(),
			ExpiresAt:     session.ExpiresAt.Unix(),
		},
	)
}
