package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/config"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/logger"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/router"
	platformsecretsmanager "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/secretsmanager"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadAPI()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	appLogger, err := logger.NewLogger(cfg.App.Env)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	baseS3Client := s3.NewFromConfig(awsConfig)
	dynamoDBClient := dynamodb.NewFromConfig(awsConfig)
	secretsManagerClient := awssecretsmanager.NewFromConfig(awsConfig)
	secretLoader := platformsecretsmanager.NewLoader(secretsManagerClient)

	authHandlerSet, err := buildAuthHandlers(
		ctx,
		cfg.Auth,
		dynamoDBClient,
		secretLoader,
	)
	if err != nil {
		appLogger.Fatal("Failed to build authentication handlers", zap.Error(err))
	}

	storageHandler, err := buildStorageHandler(
		cfg.Storage,
		baseS3Client,
		authHandlerSet.sessionResolver,
		appLogger,
	)
	if err != nil {
		appLogger.Fatal("Failed to build storage handler", zap.Error(err))
	}

	slackHandlerSet, err := buildSlackHandlers(
		ctx,
		cfg.Slack,
		cfg.App.CookieSecure,
		dynamoDBClient,
		secretLoader,
		authHandlerSet.sessionResolver,
		appLogger,
	)
	if err != nil {
		appLogger.Fatal("Failed to build Slack handlers", zap.Error(err))
	}

	mux := http.NewServeMux()

	router.SetupRoutes(
		mux,
		authHandlerSet.login,
		authHandlerSet.callback,
		authHandlerSet.sessionStatus,
		authHandlerSet.logout,
		slackHandlerSet.login,
		slackHandlerSet.callback,
		storageHandler,
	)

	appLogger.Info("Application initialized successfully")

	if cfg.AWS.IsLambda {
		appLogger.Info("Starting processor on AWS Lambda")
		adapter := httpadapter.NewV2(mux)
		lambda.Start(adapter.ProxyWithContext)

		return
	}

	appLogger.Info("Starting local server on :8080")

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		appLogger.Fatal("Server failed", zap.Error(err))
	}
}
