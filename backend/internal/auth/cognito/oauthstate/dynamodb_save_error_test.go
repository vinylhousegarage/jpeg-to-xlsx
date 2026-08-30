package oauthstate

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDynamoDBStoreSaveRejectsInvalidState(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*State)
	}{
		{
			name: "empty value",
			mutate: func(state *State) {
				state.Value = ""
			},
		},
		{
			name: "whitespace value",
			mutate: func(state *State) {
				state.Value = " \t\n"
			},
		},
		{
			name: "empty code verifier",
			mutate: func(state *State) {
				state.CodeVerifier = ""
			},
		},
		{
			name: "whitespace code verifier",
			mutate: func(state *State) {
				state.CodeVerifier = " \t\n"
			},
		},
		{
			name: "empty nonce",
			mutate: func(state *State) {
				state.Nonce = ""
			},
		},
		{
			name: "whitespace nonce",
			mutate: func(state *State) {
				state.Nonce = " \t\n"
			},
		},
		{
			name: "zero created at",
			mutate: func(state *State) {
				state.CreatedAt = time.Time{}
			},
		},
		{
			name: "zero expires at",
			mutate: func(state *State) {
				state.ExpiresAt = time.Time{}
			},
		},
		{
			name: "expires at equals created at",
			mutate: func(state *State) {
				state.ExpiresAt =
					state.CreatedAt
			},
		},
		{
			name: "expires before created at",
			mutate: func(state *State) {
				state.ExpiresAt =
					state.CreatedAt.Add(
						-time.Second,
					)
			},
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

				state := newTestState()
				test.mutate(&state)

				err = store.Save(
					context.Background(),
					state,
				)
				if err == nil {
					t.Fatal(
						"Save() error = nil, want an error",
					)
				}

				if client.putItemCalls != 0 {
					t.Errorf(
						"PutItem() calls = %d, want 0",
						client.putItemCalls,
					)
				}
			},
		)
	}
}

func TestDynamoDBStoreSaveReturnsDynamoDBError(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{
		putItemErr: errTestDynamoDB,
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

	err = store.Save(
		context.Background(),
		newTestState(),
	)
	if err == nil {
		t.Fatal(
			"Save() error = nil, want an error",
		)
	}

	if !errors.Is(
		err,
		errTestDynamoDB,
	) {
		t.Errorf(
			"Save() error = %v, want wrapped %v",
			err,
			errTestDynamoDB,
		)
	}

	if client.putItemCalls != 1 {
		t.Errorf(
			"PutItem() calls = %d, want 1",
			client.putItemCalls,
		)
	}
}
