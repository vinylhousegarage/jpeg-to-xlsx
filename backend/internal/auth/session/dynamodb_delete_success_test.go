package session

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDynamoDBStoreDelete(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{}

	store := newTestDynamoDBStore(
		t,
		client,
	)

	err := store.Delete(
		context.Background(),
		testSessionIDHash,
	)
	if err != nil {
		t.Fatalf(
			"Delete() error = %v",
			err,
		)
	}

	if client.deleteItemCalls != 1 {
		t.Fatalf(
			"DeleteItem() calls = %d, want 1",
			client.deleteItemCalls,
		)
	}

	input := client.deleteItemInput
	if input == nil {
		t.Fatal(
			"DeleteItem() input = nil",
		)
	}

	if got := aws.ToString(
		input.TableName,
	); got != testSessionTableName {
		t.Errorf(
			"DeleteItem() table name = %q, want %q",
			got,
			testSessionTableName,
		)
	}

	if len(input.Key) != 1 {
		t.Errorf(
			"DeleteItem() key attribute count = %d, want 1",
			len(input.Key),
		)
	}

	keyAttribute, exists :=
		input.Key["session_id_hash"]
	if !exists {
		t.Fatal(
			"DeleteItem() key session_id_hash is missing",
		)
	}

	sessionIDHashAttribute, ok :=
		keyAttribute.(*types.AttributeValueMemberS)
	if !ok {
		t.Fatalf(
			"DeleteItem() key session_id_hash type = %T, want *types.AttributeValueMemberS",
			keyAttribute,
		)
	}

	if sessionIDHashAttribute.Value !=
		testSessionIDHash {
		t.Errorf(
			"DeleteItem() session ID hash = %q, want %q",
			sessionIDHashAttribute.Value,
			testSessionIDHash,
		)
	}
}
