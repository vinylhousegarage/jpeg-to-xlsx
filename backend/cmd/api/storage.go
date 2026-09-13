package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"

	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/platform/config"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/storage"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/storage/put"
)

func buildStorageHandler(
	storageConfig config.StorageConfig,
	baseS3Client *s3.Client,
	sessionResolver *authsession.Resolver,
	appLogger *zap.Logger,
) (
	http.Handler,
	error,
) {
	if baseS3Client == nil {
		return nil, errors.New("build storage handler: S3 client is nil")
	}

	if sessionResolver == nil {
		return nil, errors.New("build storage handler: session resolver is nil")
	}

	if appLogger == nil {
		return nil, errors.New("build storage handler: logger is nil")
	}

	if strings.TrimSpace(storageConfig.InputBucketName) == "" {
		return nil, errors.New("build storage handler: input bucket name is empty")
	}

	presignClient := storage.NewS3PresignClient(
		storageConfig.InputBucketName,
		baseS3Client,
	)

	putService := put.NewService(presignClient)
	putHandler := put.NewHandler(putService, sessionResolver, appLogger)

	return putHandler, nil
}
