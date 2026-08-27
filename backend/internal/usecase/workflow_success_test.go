package usecase

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

func TestWorkflow_Execute_Success(t *testing.T) {
	t.Parallel()

	const (
		inputBucket         = "input-bucket"
		outputBucket        = "output-bucket"
		inputKey            = "SHOT-001.jpg"
		expectedOutputKey   = "SHOT-001.json"
		expectedShotNumber  = "SHOT-001"
		expectedDownloadURL = "https://example.com/download"
		expectedContentType = "application/json; charset=utf-8"
	)

	ctx := context.Background()

	s3Putter := &mockS3Putter{}
	slackNotifier := &mockSlackNotifier{}

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: io.NopCloser(
					bytes.NewReader(
						[]byte("fake-image-bytes"),
					),
				),
			},
		},
		s3Putter,
		&mockS3Presigner{
			presignURL: expectedDownloadURL,
		},
		&mockBedrockService{
			resultMap: map[string]any{
				"key": "value",
			},
		},
		slackNotifier,
		outputBucket,
		zap.NewNop(),
	)

	err := workflow.Execute(
		ctx,
		inputBucket,
		inputKey,
	)
	if err != nil {
		t.Fatalf(
			"Execute() error = %v",
			err,
		)
	}

	if s3Putter.putInput == nil {
		t.Fatal("PutObject() input is nil")
	}

	if got := aws.ToString(
		s3Putter.putInput.Bucket,
	); got != outputBucket {
		t.Errorf(
			"PutObject() Bucket = %q, want %q",
			got,
			outputBucket,
		)
	}

	if got := aws.ToString(
		s3Putter.putInput.Key,
	); got != expectedOutputKey {
		t.Errorf(
			"PutObject() Key = %q, want %q",
			got,
			expectedOutputKey,
		)
	}

	if got := aws.ToString(
		s3Putter.putInput.ContentType,
	); got != expectedContentType {
		t.Errorf(
			"PutObject() ContentType = %q, want %q",
			got,
			expectedContentType,
		)
	}

	if !slackNotifier.called {
		t.Fatal("Notify() was not called")
	}

	if slackNotifier.ctx != ctx {
		t.Error(
			"Notify() received an unexpected context",
		)
	}

	if slackNotifier.message.ShotNumber !=
		expectedShotNumber {
		t.Errorf(
			"Notify() ShotNumber = %q, want %q",
			slackNotifier.message.ShotNumber,
			expectedShotNumber,
		)
	}

	if slackNotifier.message.DownloadURL !=
		expectedDownloadURL {
		t.Errorf(
			"Notify() DownloadURL = %q, want %q",
			slackNotifier.message.DownloadURL,
			expectedDownloadURL,
		)
	}
}
