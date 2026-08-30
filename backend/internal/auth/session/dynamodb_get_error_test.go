package session

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestDynamoDBStoreGetReturnsGetItemError(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{
		getItemErr: errGetSession,
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
		errGetSession,
	) {
		t.Errorf(
			"Get() error = %v, want wrapped %v",
			err,
			errGetSession,
		)
	}

	if errors.Is(
		err,
		ErrNotFound,
	) {
		t.Errorf(
			"Get() error = %v, do not want %v",
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

	if strings.Contains(
		err.Error(),
		testSessionIDHash,
	) {
		t.Error(
			"Get() error contains the session ID hash",
		)
	}

	if client.getItemCalls != 1 {
		t.Errorf(
			"GetItem() calls = %d, want 1",
			client.getItemCalls,
		)
	}
}
