package oauthstate

import "testing"

func TestNewDynamoDBStoreRejectsNilClient(
	t *testing.T,
) {
	t.Parallel()

	store, err := NewDynamoDBStore(
		nil,
		testOAuthStateTableName,
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

func TestNewDynamoDBStoreRejectsEmptyTableName(
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
			tableName: " \t\n",
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
