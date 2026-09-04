package session

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (
	store *DynamoDBStore,
) Delete(
	ctx context.Context,
	idHash string,
) error {
	if strings.TrimSpace(
		idHash,
	) == "" {
		return fmt.Errorf(
			"delete session: session ID hash is empty",
		)
	}

	_, err := store.client.DeleteItem(
		ctx,
		&dynamodb.DeleteItemInput{
			TableName: aws.String(
				store.tableName,
			),
			Key: map[string]types.AttributeValue{
				"id_hash": &types.AttributeValueMemberS{
					Value: idHash,
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf(
			"delete session: delete item: %w",
			err,
		)
	}

	return nil
}
