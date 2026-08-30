package session

import (
	"context"
	"crypto/subtle"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (
	store *DynamoDBStore,
) Get(
	ctx context.Context,
	sessionIDHash string,
) (Session, error) {
	if strings.TrimSpace(
		sessionIDHash,
	) == "" {
		return Session{}, fmt.Errorf(
			"get session: session ID hash is empty",
		)
	}

	output, err := store.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(
				store.tableName,
			),
			Key: map[string]types.AttributeValue{
				"session_id_hash": &types.AttributeValueMemberS{
					Value: sessionIDHash,
				},
			},
			ConsistentRead: aws.Bool(true),
		},
	)
	if err != nil {
		return Session{}, fmt.Errorf(
			"get session: get item: %w",
			err,
		)
	}

	if output == nil ||
		len(output.Item) == 0 {
		return Session{}, fmt.Errorf(
			"get session: %w",
			ErrNotFound,
		)
	}

	var item sessionItem

	if err := attributevalue.UnmarshalMap(
		output.Item,
		&item,
	); err != nil {
		return Session{}, fmt.Errorf(
			"get session: unmarshal item: %w",
			err,
		)
	}

	if err := validateStoredSessionItem(
		item,
		sessionIDHash,
	); err != nil {
		return Session{}, fmt.Errorf(
			"get session: invalid stored item: %w",
			err,
		)
	}

	session := Session{
		IDHash:     item.SessionIDHash,
		CognitoSub: item.CognitoSub,
		CreatedAt: time.Unix(
			item.CreatedAt,
			0,
		).UTC(),
		ExpiresAt: time.Unix(
			item.ExpiresAt,
			0,
		).UTC(),
	}

	if !session.ExpiresAt.After(
		store.now(),
	) {
		return Session{}, fmt.Errorf(
			"get session: expired: %w",
			ErrNotFound,
		)
	}

	return session, nil
}

func validateStoredSessionItem(
	item sessionItem,
	expectedSessionIDHash string,
) error {
	if strings.TrimSpace(
		item.SessionIDHash,
	) == "" {
		return fmt.Errorf(
			"session ID hash is empty",
		)
	}

	if subtle.ConstantTimeCompare(
		[]byte(item.SessionIDHash),
		[]byte(expectedSessionIDHash),
	) != 1 {
		return fmt.Errorf(
			"session ID hash does not match",
		)
	}

	if strings.TrimSpace(
		item.CognitoSub,
	) == "" {
		return fmt.Errorf(
			"cognito sub is empty",
		)
	}

	if item.CreatedAt <= 0 {
		return fmt.Errorf(
			"created at is invalid",
		)
	}

	if item.ExpiresAt <= 0 {
		return fmt.Errorf(
			"expires at is invalid",
		)
	}

	if item.ExpiresAt <= item.CreatedAt {
		return fmt.Errorf(
			"expires at must be after created at",
		)
	}

	return nil
}
