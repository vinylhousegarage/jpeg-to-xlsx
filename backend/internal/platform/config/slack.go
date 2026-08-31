package config

import (
	"fmt"
	"os"
	"strings"
)

func loadSlackConfig() (
	SlackConfig,
	error,
) {
	clientID, err := loadRequiredEnv(
		"SLACK_CLIENT_ID",
	)
	if err != nil {
		return SlackConfig{}, err
	}

	clientSecret := os.Getenv(
		"SLACK_CLIENT_SECRET",
	)
	clientSecretARN := os.Getenv(
		"SLACK_CLIENT_SECRET_ARN",
	)

	if strings.TrimSpace(clientSecret) == "" &&
		strings.TrimSpace(clientSecretARN) == "" {
		return SlackConfig{}, fmt.Errorf(
			"SLACK_CLIENT_SECRET or SLACK_CLIENT_SECRET_ARN is required",
		)
	}

	redirectURI, err := loadRequiredEnv(
		"SLACK_REDIRECT_URI",
	)
	if err != nil {
		return SlackConfig{}, err
	}

	tokenTableName, err :=
		loadSlackTokenTableName()
	if err != nil {
		return SlackConfig{}, err
	}

	return SlackConfig{
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		ClientSecretARN: clientSecretARN,
		RedirectURI:     redirectURI,
		TokenTableName:  tokenTableName,
	}, nil
}

func loadSlackTokenConfig() (
	SlackTokenConfig,
	error,
) {
	tokenTableName, err :=
		loadSlackTokenTableName()
	if err != nil {
		return SlackTokenConfig{}, err
	}

	return SlackTokenConfig{
		TokenTableName: tokenTableName,
	}, nil
}

func loadSlackTokenTableName() (
	string,
	error,
) {
	return loadRequiredEnv(
		"SLACK_TOKEN_TABLE_NAME",
	)
}
