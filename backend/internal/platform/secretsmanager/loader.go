package secretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Loader struct {
	client Client
}

type secretValue struct {
	ClientSecret string `json:"client_secret"`
}

func NewLoader(client Client) *Loader {
	return &Loader{
		client: client,
	}
}

func (l *Loader) LoadClientSecret(
	ctx context.Context,
	secretARN string,
) (string, error) {
	if strings.TrimSpace(secretARN) == "" {
		return "", fmt.Errorf(
			"load client secret: secret ARN is empty",
		)
	}

	if l == nil || l.client == nil {
		return "", fmt.Errorf(
			"load client secret: client is nil",
		)
	}

	output, err := l.client.GetSecretValue(
		ctx,
		&awssecretsmanager.GetSecretValueInput{
			SecretId: &secretARN,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"load client secret: get secret value: %w",
			err,
		)
	}

	if output == nil {
		return "", fmt.Errorf(
			"load client secret: output is nil",
		)
	}

	if output.SecretString == nil ||
		strings.TrimSpace(*output.SecretString) == "" {
		return "", fmt.Errorf(
			"load client secret: secret string is empty",
		)
	}

	var value secretValue

	if err := json.Unmarshal(
		[]byte(*output.SecretString),
		&value,
	); err != nil {
		return "", fmt.Errorf(
			"load client secret: decode secret: %w",
			err,
		)
	}

	if strings.TrimSpace(value.ClientSecret) == "" {
		return "", fmt.Errorf(
			"load client secret: client_secret is empty",
		)
	}

	return value.ClientSecret, nil
}
