package main

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

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
	cfg, err := config.LoadProcessor()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	l, err := logger.NewLogger(cfg.App.Env)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	defer func() {
		_ = l.Sync()
	}()

	baseS3Client := s3.NewFromConfig(awsCfg)
	baseBedrockClient := bedrockruntime.NewFromConfig(awsCfg)
	baseDynamoClient := dynamodb.NewFromConfig(awsCfg)

	s3Client := storage.NewS3Client(
		cfg.Storage.InputBucketName,
		baseS3Client,
	)

	presignClient := storage.NewS3PresignClient(
		cfg.Storage.InputBucketName,
		baseS3Client,
	)

	promptText, err := prompts.LoadPrompt(
		cfg.Bedrock.PromptFileName,
	)
	if err != nil {
		l.Fatal(
			"failed to load prompt",
			zap.Error(err),
		)
	}

	bedrockClient := bedrock.NewClient(
		baseBedrockClient,
		cfg.Bedrock.ModelID,
		promptText,
		l,
	)

	bedrockService := bedrock.NewService(
		bedrockClient,
	)

	tokenStore := tokenstore.NewStore(
		baseDynamoClient,
		cfg.SlackToken.TokenTableName,
	)

	slackClient := slackapi.NewClient(nil)

	slackNotifier := notifier.NewNotifier(
		tokenStore,
		slackClient,
	)

	workflow := usecase.NewWorkflow(
		s3Client,
		s3Client,
		presignClient,
		bedrockService,
		slackNotifier,
		cfg.Storage.OutputBucketName,
		l,
	)

	handler := processor.NewHandler(workflow)

	l.Info("Starting processor on AWS Lambda")

	lambda.Start(handler.HandleRequest)
}
