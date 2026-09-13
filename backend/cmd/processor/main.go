package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/bedrock"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/bedrock/prompts"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/config"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/logger"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/processor"
	slackapi "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/api"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/notifier"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/tokenstore"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/storage"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/usecase"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadProcessor()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	appLogger, err := logger.NewLogger(cfg.App.Env)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	defer func() {
		_ = appLogger.Sync()
	}()

	baseS3Client := s3.NewFromConfig(awsCfg)
	baseBedrockClient := bedrockruntime.NewFromConfig(awsCfg)
	baseDynamoDBClient := dynamodb.NewFromConfig(awsCfg)

	s3Client := storage.NewS3Client(
		cfg.Storage.InputBucketName,
		baseS3Client,
	)

	presignClient := storage.NewS3PresignClient(
		cfg.Storage.OutputBucketName,
		baseS3Client,
	)

	promptText, err := prompts.LoadPrompt(
		cfg.Bedrock.PromptFileName,
	)
	if err != nil {
		appLogger.Fatal("failed to load prompt", zap.Error(err))
	}

	bedrockClient := bedrock.NewClient(
		baseBedrockClient,
		cfg.Bedrock.ModelID,
		promptText,
		appLogger,
	)

	bedrockService := bedrock.NewService(bedrockClient)

	tokenStore := tokenstore.NewStore(
		baseDynamoDBClient,
		cfg.SlackToken.TokenTableName,
	)

	slackClient := slackapi.NewClient(nil)

	slackNotifier := notifier.NewNotifier(tokenStore, slackClient)

	workflow := usecase.NewWorkflow(
		s3Client,
		s3Client,
		presignClient,
		bedrockService,
		slackNotifier,
		cfg.Storage.OutputBucketName,
		appLogger,
	)

	handler := processor.NewHandler(workflow)

	appLogger.Info("Starting processor on AWS Lambda")

	lambda.Start(handler.HandleRequest)
}
