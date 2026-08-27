package tokenstore

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TestStore_Get_InvalidStoredToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		item      tokenItem
		wantError string
	}{
		{
			name: "missing stored team ID",
			item: tokenItem{
				ID:          defaultTokenID,
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
				UserID:      "U123",
				ChannelID:   "D123",
			},
			wantError: "get slack token: stored team ID is empty",
		},
		{
			name: "missing stored access token",
			item: tokenItem{
				ID:        defaultTokenID,
				TeamID:    "T123",
				BotUserID: "B123",
				UserID:    "U123",
				ChannelID: "D123",
			},
			wantError: "get slack token: stored access token is empty",
		},
		{
			name: "missing stored bot user ID",
			item: tokenItem{
				ID:          defaultTokenID,
				TeamID:      "T123",
				AccessToken: "xoxb-test",
				UserID:      "U123",
				ChannelID:   "D123",
			},
			wantError: "get slack token: stored bot user ID is empty",
		},
		{
			name: "missing stored user ID",
			item: tokenItem{
				ID:          defaultTokenID,
				TeamID:      "T123",
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
				ChannelID:   "D123",
			},
			wantError: "get slack token: stored user ID is empty",
		},
		{
			name: "missing stored channel ID",
			item: tokenItem{
				ID:          defaultTokenID,
				TeamID:      "T123",
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
				UserID:      "U123",
			},
			wantError: "get slack token: stored channel ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			attributes, err := attributevalue.MarshalMap(tt.item)
			if err != nil {
				t.Fatalf(
					"failed to marshal test token item: %v",
					err,
				)
			}

			client := &stubDynamoDBClient{
				getItemOutput: &dynamodb.GetItemOutput{
					Item: attributes,
				},
			}

			store := NewStore(
				client,
				testTableName,
			)

			got, err := store.Get(
				context.Background(),
			)
			if err == nil {
				t.Fatal(
					"Get() error = nil, want an error",
				)
			}

			if got != nil {
				t.Errorf(
					"Get() token = %#v, want nil",
					got,
				)
			}

			if err.Error() != tt.wantError {
				t.Errorf(
					"Get() error = %q, want %q",
					err.Error(),
					tt.wantError,
				)
			}

			if !client.getItemCalled {
				t.Fatal(
					"GetItem() was not called",
				)
			}

			if client.putItemCalled {
				t.Error(
					"PutItem() was called by Get()",
				)
			}
		})
	}
}
