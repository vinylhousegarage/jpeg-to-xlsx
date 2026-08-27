package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

func TestWorkflow_Execute_GetObjectError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getErr: errors.New("s3 get error"),
		},
		&mockS3Putter{},
		&mockS3Presigner{},
		&mockBedrockService{},
		&mockSlackNotifier{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want an error")
	}

	if !strings.Contains(
		err.Error(),
		"failed to get object from s3",
	) {
		t.Errorf(
			"Execute() error = %q, want S3 get error",
			err.Error(),
		)
	}
}

func TestWorkflow_Execute_ReadImageError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: &errorReadCloser{
					err: errors.New("read error"),
				},
			},
		},
		&mockS3Putter{},
		&mockS3Presigner{},
		&mockBedrockService{},
		&mockSlackNotifier{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want an error")
	}

	if !strings.Contains(
		err.Error(),
		"failed to read image body",
	) {
		t.Errorf(
			"Execute() error = %q, want image read error",
			err.Error(),
		)
	}
}

func TestWorkflow_Execute_PutObjectError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: io.NopCloser(
					bytes.NewReader([]byte("fake-image-bytes")),
				),
			},
		},
		&mockS3Putter{
			putErr: errors.New("s3 put error"),
		},
		&mockS3Presigner{},
		&mockBedrockService{
			resultMap: map[string]any{
				"key": "value",
			},
		},
		&mockSlackNotifier{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want an error")
	}

	if !strings.Contains(
		err.Error(),
		"failed to put xlsx to output s3",
	) {
		t.Errorf(
			"Execute() error = %q, want S3 put error",
			err.Error(),
		)
	}
}

func TestWorkflow_Execute_PresignError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: io.NopCloser(
					bytes.NewReader([]byte("fake-image-bytes")),
				),
			},
		},
		&mockS3Putter{},
		&mockS3Presigner{
			presignErr: errors.New("presign error"),
		},
		&mockBedrockService{
			resultMap: map[string]any{
				"key": "value",
			},
		},
		&mockSlackNotifier{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want an error")
	}

	if !strings.Contains(
		err.Error(),
		"failed to generate presigned url",
	) {
		t.Errorf(
			"Execute() error = %q, want presign error",
			err.Error(),
		)
	}
}
