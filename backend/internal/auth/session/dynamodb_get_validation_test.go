package session

import (
	"context"
	"testing"
)

func TestDynamoDBStoreGetRejectsInvalidSessionIDHash(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name          string
		sessionIDHash string
	}{
		{
			name:          "empty",
			sessionIDHash: "",
		},
		{
			name:          "whitespace",
			sessionIDHash: "   ",
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

				session, err := store.Get(
					context.Background(),
					test.sessionIDHash,
				)
				if err == nil {
					t.Fatal(
						"Get() error = nil, want an error",
					)
				}

				if session != (Session{}) {
					t.Errorf(
						"Get() session = %+v, want zero value",
						session,
					)
				}

				if client.getItemCalls != 0 {
					t.Errorf(
						"GetItem() calls = %d, want 0",
						client.getItemCalls,
					)
				}
			},
		)
	}
}
