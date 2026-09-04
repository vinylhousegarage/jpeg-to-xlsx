package session

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

func TestDynamoDBStoreSave(
	t *testing.T,
) {
	t.Parallel()

	client := &fakeDynamoDBClient{}

	store := newTestDynamoDBStore(
		t,
		client,
	)

	session := newTestSession()

	err := store.Save(
		context.Background(),
		session,
	)
	if err != nil {
		t.Fatalf(
			"Save() error = %v",
			err,
		)
	}

	if client.putItemCalls != 1 {
		t.Fatalf(
			"PutItem() calls = %d, want 1",
			client.putItemCalls,
		)
	}

	input := client.putItemInput
	if input == nil {
		t.Fatal(
			"PutItem() input = nil",
		)
	}

	if got := aws.ToString(
		input.TableName,
	); got != testSessionTableName {
		t.Errorf(
			"PutItem() table name = %q, want %q",
			got,
			testSessionTableName,
		)
	}

	if len(input.Item) != 4 {
		t.Errorf(
			"PutItem() item attribute count = %d, want 4",
			len(input.Item),
		)
	}

	var item struct {
		IDHash     string `dynamodbav:"id_hash"`
		CognitoSub string `dynamodbav:"cognito_sub"`
		CreatedAt  int64  `dynamodbav:"created_at"`
		ExpiresAt  int64  `dynamodbav:"expires_at"`
	}

	if err := attributevalue.UnmarshalMap(
		input.Item,
		&item,
	); err != nil {
		t.Fatalf(
			"unmarshal PutItem item: %v",
			err,
		)
	}

	if item.IDHash !=
		session.IDHash {
		t.Errorf(
			"item session ID hash = %q, want %q",
			item.IDHash,
			session.IDHash,
		)
	}

	if item.CognitoSub !=
		session.CognitoSub {
		t.Errorf(
			"item Cognito sub = %q, want %q",
			item.CognitoSub,
			session.CognitoSub,
		)
	}

	if item.CreatedAt !=
		session.CreatedAt.Unix() {
		t.Errorf(
			"item created at = %d, want %d",
			item.CreatedAt,
			session.CreatedAt.Unix(),
		)
	}

	if item.ExpiresAt !=
		session.ExpiresAt.Unix() {
		t.Errorf(
			"item expires at = %d, want %d",
			item.ExpiresAt,
			session.ExpiresAt.Unix(),
		)
	}

	const wantConditionExpression = "attribute_not_exists(#id_hash)"

	if got := aws.ToString(
		input.ConditionExpression,
	); got != wantConditionExpression {
		t.Errorf(
			"PutItem() condition expression = %q, want %q",
			got,
			wantConditionExpression,
		)
	}

	gotAttributeName :=
		input.ExpressionAttributeNames["#id_hash"]

	if gotAttributeName !=
		"id_hash" {
		t.Errorf(
			"PutItem() expression attribute name #id_hash = %q, want %q",
			gotAttributeName,
			"id_hash",
		)
	}
}
