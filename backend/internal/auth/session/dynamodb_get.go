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
	idHash string,
) (Session, error) {
	if strings.TrimSpace(
		idHash,
	) == "" {
		return Session{}, fmt.Errorf(
			"get session: ID hash is empty",
		)
	}

	output, err := store.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(
				store.tableName,
			),
			Key: map[string]types.AttributeValue{
				"id_hash": &types.AttributeValueMemberS{
					Value: idHash,
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
		idHash,
	); err != nil {
		return Session{}, fmt.Errorf(
			"get session: invalid stored item: %w",
			err,
		)
	}

	session := Session{
		IDHash:     item.IDHash,
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
	expectedIDHash string,
) error {
	if strings.TrimSpace(
		item.IDHash,
	) == "" {
		return fmt.Errorf(
			"session ID hash is empty",
		)
	}

	if subtle.ConstantTimeCompare(
		[]byte(item.IDHash),
		[]byte(expectedIDHash),
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
