package processor

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

type Workflow interface {
	Execute(
		ctx context.Context,
		inputBucket string,
		inputKey string,
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
		return fmt.Errorf("s3 event contains no records")
	}

	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		if bucket == "" {
			return fmt.Errorf("s3 event bucket name is empty")
		}

		if key == "" {
			return fmt.Errorf("s3 event object key is empty")
		}

		decodedKey, err := url.QueryUnescape(key)
		if err != nil {
			return fmt.Errorf("failed to decode s3 object key: %w", err)
		}

		if err := h.workflow.Execute(
			ctx,
			bucket,
			decodedKey,
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
