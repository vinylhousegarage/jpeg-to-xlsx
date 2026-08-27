package processor

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler_HandleRequest_Success(t *testing.T) {
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

	err := handler.HandleRequest(
		context.Background(),
		event,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mock.called {
		t.Fatal("expected workflow to be called")
	}

	if mock.callCount != 1 {
		t.Errorf(
			"expected workflow to be called once, got %d",
			mock.callCount,
		)
	}

	if mock.bucket != "input-bucket" {
		t.Errorf(
			"expected bucket %q, got %q",
			"input-bucket",
			mock.bucket,
		)
	}

	if mock.key != "SHOT-001.jpg" {
		t.Errorf(
			"expected key %q, got %q",
			"SHOT-001.jpg",
			mock.key,
		)
	}
}

func TestHandler_HandleRequest_DecodesObjectKey(t *testing.T) {
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
						Key: "images%2FSHOT-001.jpg",
					},
				},
			},
		},
	}

	err := handler.HandleRequest(
		context.Background(),
		event,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if mock.key != "images/SHOT-001.jpg" {
		t.Errorf(
			"expected key %q, got %q",
			"images/SHOT-001.jpg",
			mock.key,
		)
	}
}
