package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

func TestWorkflow_Execute_SlackError(
	t *testing.T,
) {
	t.Parallel()

	const (
		cognitoSub = "test-cognito-sub"
		inputKey   = "test-cognito-sub/SHOT-001.jpg"
	)

	slackErr := errors.New("slack error")

	slackNotifier := &mockSlackNotifier{
		notifyErr: slackErr,
	}

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader([]byte("fake-image-bytes"))),
			},
		},
		&mockS3Putter{},
		&mockS3Presigner{
			presignURL: "https://example.com/download",
		},
		&mockBedrockService{
			resultMap: map[string]any{
				"key": "value",
			},
		},
		slackNotifier,
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		inputKey,
		cognitoSub,
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want an error")
	}

	if !errors.Is(err, slackErr) {
		t.Errorf("Execute() error = %v, want wrapped error %v", err, slackErr)
	}

	const wantError = "failed to send slack notification: slack error"

	if err.Error() != wantError {
		t.Errorf("Execute() error = %q, want %q", err.Error(), wantError)
	}

	if !slackNotifier.called {
		t.Fatal("Notify() was not called")
	}

	if slackNotifier.cognitoSub != cognitoSub {
		t.Errorf(
			"Notify() cognitoSub = %q, want %q",
			slackNotifier.cognitoSub,
			cognitoSub,
		)
	}

	if slackNotifier.message.ShotNumber != "SHOT-001" {
		t.Errorf(
			"Notify() ShotNumber = %q, want %q",
			slackNotifier.message.ShotNumber,
			"SHOT-001",
		)
	}

	if slackNotifier.message.DownloadURL != "https://example.com/download" {
		t.Errorf(
			"Notify() DownloadURL = %q, want %q",
			slackNotifier.message.DownloadURL,
			"https://example.com/download",
		)
	}
}
