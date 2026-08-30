package oauthstate

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDynamoDBStoreConsumeRejectsEmptyValue(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "empty",
			value: "",
		},
		{
			name:  "whitespace",
			value: " \t\n",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				client :=
					&fakeDynamoDBClient{}

				store, err :=
					NewDynamoDBStore(
						client,
						testOAuthStateTableName,
					)
				if err != nil {
					t.Fatalf(
						"NewDynamoDBStore() error = %v",
						err,
					)
				}

				state, err := store.Consume(
					context.Background(),
					test.value,
				)
				if err == nil {
					t.Fatal(
						"Consume() error = nil, want an error",
					)
				}

				if state != (State{}) {
					t.Errorf(
						"Consume() state = %+v, want zero value",
						state,
					)
				}

				if client.deleteItemCalls != 0 {
					t.Errorf(
						"DeleteItem() calls = %d, want 0",
						client.deleteItemCalls,
					)
				}
			},
		)
	}
}

func TestDynamoDBStoreConsumeReturnsDynamoDBError(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{
		deleteItemErr: errTestDynamoDB,
	}

	store, err := NewDynamoDBStore(
		client,
		testOAuthStateTableName,
	)
	if err != nil {
		t.Fatalf(
			"NewDynamoDBStore() error = %v",
			err,
		)
	}

	state, err := store.Consume(
		context.Background(),
		"test-state",
	)
	if err == nil {
		t.Fatal(
			"Consume() error = nil, want an error",
		)
	}

	if state != (State{}) {
		t.Errorf(
			"Consume() state = %+v, want zero value",
			state,
		)
	}

	if !errors.Is(
		err,
		errTestDynamoDB,
	) {
		t.Errorf(
			"Consume() error = %v, want wrapped %v",
			err,
			errTestDynamoDB,
		)
	}

	if client.deleteItemCalls != 1 {
		t.Errorf(
			"DeleteItem() calls = %d, want 1",
			client.deleteItemCalls,
		)
	}
}

func TestDynamoDBStoreConsumeReturnsNotFound(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{
		deleteItemOutput: &dynamodb.DeleteItemOutput{},
	}

	store, err := NewDynamoDBStore(
		client,
		testOAuthStateTableName,
	)
	if err != nil {
		t.Fatalf(
			"NewDynamoDBStore() error = %v",
			err,
		)
	}

	state, err := store.Consume(
		context.Background(),
		"missing-state",
	)
	if err == nil {
		t.Fatal(
			"Consume() error = nil, want an error",
		)
	}

	if state != (State{}) {
		t.Errorf(
			"Consume() state = %+v, want zero value",
			state,
		)
	}

	if !errors.Is(err, ErrNotFound) {
		t.Errorf(
			"Consume() error = %v, want %v",
			err,
			ErrNotFound,
		)
	}
}

func TestDynamoDBStoreConsumeRejectsInvalidItem(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{
		deleteItemOutput: &dynamodb.DeleteItemOutput{
			Attributes: map[string]dynamodbtypes.AttributeValue{
				"state": &dynamodbtypes.AttributeValueMemberS{
					Value: "test-state",
				},
				"created_at": &dynamodbtypes.AttributeValueMemberS{
					Value: "invalid",
				},
			},
		},
	}

	store, err := NewDynamoDBStore(
		client,
		testOAuthStateTableName,
	)
	if err != nil {
		t.Fatalf(
			"NewDynamoDBStore() error = %v",
			err,
		)
	}

	state, err := store.Consume(
		context.Background(),
		"test-state",
	)
	if err == nil {
		t.Fatal(
			"Consume() error = nil, want an error",
		)
	}

	if state != (State{}) {
		t.Errorf(
			"Consume() state = %+v, want zero value",
			state,
		)
	}
}
