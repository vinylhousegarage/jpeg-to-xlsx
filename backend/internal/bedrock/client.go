package bedrock

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type BedrockRuntimeClient interface {
	Converse(
		ctx context.Context,
		params *bedrockruntime.ConverseInput,
		optFns ...func(*bedrockruntime.Options),
	) (*bedrockruntime.ConverseOutput, error)
}

type BedrockClient struct {
	sdkClient BedrockRuntimeClient
	modelID   string
	prompt    string
	logger    *zap.Logger
}

func NewClient(
	sdkClient BedrockRuntimeClient,
	modelID string,
	prompt string,
	logger *zap.Logger,
) *BedrockClient {
	return &BedrockClient{
		sdkClient: sdkClient,
		modelID:   modelID,
		prompt:    prompt,
		logger:    logger,
	}
}

func (c *BedrockClient) Analyze(
	ctx context.Context,
	imgData []byte,
) (string, error) {
	input := &bedrockruntime.ConverseInput{
		ModelId: &c.modelID,
		Messages: []types.Message{
			{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberImage{
						Value: types.ImageBlock{
							Format: types.ImageFormatJpeg,
							Source: &types.ImageSourceMemberBytes{
								Value: imgData,
							},
						},
					},
					&types.ContentBlockMemberText{
						Value: c.prompt,
					},
				},
			},
		},
		InferenceConfig: &types.InferenceConfiguration{
			MaxTokens: aws.Int32(2000),
		},
	}

	output, err := c.sdkClient.Converse(ctx, input)
	if err != nil {
		return "", fmt.Errorf(
			"failed to invoke bedrock converse: %w",
			err,
		)
	}

	message, ok := output.Output.(*types.ConverseOutputMemberMessage)
	if !ok || len(message.Value.Content) == 0 {
		return "", fmt.Errorf("no content in bedrock response")
	}

	textBlock, ok := message.Value.Content[0].(*types.ContentBlockMemberText)
	if !ok {
		return "", fmt.Errorf("response content is not text")
	}

	if output.Usage != nil {
		c.logger.Info(
			"bedrock inference success",
			zap.Int32(
				"input_tokens",
				aws.ToInt32(output.Usage.InputTokens),
			),
			zap.Int32(
				"output_tokens",
				aws.ToInt32(output.Usage.OutputTokens),
			),
		)
	}

	return textBlock.Value, nil
}
