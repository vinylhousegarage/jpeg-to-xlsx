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

func TestWorkflow_Execute_BedrockError(t *testing.T) {
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
		&mockS3Presigner{},
		&mockBedrockService{
			processErr: errors.New("bedrock error"),
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
		"failed to process image with bedrock",
	) {
		t.Errorf(
			"Execute() error = %q, want Bedrock processing error",
			err.Error(),
		)
	}
}

func TestWorkflow_Execute_MarshalError(t *testing.T) {
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
		&mockS3Presigner{},
		&mockBedrockService{
			resultMap: map[string]any{
				"invalid": make(chan int),
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
		"failed to marshal json",
	) {
		t.Errorf(
			"Execute() error = %q, want JSON marshal error",
			err.Error(),
		)
	}
}
