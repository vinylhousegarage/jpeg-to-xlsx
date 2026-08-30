package oauthstate

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDynamoDBStoreConsume(t *testing.T) {
	t.Parallel()

	want := newTestState()

	storedItem := struct {
		State        string `dynamodbav:"state"`
		CodeVerifier string `dynamodbav:"code_verifier"`
		Nonce        string `dynamodbav:"nonce"`
		CreatedAt    int64  `dynamodbav:"created_at"`
		ExpiresAt    int64  `dynamodbav:"expires_at"`
	}{
		State:        want.Value,
		CodeVerifier: want.CodeVerifier,
		Nonce:        want.Nonce,
		CreatedAt:    want.CreatedAt.Unix(),
		ExpiresAt:    want.ExpiresAt.Unix(),
	}

	attributes, err :=
		attributevalue.MarshalMap(
			storedItem,
		)
	if err != nil {
		t.Fatalf(
			"marshal test item: %v",
			err,
		)
	}

	client := &fakeDynamoDBClient{
		deleteItemOutput: &dynamodb.DeleteItemOutput{
			Attributes: attributes,
		},
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

	got, err := store.Consume(
		context.Background(),
		want.Value,
	)
	if err != nil {
		t.Fatalf(
			"Consume() error = %v",
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

	if gotTableName := aws.ToString(
		input.TableName,
	); gotTableName !=
		testOAuthStateTableName {
		t.Errorf(
			"DeleteItem() table name = %q, want %q",
			gotTableName,
			testOAuthStateTableName,
		)
	}

	if len(input.Key) != 1 {
		t.Errorf(
			"DeleteItem() key attribute count = %d, want 1",
			len(input.Key),
		)
	}

	var key struct {
		State string `dynamodbav:"state"`
	}

	if err := attributevalue.UnmarshalMap(
		input.Key,
		&key,
	); err != nil {
		t.Fatalf(
			"unmarshal DeleteItem key: %v",
			err,
		)
	}

	if key.State != want.Value {
		t.Errorf(
			"DeleteItem() key state = %q, want %q",
			key.State,
			want.Value,
		)
	}

	if input.ReturnValues !=
		dynamodbtypes.ReturnValueAllOld {
		t.Errorf(
			"DeleteItem() return values = %q, want %q",
			input.ReturnValues,
			dynamodbtypes.ReturnValueAllOld,
		)
	}

	if got.Value != want.Value {
		t.Errorf(
			"Consume() state value = %q, want %q",
			got.Value,
			want.Value,
		)
	}

	if got.CodeVerifier !=
		want.CodeVerifier {
		t.Errorf(
			"Consume() code verifier = %q, want %q",
			got.CodeVerifier,
			want.CodeVerifier,
		)
	}

	if got.Nonce != want.Nonce {
		t.Errorf(
			"Consume() nonce = %q, want %q",
			got.Nonce,
			want.Nonce,
		)
	}

	if !got.CreatedAt.Equal(
		want.CreatedAt,
	) {
		t.Errorf(
			"Consume() created at = %v, want %v",
			got.CreatedAt,
			want.CreatedAt,
		)
	}

	if !got.ExpiresAt.Equal(
		want.ExpiresAt,
	) {
		t.Errorf(
			"Consume() expires at = %v, want %v",
			got.ExpiresAt,
			want.ExpiresAt,
		)
	}
}
