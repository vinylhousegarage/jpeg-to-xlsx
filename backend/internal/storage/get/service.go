package get

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/storage"
)

type S3Presigner interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions),
	) (
		*v4.PresignedHTTPRequest,
		error,
	)
}

type Service struct {
	s3Presigner S3Presigner
}

func NewService(
	s3Presigner S3Presigner,
) *Service {
	return &Service{
		s3Presigner: s3Presigner,
	}
}

func (s *Service) GeneratePresignURL(
	ctx context.Context,
	cognitoSub string,
	shotNumber string,
) (
	string,
	time.Time,
	error,
) {
	now := time.Now()
	duration := 15 * time.Minute
	objectKey := storage.BuildObjectKey(cognitoSub, shotNumber)

	request, err := s.s3Presigner.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Key: aws.String(objectKey),
		},
		s3.WithPresignExpires(duration),
	)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign get request: %w", err)
	}

	return request.URL, now.Add(duration), nil
}
