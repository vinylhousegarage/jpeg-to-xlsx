package tokenstore

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestStore_Get_NotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		output *dynamodb.GetItemOutput
	}{
		{
			name:   "nil output",
			output: nil,
		},
		{
			name: "nil item",
			output: &dynamodb.GetItemOutput{
				Item: nil,
			},
		},
		{
			name: "empty item",
			output: &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &stubDynamoDBClient{
				getItemOutput: tt.output,
			}

			store := NewStore(
				client,
				testTableName,
			)

			got, err := store.Get(
				context.Background(),
				testCognitoSub,
			)
			if err == nil {
				t.Fatal("Get() error = nil, want an error")
			}

			if got != nil {
				t.Errorf("Get() token = %#v, want nil", got)
			}

			if !errors.Is(err, ErrTokenNotFound) {
				t.Errorf("Get() error = %v, want ErrTokenNotFound", err)
			}

			const wantError = "get slack token: slack token not found"

			if err.Error() != wantError {
				t.Errorf("Get() error = %q, want %q", err.Error(), wantError)
			}

			if !client.getItemCalled {
				t.Fatal("GetItem() was not called")
			}

			if client.putItemCalled {
				t.Error("PutItem() was called by Get()")
			}
		})
	}
}

func TestStore_Get_GetItemError(t *testing.T) {
	t.Parallel()

	getErr := errors.New("dynamodb unavailable")

	client := &stubDynamoDBClient{
		getItemErr: getErr,
	}

	store := NewStore(client, testTableName)

	got, err := store.Get(context.Background(), testCognitoSub)
	if err == nil {
		t.Fatal("Get() error = nil, want an error")
	}

	if got != nil {
		t.Errorf("Get() token = %#v, want nil", got)
	}

	if !errors.Is(err, getErr) {
		t.Errorf("Get() error = %v, want wrapped error %v", err, getErr)
	}

	const wantError = "get slack token item: dynamodb unavailable"

	if err.Error() != wantError {
		t.Errorf("Get() error = %q, want %q", err.Error(), wantError)
	}

	if !client.getItemCalled {
		t.Fatal("GetItem() was not called")
	}

	if client.putItemCalled {
		t.Error("PutItem() was called by Get()")
	}
}

func TestStore_Get_UnmarshalError(t *testing.T) {
	t.Parallel()

	client := &stubDynamoDBClient{
		getItemOutput: &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"cognito_sub": &types.AttributeValueMemberS{
					Value: testCognitoSub,
				},
				"team_id": &types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberS{
							Value: "T123",
						},
					},
				},
				"access_token": &types.AttributeValueMemberS{
					Value: "xoxb-test",
				},
				"bot_user_id": &types.AttributeValueMemberS{
					Value: "B123",
				},
				"user_id": &types.AttributeValueMemberS{
					Value: "U123",
				},
				"channel_id": &types.AttributeValueMemberS{
					Value: "D123",
				},
			},
		},
	}

	store := NewStore(client, testTableName)

	got, err := store.Get(context.Background(), testCognitoSub)
	if err == nil {
		t.Fatal("Get() error = nil, want an error")
	}

	if got != nil {
		t.Errorf("Get() token = %#v, want nil", got)
	}

	const wantPrefix = "unmarshal slack token item:"

	if !strings.HasPrefix(err.Error(), wantPrefix) {
		t.Errorf("Get() error = %q, want prefix %q", err.Error(), wantPrefix)
	}

	if !client.getItemCalled {
		t.Fatal("GetItem() was not called")
	}

	if client.putItemCalled {
		t.Error("PutItem() was called by Get()")
	}
}
