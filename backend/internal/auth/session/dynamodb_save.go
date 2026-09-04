package session

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (
	store *DynamoDBStore,
) Save(
	ctx context.Context,
	session Session,
) error {
	if err := store.validateSessionForSave(
		session,
	); err != nil {
		return fmt.Errorf(
			"save session: %w",
			err,
		)
	}

	item, err := attributevalue.MarshalMap(
		sessionItem{
			IDHash:     session.IDHash,
			CognitoSub: session.CognitoSub,
			CreatedAt: session.
				CreatedAt.
				Unix(),
			ExpiresAt: session.
				ExpiresAt.
				Unix(),
		},
	)
	if err != nil {
		return fmt.Errorf(
			"save session: marshal item: %w",
			err,
		)
	}

	_, err = store.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String(
				store.tableName,
			),
			Item: item,
			ConditionExpression: aws.String(
				"attribute_not_exists(#id_hash)",
			),
			ExpressionAttributeNames: map[string]string{
				"#id_hash": "id_hash",
			},
		},
	)
	if err != nil {
		return fmt.Errorf(
			"save session: put item: %w",
			err,
		)
	}

	return nil
}

func (
	store *DynamoDBStore,
) validateSessionForSave(
	session Session,
) error {
	if strings.TrimSpace(
		session.IDHash,
	) == "" {
		return fmt.Errorf(
			"session ID hash is empty",
		)
	}

	if strings.TrimSpace(
		session.CognitoSub,
	) == "" {
		return fmt.Errorf(
			"cognito sub is empty",
		)
	}

	if session.CreatedAt.IsZero() {
		return fmt.Errorf(
			"created at is zero",
		)
	}

	if session.ExpiresAt.IsZero() {
		return fmt.Errorf(
			"expires at is zero",
		)
	}

	if !session.ExpiresAt.After(
		session.CreatedAt,
	) {
		return fmt.Errorf(
			"expires at must be after created at",
		)
	}

	if !session.ExpiresAt.After(
		store.now(),
	) {
		return fmt.Errorf(
			"session is already expired",
		)
	}

	return nil
}
