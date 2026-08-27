package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3用構造体を定義
type S3Client struct {
	bucketName string
	client     *s3.Client
}

// 署名付きURL用構造体を定義
type S3PresignClient struct {
	bucketName    string
	presignClient *s3.PresignClient
}

// S3用構造体を初期化
func NewS3Client(bucketName string, client *s3.Client) *S3Client {
	return &S3Client{
		bucketName: bucketName,
		client:     client,
	}
}

// 署名付きURL用構造体を初期化
func NewS3PresignClient(bucketName string, client *s3.Client) *S3PresignClient {
	return &S3PresignClient{
		bucketName:    bucketName,
		presignClient: s3.NewPresignClient(client),
	}
}

// Getメソッド
func (c *S3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if params.Bucket == nil || *params.Bucket == "" {
		params.Bucket = aws.String(c.bucketName)
	}
	return c.client.GetObject(ctx, params, optFns...)
}

// Putメソッド
func (c *S3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if params.Bucket == nil || *params.Bucket == "" {
		params.Bucket = aws.String(c.bucketName)
	}
	return c.client.PutObject(ctx, params, optFns...)
}

// PresignGetメソッド
func (c *S3PresignClient) PresignGetObject(
	ctx context.Context,
	params *s3.GetObjectInput,
	optFns ...func(*s3.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {
	if params.Bucket == nil || *params.Bucket == "" {
		params.Bucket = aws.String(c.bucketName)
	}
	return c.presignClient.PresignGetObject(ctx, params, optFns...)
}

// PresignPutメソッド
func (c *S3PresignClient) PresignPutObject(
	ctx context.Context,
	params *s3.PutObjectInput,
	optFns ...func(*s3.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {
	if params.Bucket == nil || *params.Bucket == "" {
		params.Bucket = aws.String(c.bucketName)
	}
	return c.presignClient.PresignPutObject(ctx, params, optFns...)
}
