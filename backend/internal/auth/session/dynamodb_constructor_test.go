package session

import (
	"testing"
)

func TestNewDynamoDBStore(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{}

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

	if store == nil {
		t.Fatal(
			"NewDynamoDBStore() store = nil, want non-nil",
		)
	}
}

func TestNewDynamoDBStoreRejectsNilClient(
	t *testing.T,
) {
	t.Parallel()

	store, err := NewDynamoDBStore(
		nil,
		testSessionTableName,
	)
	if err == nil {
		t.Fatal(
			"NewDynamoDBStore() error = nil, want an error",
		)
	}

	if store != nil {
		t.Errorf(
			"NewDynamoDBStore() store = %v, want nil",
			store,
		)
	}
}

func TestNewDynamoDBStoreRejectsInvalidTableName(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name      string
		tableName string
	}{
		{
			name:      "empty",
			tableName: "",
		},
		{
			name:      "whitespace",
			tableName: "   ",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				store, err :=
					NewDynamoDBStore(
						&fakeDynamoDBClient{},
						test.tableName,
					)
				if err == nil {
					t.Fatal(
						"NewDynamoDBStore() error = nil, want an error",
					)
				}

				if store != nil {
					t.Errorf(
						"NewDynamoDBStore() store = %v, want nil",
						store,
					)
				}
			},
		)
	}
}
