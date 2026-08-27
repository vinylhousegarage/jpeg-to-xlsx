package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/notifier"
)

// インターフェースを定義
type S3Getter interface {
	GetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.GetObjectOutput, error)
}

type S3Putter interface {
	PutObject(
		ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.PutObjectOutput, error)
}

type S3Presigner interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions),
	) (*v4.PresignedHTTPRequest, error)
}

type BedrockService interface {
	ProcessImage(
		ctx context.Context,
		rawImage []byte,
	) (map[string]any, error)
}

type SlackNotifier interface {
	Notify(
		ctx context.Context,
		message notifier.Message,
	) error
}

// 構造体を定義
type Workflow struct {
	s3Getter       S3Getter
	s3Putter       S3Putter
	s3Presigner    S3Presigner
	bedrockService BedrockService
	slackNotifier  SlackNotifier
	outputBucket   string
	logger         *zap.Logger
}

// 構造体を初期化
func NewWorkflow(
	s3Getter S3Getter,
	s3Putter S3Putter,
	s3Presigner S3Presigner,
	bedrockService BedrockService,
	slackNotifier SlackNotifier,
	outputBucket string,
	logger *zap.Logger,
) *Workflow {
	return &Workflow{
		s3Getter:       s3Getter,
		s3Putter:       s3Putter,
		s3Presigner:    s3Presigner,
		bedrockService: bedrockService,
		slackNotifier:  slackNotifier,
		outputBucket:   outputBucket,
		logger:         logger,
	}
}

// Execute ワークフローの実行関数
func (w *Workflow) Execute(
	ctx context.Context,
	inputBucket string,
	inputKey string,
) error {
	w.logger.Info(
		"starting workflow",
		zap.String("bucket", inputBucket),
		zap.String("key", inputKey),
	)

	// S3オブジェクトキーから撮影番号を取得
	shotNumber := strings.TrimSuffix(
		filepath.Base(inputKey),
		filepath.Ext(inputKey),
	)

	// アップロードされた画像をS3から取得
	objOutput, err := w.s3Getter.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(inputBucket),
			Key:    aws.String(inputKey),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get object from s3: %w", err)
	}
	defer func() {
		if err := objOutput.Body.Close(); err != nil {
			w.logger.Warn(
				"failed to close s3 object body",
				zap.Error(err),
			)
		}
	}()

	// io.ReadCloserから画像データを読み込み
	imgData, err := io.ReadAll(objOutput.Body)
	if err != nil {
		return fmt.Errorf("failed to read image body: %w", err)
	}

	// Bedrockで画像を解析してJSONへ変換
	jsonMap, err := w.bedrockService.ProcessImage(ctx, imgData)
	if err != nil {
		return fmt.Errorf("failed to process image with bedrock: %w", err)
	}

	// mapをJSONバイト列へ変換
	jsonBytes, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	// JSON保存先のオブジェクトキーを生成
	outputKey := strings.TrimSuffix(
		inputKey,
		filepath.Ext(inputKey),
	) + ".json"

	// JSONをS3へ保存
	_, err = w.s3Putter.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(w.outputBucket),
			Key:    aws.String(outputKey),
			Body:   bytes.NewReader(jsonBytes),
			ContentType: aws.String(
				"application/json; charset=utf-8",
			),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to put json to output s3: %w", err)
	}

	// JSONダウンロード用の署名付きURLを生成
	presignReq, err := w.s3Presigner.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(w.outputBucket),
			Key:    aws.String(outputKey),
		},
		s3.WithPresignExpires(24*time.Hour),
	)
	if err != nil {
		return fmt.Errorf("failed to generate presigned url: %w", err)
	}

	// Slack OAuthで連携済みのSlackへ通知
	if err := w.slackNotifier.Notify(
		ctx,
		notifier.Message{
			ShotNumber:  shotNumber,
			DownloadURL: presignReq.URL,
		},
	); err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}

	w.logger.Info(
		"workflow completed successfully",
		zap.String("shot_number", shotNumber),
		zap.String("output_key", outputKey),
	)

	return nil
}
