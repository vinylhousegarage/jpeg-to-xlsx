package session

import (
	"context"
	"testing"
	"time"
)

func TestDynamoDBStoreSaveRejectsInvalidSession(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*Session)
	}{
		{
			name: "empty session ID hash",
			mutate: func(
				session *Session,
			) {
				session.IDHash = ""
			},
		},
		{
			name: "whitespace session ID hash",
			mutate: func(
				session *Session,
			) {
				session.IDHash = "   "
			},
		},
		{
			name: "empty Cognito sub",
			mutate: func(
				session *Session,
			) {
				session.CognitoSub = ""
			},
		},
		{
			name: "whitespace Cognito sub",
			mutate: func(
				session *Session,
			) {
				session.CognitoSub = "   "
			},
		},
		{
			name: "zero created at",
			mutate: func(
				session *Session,
			) {
				session.CreatedAt =
					time.Time{}
			},
		},
		{
			name: "zero expires at",
			mutate: func(
				session *Session,
			) {
				session.ExpiresAt =
					time.Time{}
			},
		},
		{
			name: "expires at equals created at",
			mutate: func(
				session *Session,
			) {
				session.ExpiresAt =
					session.CreatedAt
			},
		},
		{
			name: "expires before created at",
			mutate: func(
				session *Session,
			) {
				session.ExpiresAt =
					session.CreatedAt.Add(
						-time.Minute,
					)
			},
		},
		{
			name: "already expired",
			mutate: func(
				session *Session,
			) {
				session.CreatedAt =
					testSessionNow().
						Add(-2 * time.Hour)
				session.ExpiresAt =
					testSessionNow().
						Add(-time.Hour)
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				client :=
					&fakeDynamoDBClient{}

				store :=
					newTestDynamoDBStore(
						t,
						client,
					)

				session := newTestSession()
				test.mutate(&session)

				err := store.Save(
					context.Background(),
					session,
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
