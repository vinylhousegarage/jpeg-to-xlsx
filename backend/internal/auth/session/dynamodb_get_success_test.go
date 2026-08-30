package session

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDynamoDBStoreGet(
	t *testing.T,
) {
	t.Parallel()

	wantSession := newTestSession()

	item, err := attributevalue.MarshalMap(
		struct {
			SessionIDHash string `dynamodbav:"session_id_hash"`
			CognitoSub    string `dynamodbav:"cognito_sub"`
			CreatedAt     int64  `dynamodbav:"created_at"`
			ExpiresAt     int64  `dynamodbav:"expires_at"`
		}{
			SessionIDHash: wantSession.IDHash,
			CognitoSub:    wantSession.CognitoSub,
			CreatedAt: wantSession.
				CreatedAt.
				Unix(),
			ExpiresAt: wantSession.
				ExpiresAt.
				Unix(),
		},
	)
	if err != nil {
		t.Fatalf(
			"marshal GetItem output: %v",
			err,
		)
	}

	client := &fakeDynamoDBClient{
		getItemOutput: &dynamodb.GetItemOutput{
			Item: item,
		},
	}

	store := newTestDynamoDBStore(
		t,
		client,
	)

	gotSession, err := store.Get(
		context.Background(),
		testSessionIDHash,
	)
	if err != nil {
		t.Fatalf(
			"Get() error = %v",
			err,
		)
	}

	if gotSession != wantSession {
		t.Errorf(
			"Get() session = %+v, want %+v",
			gotSession,
			wantSession,
		)
	}

	if client.getItemCalls != 1 {
		t.Fatalf(
			"GetItem() calls = %d, want 1",
			client.getItemCalls,
		)
	}

	input := client.getItemInput
	if input == nil {
		t.Fatal(
			"GetItem() input = nil",
		)
	}

	if got := aws.ToString(
		input.TableName,
	); got != testSessionTableName {
		t.Errorf(
			"GetItem() table name = %q, want %q",
			got,
			testSessionTableName,
		)
	}

	if !aws.ToBool(input.ConsistentRead) {
		t.Error(
			"GetItem() ConsistentRead = false, want true",
		)
	}

	if len(input.Key) != 1 {
		t.Errorf(
			"GetItem() key attribute count = %d, want 1",
			len(input.Key),
		)
	}

	keyAttribute, exists :=
		input.Key["session_id_hash"]
	if !exists {
		t.Fatal(
			"GetItem() key session_id_hash is missing",
		)
	}

	sessionIDHashAttribute, ok :=
		keyAttribute.(*types.AttributeValueMemberS)
	if !ok {
		t.Fatalf(
			"GetItem() key session_id_hash type = %T, want *types.AttributeValueMemberS",
			keyAttribute,
		)
	}

	if sessionIDHashAttribute.Value !=
		testSessionIDHash {
		t.Errorf(
			"GetItem() session ID hash = %q, want %q",
			sessionIDHashAttribute.Value,
			testSessionIDHash,
		)
	}
}
