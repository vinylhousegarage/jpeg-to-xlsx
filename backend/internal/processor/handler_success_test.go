package processor

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler_HandleRequest_Success(t *testing.T) {
	t.Parallel()

	const (
		expectedBucket     = "input-bucket"
		expectedKey        = "test-cognito-sub/SHOT-001.jpg"
		expectedCognitoSub = "test-cognito-sub"
	)

	mock := &mockWorkflow{}
	handler := NewHandler(mock)

	event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: expectedBucket,
					},
					Object: events.S3Object{
						Key: expectedKey,
					},
				},
			},
		},
	}

	ctx := context.Background()

	err := handler.HandleRequest(ctx, event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mock.called {
		t.Fatal("expected workflow to be called")
	}

	if mock.callCount != 1 {
		t.Errorf("expected workflow to be called once, got %d", mock.callCount)
	}

	if mock.ctx != ctx {
		t.Error("Execute() received an unexpected context")
	}

	if mock.inputBucket != expectedBucket {
		t.Errorf(
			"expected bucket %q, got %q",
			expectedBucket,
			mock.inputBucket,
		)
	}

	if mock.inputKey != expectedKey {
		t.Errorf(
			"expected key %q, got %q",
			expectedKey,
			mock.inputKey,
		)
	}

	if mock.cognitoSub != expectedCognitoSub {
		t.Errorf(
			"expected Cognito sub %q, got %q",
			expectedCognitoSub,
			mock.cognitoSub,
		)
	}
}

func TestHandler_HandleRequest_DecodesObjectKey(t *testing.T) {
	t.Parallel()

	const (
		encodedKey         = "test-cognito-sub%2FSHOT-001.jpg"
		expectedKey        = "test-cognito-sub/SHOT-001.jpg"
		expectedCognitoSub = "test-cognito-sub"
	)

	mock := &mockWorkflow{}
	handler := NewHandler(mock)

	event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "input-bucket",
					},
					Object: events.S3Object{
						Key: encodedKey,
					},
				},
			},
		},
	}

	err := handler.HandleRequest(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mock.called {
		t.Fatal("expected workflow to be called")
	}

	if mock.callCount != 1 {
		t.Errorf("expected workflow to be called once, got %d", mock.callCount)
	}

	if mock.inputKey != expectedKey {
		t.Errorf("expected key %q, got %q", expectedKey, mock.inputKey)
	}

	if mock.cognitoSub != expectedCognitoSub {
		t.Errorf("expected Cognito sub %q, got %q", expectedCognitoSub, mock.cognitoSub)
	}
}
