package tokenstore

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

func TestStore_Save_PutItemError(t *testing.T) {
	t.Parallel()

	putErr := errors.New("dynamodb unavailable")

	client := &stubDynamoDBClient{
		putItemErr: putErr,
	}

	store := NewStore(client, testTableName)
	store.now = func() time.Time {
		return time.Date(
			2026,
			time.August,
			2,
			12,
			0,
			0,
			0,
			time.UTC,
		)
	}

	token := &oauth.Token{
		TeamID:      "T123",
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
		UserID:      "U123",
		ChannelID:   "D123",
	}

	err := store.Save(context.Background(), testCognitoSub, token)
	if err == nil {
		t.Fatal("Save() error = nil, want an error")
	}

	if !errors.Is(err, putErr) {
		t.Errorf("Save() error = %v, want wrapped error %v", err, putErr)
	}

	const wantError = "put slack token item: dynamodb unavailable"

	if got := err.Error(); got != wantError {
		t.Errorf("Save() error = %q, want %q", got, wantError)
	}

	if !client.putItemCalled {
		t.Fatal("PutItem() was not called")
	}

	if client.putItemInput == nil {
		t.Fatal("PutItem() input is nil")
	}

	cognitoSubAttribute, ok := client.putItemInput.Item["cognito_sub"]
	if !ok {
		t.Fatal(`PutItem() item does not contain "cognito_sub"`)
	}

	cognitoSub, ok := cognitoSubAttribute.(*types.AttributeValueMemberS)
	if !ok {
		t.Fatalf(
			`PutItem() item "cognito_sub" type = %T, want *types.AttributeValueMemberS`,
			cognitoSubAttribute,
		)
	}

	if cognitoSub.Value != testCognitoSub {
		t.Errorf(`PutItem() item "cognito_sub" = %q, want %q`, cognitoSub.Value, testCognitoSub)
	}

	if client.getItemCalled {
		t.Error("GetItem() was called by Save()")
	}
}
