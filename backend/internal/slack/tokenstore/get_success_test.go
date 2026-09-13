package tokenstore

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestStore_Get_Success(t *testing.T) {
	t.Parallel()

	wantItem := tokenItem{
		CognitoSub:  testCognitoSub,
		TeamID:      "T123",
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
		UserID:      "U123",
		ChannelID:   "D123",
		UpdatedAt:   "2026-08-02T12:00:00Z",
	}

	attributes, err := attributevalue.MarshalMap(wantItem)
	if err != nil {
		t.Fatalf("failed to marshal test token item: %v", err)
	}

	client := &stubDynamoDBClient{
		getItemOutput: &dynamodb.GetItemOutput{
			Item: attributes,
		},
	}

	store := NewStore(client, testTableName)

	got, err := store.Get(context.Background(), testCognitoSub)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if !client.getItemCalled {
		t.Fatal("GetItem() was not called")
	}

	if client.getItemInput == nil {
		t.Fatal("GetItem() input is nil")
	}

	if client.getItemInput.TableName == nil {
		t.Fatal("GetItem() TableName is nil")
	}

	if gotTableName := *client.getItemInput.TableName; gotTableName != testTableName {
		t.Errorf(
			"GetItem() TableName = %q, want %q",
			gotTableName,
			testTableName,
		)
	}

	cognitoSubAttribute, ok := client.getItemInput.Key["cognito_sub"]
	if !ok {
		t.Fatal(`GetItem() key does not contain "cognito_sub"`)
	}

	cognitoSub, ok := cognitoSubAttribute.(*types.AttributeValueMemberS)
	if !ok {
		t.Fatalf(
			`GetItem() key "cognito_sub" type = %T, want *types.AttributeValueMemberS`,
			cognitoSubAttribute,
		)
	}

	if cognitoSub.Value != testCognitoSub {
		t.Errorf(
			`GetItem() key "cognito_sub" = %q, want %q`,
			cognitoSub.Value,
			testCognitoSub,
		)
	}

	if got == nil {
		t.Fatal("Get() token is nil")
	}

	if got.TeamID != wantItem.TeamID {
		t.Errorf("Get() TeamID = %q, want %q", got.TeamID, wantItem.TeamID)
	}

	if got.AccessToken != wantItem.AccessToken {
		t.Errorf("Get() AccessToken = %q, want %q", got.AccessToken, wantItem.AccessToken)
	}

	if got.BotUserID != wantItem.BotUserID {
		t.Errorf("Get() BotUserID = %q, want %q", got.BotUserID, wantItem.BotUserID)
	}

	if got.UserID != wantItem.UserID {
		t.Errorf("Get() UserID = %q, want %q", got.UserID, wantItem.UserID)
	}

	if got.ChannelID != wantItem.ChannelID {
		t.Errorf("Get() ChannelID = %q, want %q", got.ChannelID, wantItem.ChannelID)
	}

	if client.putItemCalled {
		t.Error("PutItem() was called by Get()")
	}
}
