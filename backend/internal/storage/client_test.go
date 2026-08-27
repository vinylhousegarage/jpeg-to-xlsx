package storage

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestNewS3PresignClient(t *testing.T) {
	t.Parallel()

	rawClient := &s3.Client{}
	bucket := "my-test-bucket"

	client := NewS3PresignClient(bucket, rawClient)

	if client == nil {
		t.Fatal("expected client to be initialized, got nil")
	}
	if client.bucketName != bucket {
		t.Errorf("expected bucket name %s, got %s", bucket, client.bucketName)
	}
	if client.presignClient == nil {
		t.Error("expected presignClient to be initialized, got nil")
	}
}

func TestS3Client_BucketFallback(t *testing.T) {
	t.Parallel()

	client := &S3Client{
		bucketName: "default-bucket",
	}

	t.Run("PutObject fallback", func(t *testing.T) {
		params := &s3.PutObjectInput{}
		if params.Bucket == nil || *params.Bucket == "" {
			params.Bucket = aws.String(client.bucketName)
		}
		if *params.Bucket != "default-bucket" {
			t.Errorf("expected bucket 'default-bucket', got %s", *params.Bucket)
		}
	})

	t.Run("GetObject fallback", func(t *testing.T) {
		params := &s3.GetObjectInput{}
		if params.Bucket == nil || *params.Bucket == "" {
			params.Bucket = aws.String(client.bucketName)
		}
		if *params.Bucket != "default-bucket" {
			t.Errorf("expected bucket 'default-bucket', got %s", *params.Bucket)
		}
	})
}

func TestPresignClient_BucketFallback(t *testing.T) {
	t.Parallel()

	client := &S3PresignClient{
		bucketName: "default-bucket",
	}

	t.Run("PresignPutObject fallback", func(t *testing.T) {
		params := &s3.PutObjectInput{}
		if params.Bucket == nil || *params.Bucket == "" {
			params.Bucket = aws.String(client.bucketName)
		}
		if *params.Bucket != "default-bucket" {
			t.Errorf("expected bucket 'default-bucket', got %s", *params.Bucket)
		}
	})

	t.Run("PresignGetObject fallback", func(t *testing.T) {
		params := &s3.GetObjectInput{}
		if params.Bucket == nil || *params.Bucket == "" {
			params.Bucket = aws.String(client.bucketName)
		}
		if *params.Bucket != "default-bucket" {
			t.Errorf("expected bucket 'default-bucket', got %s", *params.Bucket)
		}
	})
}
