package secretsmanager

import (
	"context"
	"testing"

	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func TestLoaderLoadClientSecret(t *testing.T) {
	t.Parallel()

	const (
		secretARN        = "arn:aws:secretsmanager:ap-northeast-1:123456789012:secret:test"
		wantClientSecret = "test-client-secret"
	)

	secretString := `{"client_secret":"test-client-secret"}`

	client := &mockClient{
		getSecretValueFunc: func(
			_ context.Context,
			params *awssecretsmanager.GetSecretValueInput,
			_ ...func(*awssecretsmanager.Options),
		) (*awssecretsmanager.GetSecretValueOutput, error) {
			if params == nil {
				t.Fatal("GetSecretValueInput is nil")
			}

			if params.SecretId == nil {
				t.Fatal("SecretId is nil")
			}

			if *params.SecretId != secretARN {
				t.Errorf(
					"SecretId = %q, want %q",
					*params.SecretId,
					secretARN,
				)
			}

			return &awssecretsmanager.GetSecretValueOutput{
				SecretString: &secretString,
			}, nil
		},
	}

	loader := NewLoader(client)

	got, err := loader.LoadClientSecret(
		context.Background(),
		secretARN,
	)
	if err != nil {
		t.Fatalf(
			"LoadClientSecret() error = %v",
			err,
		)
	}

	if got != wantClientSecret {
		t.Errorf(
			"LoadClientSecret() = %q, want %q",
			got,
			wantClientSecret,
		)
	}
}
