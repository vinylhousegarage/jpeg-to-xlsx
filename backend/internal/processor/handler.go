package processor

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/storage"
)

type Workflow interface {
	Execute(
		ctx context.Context,
		inputBucket string,
		inputKey string,
		cognitoSub string,
	) error
}

type Handler struct {
	workflow Workflow
}

func NewHandler(workflow Workflow) *Handler {
	return &Handler{
		workflow: workflow,
	}
}

func (h *Handler) HandleRequest(
	ctx context.Context,
	event events.S3Event,
) error {
	if len(event.Records) == 0 {
		return errors.New("s3 event contains no records")
	}

	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		encodedKey := record.S3.Object.Key

		if bucket == "" {
			return errors.New("s3 event bucket name is empty")
		}

		if encodedKey == "" {
			return errors.New("s3 event object key is empty")
		}

		decodedKey, err := url.QueryUnescape(encodedKey)
		if err != nil {
			return fmt.Errorf("failed to decode s3 object key: %w", err)
		}

		cognitoSub, err := storage.ExtractCognitoSub(decodedKey)
		if err != nil {
			return fmt.Errorf(
				"failed to extract Cognito sub from S3 object key %q: %w",
				decodedKey,
				err,
			)
		}

		if err := h.workflow.Execute(
			ctx,
			bucket,
			decodedKey,
			cognitoSub,
		); err != nil {
			return fmt.Errorf(
				"failed to execute workflow for %s/%s: %w",
				bucket,
				decodedKey,
				err,
			)
		}
	}

	return nil
}
