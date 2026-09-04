package session

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDynamoDBStoreGetRejectsInvalidStoredSession(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(
			map[string]types.AttributeValue,
		)
	}{
		{
			name: "missing session ID hash",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				delete(
					item,
					"id_hash",
				)
			},
		},
		{
			name: "empty session ID hash",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				item["id_hash"] =
					&types.AttributeValueMemberS{
						Value: "",
					}
			},
		},
		{
			name: "whitespace session ID hash",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				item["id_hash"] =
					&types.AttributeValueMemberS{
						Value: "   ",
					}
			},
		},
		{
			name: "session ID hash mismatch",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				item["id_hash"] =
					&types.AttributeValueMemberS{
						Value: "unexpected-session-id-hash",
					}
			},
		},
		{
			name: "missing Cognito sub",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				delete(
					item,
					"cognito_sub",
				)
			},
		},
		{
			name: "empty Cognito sub",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				item["cognito_sub"] =
					&types.AttributeValueMemberS{
						Value: "",
					}
			},
		},
		{
			name: "whitespace Cognito sub",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				item["cognito_sub"] =
					&types.AttributeValueMemberS{
						Value: "   ",
					}
			},
		},
		{
			name: "missing created at",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				delete(
					item,
					"created_at",
				)
			},
		},
		{
			name: "invalid created at type",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				item["created_at"] =
					&types.AttributeValueMemberS{
						Value: "not-a-number",
					}
			},
		},
		{
			name: "missing expires at",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				delete(
					item,
					"expires_at",
				)
			},
		},
		{
			name: "invalid expires at type",
			mutate: func(
				item map[string]types.AttributeValue,
			) {
				item["expires_at"] =
					&types.AttributeValueMemberS{
						Value: "not-a-number",
					}
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				item, err :=
					marshalTestSessionItem(
						newTestSession(),
					)
				if err != nil {
					t.Fatalf(
						"marshal GetItem output: %v",
						err,
					)
				}

				test.mutate(item)

				client :=
					&fakeDynamoDBClient{
						getItemOutput: &dynamodb.GetItemOutput{
							Item: item,
						},
					}

				store :=
					newTestDynamoDBStore(
						t,
						client,
					)

				session, err := store.Get(
					context.Background(),
					testSessionIDHash,
				)
				if err == nil {
					t.Fatal(
						"Get() error = nil, want an error",
					)
				}

				if errors.Is(
					err,
					ErrNotFound,
				) {
					t.Errorf(
						"Get() error = %v, do not want %v",
						err,
						ErrNotFound,
					)
				}

				if session != (Session{}) {
					t.Errorf(
						"Get() session = %+v, want zero value",
						session,
					)
				}

				if strings.Contains(
					err.Error(),
					testSessionIDHash,
				) {
					t.Error(
						"Get() error contains the session ID hash",
					)
				}

				if client.getItemCalls != 1 {
					t.Errorf(
						"GetItem() calls = %d, want 1",
						client.getItemCalls,
					)
				}
			},
		)
	}
}
