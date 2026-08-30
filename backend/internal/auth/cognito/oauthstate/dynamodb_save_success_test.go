package oauthstate

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

func TestDynamoDBStoreSave(t *testing.T) {
	t.Parallel()

	client := &fakeDynamoDBClient{}

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

	state := newTestState()

	err = store.Save(
		context.Background(),
		state,
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
	); got != testOAuthStateTableName {
		t.Errorf(
			"PutItem() table name = %q, want %q",
			got,
			testOAuthStateTableName,
		)
	}

	if len(input.Item) != 5 {
		t.Errorf(
			"PutItem() item attribute count = %d, want 5",
			len(input.Item),
		)
	}

	var item struct {
		State        string `dynamodbav:"state"`
		CodeVerifier string `dynamodbav:"code_verifier"`
		Nonce        string `dynamodbav:"nonce"`
		CreatedAt    int64  `dynamodbav:"created_at"`
		ExpiresAt    int64  `dynamodbav:"expires_at"`
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

	if item.State != state.Value {
		t.Errorf(
			"item state = %q, want %q",
			item.State,
			state.Value,
		)
	}

	if item.CodeVerifier !=
		state.CodeVerifier {
		t.Errorf(
			"item code verifier = %q, want %q",
			item.CodeVerifier,
			state.CodeVerifier,
		)
	}

	if item.Nonce != state.Nonce {
		t.Errorf(
			"item nonce = %q, want %q",
			item.Nonce,
			state.Nonce,
		)
	}

	if item.CreatedAt !=
		state.CreatedAt.Unix() {
		t.Errorf(
			"item created at = %d, want %d",
			item.CreatedAt,
			state.CreatedAt.Unix(),
		)
	}

	if item.ExpiresAt !=
		state.ExpiresAt.Unix() {
		t.Errorf(
			"item expires at = %d, want %d",
			item.ExpiresAt,
			state.ExpiresAt.Unix(),
		)
	}

	const wantConditionExpression = "attribute_not_exists(#state)"

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
		input.ExpressionAttributeNames["#state"]

	if gotAttributeName != "state" {
		t.Errorf(
			"PutItem() expression attribute name #state = %q, want %q",
			gotAttributeName,
			"state",
		)
	}
}
