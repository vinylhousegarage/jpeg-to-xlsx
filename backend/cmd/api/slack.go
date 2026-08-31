package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/config"
	platformsecretsmanager "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/secretsmanager"
	slackapi "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/api"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
	slackcallback "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth/callback"
	slacklogin "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth/login"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/tokenstore"
)

type slackHandlers struct {
	login    http.Handler
	callback http.Handler
}

func buildSlackHandlers(
	ctx context.Context,
	slackConfig config.SlackConfig,
	cookieSecure bool,
	dynamoDBClient *dynamodb.Client,
	secretLoader *platformsecretsmanager.Loader,
	appLogger *zap.Logger,
) (
	slackHandlers,
	error,
) {
	if ctx == nil {
		return slackHandlers{}, fmt.Errorf(
			"build Slack handlers: context is nil",
		)
	}

	if dynamoDBClient == nil {
		return slackHandlers{}, fmt.Errorf(
			"build Slack handlers: DynamoDB client is nil",
		)
	}

	if appLogger == nil {
		return slackHandlers{}, fmt.Errorf(
			"build Slack handlers: logger is nil",
		)
	}

	clientSecret := slackConfig.ClientSecret

	if strings.TrimSpace(
		slackConfig.ClientSecretARN,
	) != "" {
		if secretLoader == nil {
			return slackHandlers{}, fmt.Errorf(
				"build Slack handlers: secret loader is nil",
			)
		}

		loadedClientSecret, err :=
			secretLoader.LoadClientSecret(
				ctx,
				slackConfig.ClientSecretARN,
			)
		if err != nil {
			return slackHandlers{}, fmt.Errorf(
				"build Slack handlers: load Slack client secret: %w",
				err,
			)
		}

		clientSecret = loadedClientSecret
	}

	if strings.TrimSpace(clientSecret) == "" {
		return slackHandlers{}, fmt.Errorf(
			"build Slack handlers: client secret is empty",
		)
	}

	oauthClient := oauth.NewClient(
		http.DefaultClient,
		slackConfig.ClientID,
		clientSecret,
	)

	slackAPIClient := slackapi.NewClient(
		http.DefaultClient,
	)

	tokenStore := tokenstore.NewStore(
		dynamoDBClient,
		slackConfig.TokenTableName,
	)

	loginHandler := slacklogin.NewHandler(
		slackConfig.ClientID,
		slackConfig.RedirectURI,
		cookieSecure,
		appLogger,
	)

	callbackHandler :=
		slackcallback.NewHandler(
			slackConfig.RedirectURI,
			cookieSecure,
			oauthClient,
			slackAPIClient,
			tokenStore,
			appLogger,
		)

	return slackHandlers{
		login:    loginHandler,
		callback: callbackHandler,
	}, nil
}
