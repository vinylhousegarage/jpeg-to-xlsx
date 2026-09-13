package usecase

import (
	"context"
	"io"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/notifier"
)

type mockS3Getter struct {
	getOutput *s3.GetObjectOutput
	getErr    error
}

func (m *mockS3Getter) GetObject(
	ctx context.Context,
	params *s3.GetObjectInput,
	optFns ...func(*s3.Options),
) (*s3.GetObjectOutput, error) {
	return m.getOutput, m.getErr
}

type mockS3Putter struct {
	putErr   error
	putInput *s3.PutObjectInput
}

func (m *mockS3Putter) PutObject(
	_ context.Context,
	params *s3.PutObjectInput,
	_ ...func(*s3.Options),
) (*s3.PutObjectOutput, error) {
	m.putInput = params

	return &s3.PutObjectOutput{}, m.putErr
}

type mockS3Presigner struct {
	presignURL string
	presignErr error
}

func (m *mockS3Presigner) PresignGetObject(
	ctx context.Context,
	params *s3.GetObjectInput,
	optFns ...func(*s3.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {
	if m.presignErr != nil {
		return nil, m.presignErr
	}

	return &v4.PresignedHTTPRequest{
		URL: m.presignURL,
	}, nil
}

type mockBedrockService struct {
	resultMap  map[string]any
	processErr error
}

func (m *mockBedrockService) ProcessImage(
	ctx context.Context,
	rawImage []byte,
) (map[string]any, error) {
	return m.resultMap, m.processErr
}

type mockSlackNotifier struct {
	notifyErr  error
	called     bool
	ctx        context.Context
	cognitoSub string
	message    notifier.Message
}

func (m *mockSlackNotifier) Notify(
	ctx context.Context,
	cognitoSub string,
	message notifier.Message,
) error {
	m.called = true
	m.ctx = ctx
	m.cognitoSub = cognitoSub
	m.message = message

	return m.notifyErr
}

type errorReadCloser struct {
	err error
}

func (r *errorReadCloser) Read(p []byte) (int, error) {
	return 0, r.err
}

func (r *errorReadCloser) Close() error {
	return nil
}

var _ io.ReadCloser = (*errorReadCloser)(nil)
