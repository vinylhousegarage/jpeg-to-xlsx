package session

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestDynamoDBStoreDeleteReturnsDeleteItemError(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{
		deleteItemErr: errDeleteSession,
	}

	store := newTestDynamoDBStore(
		t,
		client,
	)

	err := store.Delete(
		context.Background(),
		testSessionIDHash,
	)
	if err == nil {
		t.Fatal(
			"Delete() error = nil, want an error",
		)
	}

	if !errors.Is(
		err,
		errDeleteSession,
	) {
		t.Errorf(
			"Delete() error = %v, want wrapped %v",
			err,
			errDeleteSession,
		)
	}

	if strings.Contains(
		err.Error(),
		testSessionIDHash,
	) {
		t.Error(
			"Delete() error contains the session ID hash",
		)
	}

	if client.deleteItemCalls != 1 {
		t.Errorf(
			"DeleteItem() calls = %d, want 1",
			client.deleteItemCalls,
		)
	}
}
