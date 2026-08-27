package secretsmanager

import (
	"context"
	"errors"
	"strings"
	"testing"

	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func TestLoaderLoadClientSecretGetError(t *testing.T) {
	t.Parallel()

	getError := errors.New("get secret failed")

	client := &mockClient{
		getSecretValueFunc: func(
			_ context.Context,
			_ *awssecretsmanager.GetSecretValueInput,
			_ ...func(*awssecretsmanager.Options),
		) (*awssecretsmanager.GetSecretValueOutput, error) {
			return nil, getError
		},
	}

	loader := NewLoader(client)

	_, err := loader.LoadClientSecret(
		context.Background(),
		"test-secret-arn",
	)
	if err == nil {
		t.Fatal(
			"LoadClientSecret() error = nil, want error",
		)
	}

	if !errors.Is(err, getError) {
		t.Errorf(
			"LoadClientSecret() error does not wrap %v: %v",
			getError,
			err,
		)
	}

	const wantError = "load client secret: get secret value"

	if !strings.Contains(err.Error(), wantError) {
		t.Errorf(
			"LoadClientSecret() error = %q, want containing %q",
			err.Error(),
			wantError,
		)
	}
}

func TestLoaderLoadClientSecretResponseErrors(t *testing.T) {
	t.Parallel()

	emptyString := ""
	whitespaceString := " "
	invalidJSON := "{"
	missingClientSecret := `{"other":"value"}`
	emptyClientSecret := `{"client_secret":""}`
	whitespaceClientSecret := `{"client_secret":" "}`

	tests := []struct {
		name      string
		output    *awssecretsmanager.GetSecretValueOutput
		wantError string
	}{
		{
			name:      "nil output",
			output:    nil,
			wantError: "load client secret: output is nil",
		},
		{
			name: "nil secret string",
			output: &awssecretsmanager.GetSecretValueOutput{
				SecretString: nil,
			},
			wantError: "load client secret: secret string is empty",
		},
		{
			name: "empty secret string",
			output: &awssecretsmanager.GetSecretValueOutput{
				SecretString: &emptyString,
			},
			wantError: "load client secret: secret string is empty",
		},
		{
			name: "whitespace secret string",
			output: &awssecretsmanager.GetSecretValueOutput{
				SecretString: &whitespaceString,
			},
			wantError: "load client secret: secret string is empty",
		},
		{
			name: "invalid JSON",
			output: &awssecretsmanager.GetSecretValueOutput{
				SecretString: &invalidJSON,
			},
			wantError: "load client secret: decode secret",
		},
		{
			name: "missing client secret",
			output: &awssecretsmanager.GetSecretValueOutput{
				SecretString: &missingClientSecret,
			},
			wantError: "load client secret: client_secret is empty",
		},
		{
			name: "empty client secret",
			output: &awssecretsmanager.GetSecretValueOutput{
				SecretString: &emptyClientSecret,
			},
			wantError: "load client secret: client_secret is empty",
		},
		{
			name: "whitespace client secret",
			output: &awssecretsmanager.GetSecretValueOutput{
				SecretString: &whitespaceClientSecret,
			},
			wantError: "load client secret: client_secret is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &mockClient{
				getSecretValueFunc: func(
					_ context.Context,
					_ *awssecretsmanager.GetSecretValueInput,
					_ ...func(*awssecretsmanager.Options),
				) (*awssecretsmanager.GetSecretValueOutput, error) {
					return tt.output, nil
				},
			}

			loader := NewLoader(client)

			_, err := loader.LoadClientSecret(
				context.Background(),
				"test-secret-arn",
			)
			if err == nil {
				t.Fatal(
					"LoadClientSecret() error = nil, want error",
				)
			}

			if !strings.Contains(
				err.Error(),
				tt.wantError,
			) {
				t.Errorf(
					"LoadClientSecret() error = %q, want containing %q",
					err.Error(),
					tt.wantError,
				)
			}
		})
	}
}
