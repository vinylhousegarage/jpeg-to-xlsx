package secretsmanager

import (
	"context"

	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type mockClient struct {
	getSecretValueFunc func(
		ctx context.Context,
		params *awssecretsmanager.GetSecretValueInput,
		optFns ...func(*awssecretsmanager.Options),
	) (*awssecretsmanager.GetSecretValueOutput, error)
}

func (m *mockClient) GetSecretValue(
	ctx context.Context,
	params *awssecretsmanager.GetSecretValueInput,
	optFns ...func(*awssecretsmanager.Options),
) (*awssecretsmanager.GetSecretValueOutput, error) {
	return m.getSecretValueFunc(ctx, params, optFns...)
}
