package session

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDynamoDBStoreSaveReturnsPutItemError(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{
		putItemErr: errPutSession,
	}

	store := newTestDynamoDBStore(
		t,
		client,
	)

	err := store.Save(
		context.Background(),
		newTestSession(),
	)
	if err == nil {
		t.Fatal(
			"Save() error = nil, want an error",
		)
	}

	if !errors.Is(
		err,
		errPutSession,
	) {
		t.Errorf(
			"Save() error = %v, want wrapped %v",
			err,
			errPutSession,
		)
	}

	if client.putItemCalls != 1 {
		t.Errorf(
			"PutItem() calls = %d, want 1",
			client.putItemCalls,
		)
	}
}

func TestDynamoDBStoreSaveReturnsConditionalCheckError(
	t *testing.T,
) {
	t.Parallel()

	conditionalError :=
		&types.ConditionalCheckFailedException{
			Message: stringPointer(
				"session already exists",
			),
		}

	client := &fakeDynamoDBClient{
		putItemErr: conditionalError,
	}

	store := newTestDynamoDBStore(
		t,
		client,
	)

	err := store.Save(
		context.Background(),
		newTestSession(),
	)
	if err == nil {
		t.Fatal(
			"Save() error = nil, want an error",
		)
	}

	if !errors.Is(
		err,
		conditionalError,
	) {
		t.Errorf(
			"Save() error = %v, want wrapped conditional check error",
			err,
		)
	}

	if strings.Contains(
		err.Error(),
		testSessionIDHash,
	) {
		t.Error(
			"Save() error contains the session ID hash",
		)
	}

	if client.putItemCalls != 1 {
		t.Errorf(
			"PutItem() calls = %d, want 1",
			client.putItemCalls,
		)
	}
}

func stringPointer(
	value string,
) *string {
	return &value
}
