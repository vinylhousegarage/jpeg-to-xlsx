package processor

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler_HandleRequest_NoRecords(t *testing.T) {
	t.Parallel()

	mock := &mockWorkflow{}
	handler := NewHandler(mock)

	event := events.S3Event{}

	err := handler.HandleRequest(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "s3 event contains no records") {
		t.Errorf("unexpected error: %v", err)
	}

	if mock.called {
		t.Error("expected workflow not to be called")
	}
}

func TestHandler_HandleRequest_EmptyBucket(t *testing.T) {
	t.Parallel()

	mock := &mockWorkflow{}
	handler := NewHandler(mock)

	event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "",
					},
					Object: events.S3Object{
						Key: "test-cognito-sub/SHOT-001.jpg",
					},
				},
			},
		},
	}

	err := handler.HandleRequest(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "s3 event bucket name is empty") {
		t.Errorf("unexpected error: %v", err)
	}

	if mock.called {
		t.Error("expected workflow not to be called")
	}
}

func TestHandler_HandleRequest_EmptyObjectKey(t *testing.T) {
	t.Parallel()

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
						Key: "",
					},
				},
			},
		},
	}

	err := handler.HandleRequest(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "s3 event object key is empty") {
		t.Errorf("unexpected error: %v", err)
	}

	if mock.called {
		t.Error("expected workflow not to be called")
	}
}

func TestHandler_HandleRequest_InvalidObjectKeyEncoding(
	t *testing.T,
) {
	t.Parallel()

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
						Key: "%ZZ",
					},
				},
			},
		},
	}

	err := handler.HandleRequest(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to decode s3 object key") {
		t.Errorf("unexpected error: %v", err)
	}

	if mock.called {
		t.Error("expected workflow not to be called")
	}
}

func TestHandler_HandleRequest_InvalidObjectKeyFormat(
	t *testing.T,
) {
	t.Parallel()

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
						Key: "SHOT-001.jpg",
					},
				},
			},
		},
	}

	err := handler.HandleRequest(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"failed to extract Cognito sub from S3 object key",
	) {
		t.Errorf("unexpected error: %v", err)
	}

	if mock.called {
		t.Error("expected workflow not to be called")
	}
}

func TestHandler_HandleRequest_WorkflowError(t *testing.T) {
	t.Parallel()

	workflowErr := errors.New("workflow error")

	mock := &mockWorkflow{
		executeErr: workflowErr,
	}
	handler := NewHandler(mock)

	event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "input-bucket",
					},
					Object: events.S3Object{
						Key: "test-cognito-sub/SHOT-001.jpg",
					},
				},
			},
		},
	}

	err := handler.HandleRequest(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, workflowErr) {
		t.Errorf(
			"HandleRequest() error = %v, want wrapped error %v",
			err,
			workflowErr,
		)
	}

	if !strings.Contains(err.Error(), "failed to execute workflow") {
		t.Errorf("unexpected error: %v", err)
	}

	if !mock.called {
		t.Fatal("expected workflow to be called")
	}

	if mock.callCount != 1 {
		t.Errorf("expected workflow to be called once, got %d", mock.callCount)
	}

	if mock.inputBucket != "input-bucket" {
		t.Errorf("Execute() inputBucket = %q, want %q", mock.inputBucket, "input-bucket")
	}

	if mock.inputKey != "test-cognito-sub/SHOT-001.jpg" {
		t.Errorf("Execute() inputKey = %q, want %q", mock.inputKey, "test-cognito-sub/SHOT-001.jpg")
	}

	if mock.cognitoSub != "test-cognito-sub" {
		t.Errorf("Execute() cognitoSub = %q, want %q", mock.cognitoSub, "test-cognito-sub")
	}
}
