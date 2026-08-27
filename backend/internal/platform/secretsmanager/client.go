package secretsmanager

import (
	"context"

	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Client interface {
	GetSecretValue(
		ctx context.Context,
		params *awssecretsmanager.GetSecretValueInput,
		optFns ...func(*awssecretsmanager.Options),
	) (*awssecretsmanager.GetSecretValueOutput, error)
}
