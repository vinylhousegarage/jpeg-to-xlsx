package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/bedrock"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/bedrock/prompts"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/config"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/logger"
)

func main() {
	// 設定の初期化
	cfg, err := config.LoadBedrockPlayground()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// loggerの初期化
	l, err := logger.NewLogger(cfg.App.Env)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer func() {
		_ = l.Sync()
	}()

	ctx := context.Background()

	// AWS設定の読み込み
	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		l.Fatal(
			"failed to load AWS config",
			zap.Error(err),
		)
	}

	// プロンプトの読み込み
	promptText, err := prompts.LoadPrompt(
		cfg.Bedrock.PromptFileName,
	)
	if err != nil {
		l.Fatal(
			"failed to load prompt",
			zap.Error(err),
		)
	}

	// AWS SDK Bedrock Runtime Client
	baseBedrockClient := bedrockruntime.NewFromConfig(awsCfg)

	// アプリケーション用Bedrock Client
	bedrockClient := bedrock.NewClient(
		baseBedrockClient,
		cfg.Bedrock.ModelID,
		promptText,
		l,
	)

	// Bedrock Service
	bedrockService := bedrock.NewService(
		bedrockClient,
	)

	// テスト用画像の読み込み
	samplesDir := "tools/bedrock-playground/samples"

	l.Info(
		"Reading test images from directory...",
		zap.String("dir", samplesDir),
	)

	files, err := os.ReadDir(samplesDir)
	if err != nil {
		l.Fatal(
			"failed to read samples directory",
			zap.Error(err),
		)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		imagePath := fmt.Sprintf(
			"%s/%s",
			samplesDir,
			file.Name(),
		)

		l.Info(
			"Processing image...",
			zap.String("path", imagePath),
		)

		imgData, err := os.ReadFile(imagePath)
		if err != nil {
			l.Error(
				"failed to read image file",
				zap.String("path", imagePath),
				zap.Error(err),
			)
			continue
		}

		result, err := bedrockService.ProcessImage(
			ctx,
			imgData,
		)
		if err != nil {
			l.Error(
				"ProcessImage failed",
				zap.String("path", imagePath),
				zap.Error(err),
			)
			continue
		}

		fmt.Printf(
			"=== Success: %s ===\n",
			file.Name(),
		)

		for k, v := range result {
			fmt.Printf(
				"  %s: %+v\n",
				k,
				v,
			)
		}

		fmt.Println()
	}
}
