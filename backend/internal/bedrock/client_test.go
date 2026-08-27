package bedrock

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type mockBedrockRuntimeClient struct {
	output *bedrockruntime.ConverseOutput
	err    error
	input  *bedrockruntime.ConverseInput
}

func (m *mockBedrockRuntimeClient) Converse(
	ctx context.Context,
	params *bedrockruntime.ConverseInput,
	optFns ...func(*bedrockruntime.Options),
) (*bedrockruntime.ConverseOutput, error) {
	m.input = params

	return m.output, m.err
}

func TestBedrockClient_Analyze_Success(t *testing.T) {
	t.Parallel()

	mock := &mockBedrockRuntimeClient{
		output: &bedrockruntime.ConverseOutput{
			Output: &types.ConverseOutputMemberMessage{
				Value: types.Message{
					Content: []types.ContentBlock{
						&types.ContentBlockMemberText{
							Value: `{"key":"value"}`,
						},
					},
				},
			},
			Usage: &types.TokenUsage{
				InputTokens:  aws.Int32(100),
				OutputTokens: aws.Int32(50),
			},
		},
	}

	client := NewClient(
		mock,
		"test-model-id",
		"test prompt",
		zap.NewNop(),
	)

	result, err := client.Analyze(
		context.Background(),
		[]byte("fake-image"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != `{"key":"value"}` {
		t.Errorf(
			"expected result %q, got %q",
			`{"key":"value"}`,
			result,
		)
	}

	if mock.input == nil {
		t.Fatal("expected Converse to be called")
	}

	if aws.ToString(mock.input.ModelId) != "test-model-id" {
		t.Errorf(
			"expected model ID %q, got %q",
			"test-model-id",
			aws.ToString(mock.input.ModelId),
		)
	}

	if len(mock.input.Messages) != 1 {
		t.Fatalf(
			"expected one message, got %d",
			len(mock.input.Messages),
		)
	}

	content := mock.input.Messages[0].Content
	if len(content) != 2 {
		t.Fatalf(
			"expected two content blocks, got %d",
			len(content),
		)
	}

	textBlock, ok := content[1].(*types.ContentBlockMemberText)
	if !ok {
		t.Fatal("expected second content block to be text")
	}

	if textBlock.Value != "test prompt" {
		t.Errorf(
			"expected prompt %q, got %q",
			"test prompt",
			textBlock.Value,
		)
	}
}

func TestBedrockClient_Analyze_ConverseError(t *testing.T) {
	t.Parallel()

	mock := &mockBedrockRuntimeClient{
		err: errors.New("bedrock error"),
	}

	client := NewClient(
		mock,
		"test-model-id",
		"test prompt",
		zap.NewNop(),
	)

	_, err := client.Analyze(
		context.Background(),
		[]byte("fake-image"),
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBedrockClient_Analyze_NoContent(t *testing.T) {
	t.Parallel()

	mock := &mockBedrockRuntimeClient{
		output: &bedrockruntime.ConverseOutput{
			Output: &types.ConverseOutputMemberMessage{
				Value: types.Message{
					Content: []types.ContentBlock{},
				},
			},
		},
	}

	client := NewClient(
		mock,
		"test-model-id",
		"test prompt",
		zap.NewNop(),
	)

	_, err := client.Analyze(
		context.Background(),
		[]byte("fake-image"),
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBedrockClient_Analyze_ResponseContentIsNotText(t *testing.T) {
	t.Parallel()

	mock := &mockBedrockRuntimeClient{
		output: &bedrockruntime.ConverseOutput{
			Output: &types.ConverseOutputMemberMessage{
				Value: types.Message{
					Content: []types.ContentBlock{
						&types.ContentBlockMemberImage{
							Value: types.ImageBlock{
								Format: types.ImageFormatJpeg,
								Source: &types.ImageSourceMemberBytes{
									Value: []byte("fake-image"),
								},
							},
						},
					},
				},
			},
		},
	}

	client := NewClient(
		mock,
		"test-model-id",
		"test prompt",
		zap.NewNop(),
	)

	_, err := client.Analyze(
		context.Background(),
		[]byte("fake-image"),
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
