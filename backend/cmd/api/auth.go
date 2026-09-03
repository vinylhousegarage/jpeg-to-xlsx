package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito"
	authcallback "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/callback"
	authlogin "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/login"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
	authlogout "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session/logout"
	authstatus "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session/status"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/config"
	platformsecretsmanager "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/secretsmanager"
)

type authHandlers struct {
	login         http.Handler
	callback      http.Handler
	sessionStatus http.Handler
	logout        http.Handler
}

func buildAuthHandlers(
	ctx context.Context,
	authConfig config.BFFAuthConfig,
	dynamoDBClient *dynamodb.Client,
	secretLoader *platformsecretsmanager.Loader,
) (
	authHandlers,
	error,
) {
	if ctx == nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: context is nil",
		)
	}

	if dynamoDBClient == nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: DynamoDB client is nil",
		)
	}

	if secretLoader == nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: secret loader is nil",
		)
	}

	clientSecret, err := secretLoader.LoadClientSecret(
		ctx,
		authConfig.CognitoClientSecretARN,
	)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: load Cognito client secret: %w",
			err,
		)
	}

	cognitoClient, err := cognito.NewClient(
		cognito.Config{
			ClientID:              authConfig.CognitoClientID,
			ClientSecret:          clientSecret,
			AuthorizationEndpoint: authConfig.CognitoAuthorizationEndpoint,
			TokenEndpoint:         authConfig.CognitoTokenEndpoint,
			RedirectURI:           authConfig.CognitoRedirectURI,
			Scopes: []string{
				"openid",
				"email",
				"profile",
			},
		},
	)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create Cognito client: %w",
			err,
		)
	}

	verifier, err := cognito.NewVerifier(
		ctx,
		cognito.VerifierConfig{
			Issuer:   authConfig.CognitoIssuer,
			ClientID: authConfig.CognitoClientID,
		},
	)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create ID token verifier: %w",
			err,
		)
	}

	stateGenerator, err := oauthstate.NewGenerator(authConfig.OAuthStateTTL)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create OAuth state generator: %w",
			err,
		)
	}

	stateStore, err := oauthstate.NewDynamoDBStore(
		dynamoDBClient,
		authConfig.OAuthStateTableName,
	)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create OAuth state store: %w",
			err,
		)
	}

	sessionStore, err := session.NewDynamoDBStore(
		dynamoDBClient,
		authConfig.SessionTableName,
	)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create session store: %w",
			err,
		)
	}

	cookieManager, err := session.NewCookieManager(authConfig.SessionLifetime)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create cookie manager: %w",
			err,
		)
	}

	sessionResolver, err := session.NewResolver(
		cookieManager,
		sessionStore,
	)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create session resolver: %w",
			err,
		)
	}

	sessionIDGenerator := session.NewIDGenerator()

	loginHandler, err := authlogin.NewHandler(
		stateGenerator,
		stateStore,
		cognitoClient,
	)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create login handler: %w",
			err,
		)
	}

	callbackHandler, err := authcallback.NewHandler(
		stateStore,
		cognitoClient,
		verifier,
		sessionIDGenerator,
		sessionStore,
		cookieManager,
		authcallback.Config{
			RedirectURL:     authConfig.PostLoginRedirectURL,
			SessionLifetime: authConfig.SessionLifetime,
		},
	)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create callback handler: %w",
			err,
		)
	}

	sessionStatusHandler, err := authstatus.NewHandler(sessionResolver)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create session status handler: %w",
			err,
		)
	}

	logoutHandler, err := authlogout.NewHandler(
		sessionResolver,
		sessionStore,
		cookieManager,
	)
	if err != nil {
		return authHandlers{}, fmt.Errorf(
			"build authentication handlers: create logout handler: %w",
			err,
		)
	}

	return authHandlers{
		login:         loginHandler,
		callback:      callbackHandler,
		sessionStatus: sessionStatusHandler,
		logout:        logoutHandler,
	}, nil
}
