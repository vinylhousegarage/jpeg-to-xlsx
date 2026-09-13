package tokenstore

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

func TestStore_Save_Success(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(
		2026,
		time.August,
		2,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	client := &stubDynamoDBClient{
		putItemOutput: &dynamodb.PutItemOutput{},
	}

	store := NewStore(client, testTableName)
	store.now = func() time.Time {
		return fixedTime
	}

	token := &oauth.Token{
		TeamID:      "T123",
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
		UserID:      "U123",
		ChannelID:   "D123",
	}

	err := store.Save(context.Background(), testCognitoSub, token)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if !client.putItemCalled {
		t.Fatal("PutItem() was not called")
	}

	if client.putItemInput == nil {
		t.Fatal("PutItem() input is nil")
	}

	if client.putItemInput.TableName == nil {
		t.Fatal("PutItem() TableName is nil")
	}

	if got := *client.putItemInput.TableName; got != testTableName {
		t.Errorf("PutItem() TableName = %q, want %q", got, testTableName)
	}

	var item tokenItem
	if err := attributevalue.UnmarshalMap(
		client.putItemInput.Item,
		&item,
	); err != nil {
		t.Fatalf("failed to unmarshal PutItem item: %v", err)
	}

	if item.CognitoSub != testCognitoSub {
		t.Errorf("CognitoSub = %q, want %q", item.CognitoSub, testCognitoSub)
	}

	if item.TeamID != token.TeamID {
		t.Errorf("TeamID = %q, want %q", item.TeamID, token.TeamID)
	}

	if item.AccessToken != token.AccessToken {
		t.Errorf("AccessToken = %q, want %q", item.AccessToken, token.AccessToken)
	}

	if item.BotUserID != token.BotUserID {
		t.Errorf("BotUserID = %q, want %q", item.BotUserID, token.BotUserID)
	}

	if item.UserID != token.UserID {
		t.Errorf("UserID = %q, want %q", item.UserID, token.UserID)
	}

	if item.ChannelID != token.ChannelID {
		t.Errorf("ChannelID = %q, want %q", item.ChannelID, token.ChannelID)
	}

	wantUpdatedAt := fixedTime.Format(time.RFC3339)
	if item.UpdatedAt != wantUpdatedAt {
		t.Errorf("UpdatedAt = %q, want %q", item.UpdatedAt, wantUpdatedAt)
	}

	if client.getItemCalled {
		t.Error("GetItem() was called by Save()")
	}
}
