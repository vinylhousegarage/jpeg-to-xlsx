package oauthstate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamodbAPI interface {
	PutItem(
		ctx context.Context,
		input *dynamodb.PutItemInput,
		optFns ...func(*dynamodb.Options),
	) (
		*dynamodb.PutItemOutput,
		error,
	)

	DeleteItem(
		ctx context.Context,
		input *dynamodb.DeleteItemInput,
		optFns ...func(*dynamodb.Options),
	) (
		*dynamodb.DeleteItemOutput,
		error,
	)
}

type DynamoDBStore struct {
	client    dynamodbAPI
	tableName string
}

type dynamodbItem struct {
	State        string `dynamodbav:"state"`
	CodeVerifier string `dynamodbav:"code_verifier"`
	Nonce        string `dynamodbav:"nonce"`
	CreatedAt    int64  `dynamodbav:"created_at"`
	ExpiresAt    int64  `dynamodbav:"expires_at"`
}

var _ Store = (*DynamoDBStore)(nil)

func NewDynamoDBStore(
	client dynamodbAPI,
	tableName string,
) (
	*DynamoDBStore,
	error,
) {
	if client == nil {
		return nil, fmt.Errorf(
			"create OAuth state store: DynamoDB client is nil",
		)
	}

	if strings.TrimSpace(tableName) == "" {
		return nil, fmt.Errorf(
			"create OAuth state store: table name is empty",
		)
	}

	return &DynamoDBStore{
		client:    client,
		tableName: tableName,
	}, nil
}

func (s *DynamoDBStore) Save(
	ctx context.Context,
	state State,
) error {
	if err := validateState(state); err != nil {
		return fmt.Errorf(
			"save OAuth state: %w",
			err,
		)
	}

	item, err := attributevalue.MarshalMap(
		dynamodbItem{
			State:        state.Value,
			CodeVerifier: state.CodeVerifier,
			Nonce:        state.Nonce,
			CreatedAt:    state.CreatedAt.Unix(),
			ExpiresAt:    state.ExpiresAt.Unix(),
		},
	)
	if err != nil {
		return fmt.Errorf(
			"save OAuth state: marshal item: %w",
			err,
		)
	}

	_, err = s.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String(
				s.tableName,
			),
			Item: item,
			ConditionExpression: aws.String(
				"attribute_not_exists(#state)",
			),
			ExpressionAttributeNames: map[string]string{
				"#state": "state",
			},
		},
	)
	if err != nil {
		return fmt.Errorf(
			"save OAuth state: put item: %w",
			err,
		)
	}

	return nil
}

func (s *DynamoDBStore) Consume(
	ctx context.Context,
	value string,
) (
	State,
	error,
) {
	if strings.TrimSpace(value) == "" {
		return State{}, fmt.Errorf(
			"consume OAuth state: state value is empty",
		)
	}

	output, err := s.client.DeleteItem(
		ctx,
		&dynamodb.DeleteItemInput{
			TableName: aws.String(
				s.tableName,
			),
			Key: map[string]dynamodbtypes.AttributeValue{
				"state": &dynamodbtypes.AttributeValueMemberS{
					Value: value,
				},
			},
			ReturnValues: dynamodbtypes.ReturnValueAllOld,
		},
	)
	if err != nil {
		return State{}, fmt.Errorf(
			"consume OAuth state: delete item: %w",
			err,
		)
	}

	if output == nil ||
		len(output.Attributes) == 0 {
		return State{}, fmt.Errorf(
			"consume OAuth state: %w",
			ErrNotFound,
		)
	}

	var item dynamodbItem

	if err := attributevalue.UnmarshalMap(
		output.Attributes,
		&item,
	); err != nil {
		return State{}, fmt.Errorf(
			"consume OAuth state: unmarshal item: %w",
			err,
		)
	}

	state := State{
		Value:        item.State,
		CodeVerifier: item.CodeVerifier,
		Nonce:        item.Nonce,
		CreatedAt: time.Unix(
			item.CreatedAt,
			0,
		).UTC(),
		ExpiresAt: time.Unix(
			item.ExpiresAt,
			0,
		).UTC(),
	}

	if err := validateState(state); err != nil {
		return State{}, fmt.Errorf(
			"consume OAuth state: invalid item: %w",
			err,
		)
	}

	return state, nil
}

func validateState(
	state State,
) error {
	if strings.TrimSpace(state.Value) == "" {
		return fmt.Errorf(
			"state value is empty",
		)
	}

	if strings.TrimSpace(
		state.CodeVerifier,
	) == "" {
		return fmt.Errorf(
			"code verifier is empty",
		)
	}

	if strings.TrimSpace(state.Nonce) == "" {
		return fmt.Errorf(
			"nonce is empty",
		)
	}

	if state.CreatedAt.IsZero() {
		return fmt.Errorf(
			"created at is zero",
		)
	}

	if state.ExpiresAt.IsZero() {
		return fmt.Errorf(
			"expires at is zero",
		)
	}

	if !state.ExpiresAt.After(
		state.CreatedAt,
	) {
		return fmt.Errorf(
			"expires at must be after created at",
		)
	}

	return nil
}
