package session

import (
	"context"
	"testing"
)

func TestDynamoDBStoreDeleteRejectsInvalidSessionIDHash(
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

				err := store.Delete(
					context.Background(),
					test.sessionIDHash,
				)
				if err == nil {
					t.Fatal(
						"Delete() error = nil, want an error",
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
