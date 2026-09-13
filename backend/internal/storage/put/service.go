package put

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
	PresignPutObject(
		ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.PresignOptions),
	) (*v4.PresignedHTTPRequest, error)
}

type Service struct{ s3Presigner S3Presigner }

func NewService(s3Presigner S3Presigner) *Service {
	return &Service{s3Presigner: s3Presigner}
}

func (s *Service) GeneratePresignURL(
	ctx context.Context,
	cognitoSub string,
	shotNumber string,
) (string, time.Time, error) {
	now := time.Now()
	duration := 15 * time.Minute
	objectKey := storage.BuildObjectKey(cognitoSub, shotNumber)

	request, err := s.s3Presigner.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Key:         aws.String(objectKey),
			ContentType: aws.String("image/jpeg"),
			Metadata: map[string]string{
				"shot-number": shotNumber,
			},
		},
		s3.WithPresignExpires(duration),
	)

	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign request: %w", err)
	}

	return request.URL, now.Add(duration), nil
}
