package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awssecretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/config"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/logger"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/router"
	platformsecretsmanager "github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/secretsmanager"
	slackapi "github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/api"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth/callback"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth/login"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/tokenstore"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/storage"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/storage/put"
)

func main() {
	// 起動処理用のルートコンテキスト
	ctx := context.Background()

	// 設定の初期化
	cfg, err := config.LoadAPI()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// AWS 設定の読み込み
	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	// Slack Client Secret の解決
	slackClientSecret := cfg.Slack.ClientSecret

	if cfg.Slack.ClientSecretARN != "" {
		secretsManagerClient :=
			awssecretsmanager.NewFromConfig(awsCfg)

		secretLoader :=
			platformsecretsmanager.NewLoader(
				secretsManagerClient,
			)

		slackClientSecret, err =
			secretLoader.LoadClientSecret(
				ctx,
				cfg.Slack.ClientSecretARN,
			)
		if err != nil {
			log.Fatalf(
				"failed to load Slack client secret: %v",
				err,
			)
		}
	}

	// logger の初期化
	l, err := logger.NewLogger(cfg.App.Env)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	l.Info("Application successfully initialized")

	// S3 クライアント
	baseS3Client := s3.NewFromConfig(awsCfg)

	// DynamoDB クライアント
	dynamoDBClient := dynamodb.NewFromConfig(awsCfg)

	// Presign 用依存
	presignClient := storage.NewS3PresignClient(
		cfg.Storage.InputBucketName,
		baseS3Client,
	)

	// Put 用 Service
	putService := put.NewService(presignClient)

	// Presign Handler
	presignHandler := put.NewHandler(
		putService,
		l,
	)

	// Slack OAuth Client
	oauthClient := oauth.NewClient(
		http.DefaultClient,
		cfg.Slack.ClientID,
		slackClientSecret,
	)

	// Slack API Client
	slackAPIClient := slackapi.NewClient(
		http.DefaultClient,
	)

	// Slack Token Store
	tokenStore := tokenstore.NewStore(
		dynamoDBClient,
		cfg.Slack.TokenTableName,
	)

	// Slack Login Handler
	slackLoginHandler := login.NewHandler(
		cfg.Slack.ClientID,
		cfg.Slack.RedirectURI,
		cfg.App.CookieSecure,
		l,
	)

	// Slack Callback Handler
	slackCallbackHandler := callback.NewHandler(
		cfg.Slack.RedirectURI,
		cfg.App.CookieSecure,
		oauthClient,
		slackAPIClient,
		tokenStore,
		l,
	)

	// サーバーの初期化
	mux := http.NewServeMux()

	// ルーティングの設定を委譲
	router.SetupRoutes(
		mux,
		presignHandler,
		slackLoginHandler,
		slackCallbackHandler,
	)

	// サーバーの起動
	if cfg.AWS.IsLambda {
		l.Info("Starting server on AWS Lambda")

		adapter := httpadapter.NewV2(mux)
		lambda.Start(adapter.ProxyWithContext)
	} else {
		l.Info("Starting local server on :8080")

		srv := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}

		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}
}
