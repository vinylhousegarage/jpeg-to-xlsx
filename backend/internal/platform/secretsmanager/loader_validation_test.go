package secretsmanager

import (
	"context"
	"testing"

	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func TestLoaderLoadClientSecretValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		loader    *Loader
		secretARN string
		wantError string
	}{
		{
			name: "empty secret ARN",
			loader: NewLoader(
				&mockClient{
					getSecretValueFunc: func(
						_ context.Context,
						_ *awssecretsmanager.GetSecretValueInput,
						_ ...func(*awssecretsmanager.Options),
					) (*awssecretsmanager.GetSecretValueOutput, error) {
						t.Fatal(
							"GetSecretValue() should not be called",
						)

						return nil, nil
					},
				},
			),
			secretARN: "",
			wantError: "load client secret: secret ARN is empty",
		},
		{
			name: "whitespace secret ARN",
			loader: NewLoader(
				&mockClient{
					getSecretValueFunc: func(
						_ context.Context,
						_ *awssecretsmanager.GetSecretValueInput,
						_ ...func(*awssecretsmanager.Options),
					) (*awssecretsmanager.GetSecretValueOutput, error) {
						t.Fatal(
							"GetSecretValue() should not be called",
						)

						return nil, nil
					},
				},
			),
			secretARN: " ",
			wantError: "load client secret: secret ARN is empty",
		},
		{
			name:      "nil client",
			loader:    NewLoader(nil),
			secretARN: "test-secret-arn",
			wantError: "load client secret: client is nil",
		},
		{
			name:      "nil loader",
			loader:    nil,
			secretARN: "test-secret-arn",
			wantError: "load client secret: client is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.loader.LoadClientSecret(
				context.Background(),
				tt.secretARN,
			)
			if err == nil {
				t.Fatal(
					"LoadClientSecret() error = nil, want error",
				)
			}

			if err.Error() != tt.wantError {
				t.Errorf(
					"LoadClientSecret() error = %q, want %q",
					err.Error(),
					tt.wantError,
				)
			}
		})
	}
}
