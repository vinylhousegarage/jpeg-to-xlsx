package config

import (
	"fmt"
	"os"
)

func LoadAPI() (*APIConfig, error) {
	appConfig := loadAppConfig()
	awsConfig := loadAWSConfig()

	storageConfig, err := loadAPIStorageConfig()
	if err != nil {
		return nil, err
	}

	slackConfig, err := loadSlackConfig()
	if err != nil {
		return nil, err
	}

	return &APIConfig{
		App:     appConfig,
		AWS:     awsConfig,
		Storage: storageConfig,
		Slack:   slackConfig,
	}, nil
}

func LoadProcessor() (*ProcessorConfig, error) {
	appConfig := loadAppConfig()
	awsConfig := loadAWSConfig()

	storageConfig, err := loadProcessorStorageConfig()
	if err != nil {
		return nil, err
	}

	bedrockConfig, err := loadBedrockConfig()
	if err != nil {
		return nil, err
	}

	slackTokenConfig, err := loadSlackTokenConfig()
	if err != nil {
		return nil, err
	}

	return &ProcessorConfig{
		App:        appConfig,
		AWS:        awsConfig,
		Storage:    storageConfig,
		Bedrock:    bedrockConfig,
		SlackToken: slackTokenConfig,
	}, nil
}

func LoadBedrockPlayground() (*BedrockPlaygroundConfig, error) {
	appConfig := loadAppConfig()
	awsConfig := loadAWSConfig()

	bedrockConfig, err := loadBedrockConfig()
	if err != nil {
		return nil, err
	}

	return &BedrockPlaygroundConfig{
		App:     appConfig,
		AWS:     awsConfig,
		Bedrock: bedrockConfig,
	}, nil
}

func loadAppConfig() AppConfig {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = appEnvLocal
	}

	return AppConfig{
		Env:          env,
		CookieSecure: env == appEnvProduction || env == appEnvStaging,
	}
}

func loadAWSConfig() AWSConfig {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultAWSRegion
	}

	return AWSConfig{
		IsLambda: os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "",
		Region:   region,
	}
}

func loadAPIStorageConfig() (StorageConfig, error) {
	inputBucketName, err := loadInputBucketName()
	if err != nil {
		return StorageConfig{}, err
	}

	return StorageConfig{
		InputBucketName: inputBucketName,
	}, nil
}

func loadProcessorStorageConfig() (StorageConfig, error) {
	inputBucketName, err := loadInputBucketName()
	if err != nil {
		return StorageConfig{}, err
	}

	outputBucketName := os.Getenv("OUTPUT_BUCKET_NAME")
	if outputBucketName == "" {
		return StorageConfig{}, fmt.Errorf(
			"OUTPUT_BUCKET_NAME is required",
		)
	}

	return StorageConfig{
		InputBucketName:  inputBucketName,
		OutputBucketName: outputBucketName,
	}, nil
}

func loadInputBucketName() (string, error) {
	inputBucketName := os.Getenv("INPUT_BUCKET_NAME")
	if inputBucketName == "" {
		return "", fmt.Errorf(
			"INPUT_BUCKET_NAME is required",
		)
	}

	return inputBucketName, nil
}

func loadBedrockConfig() (BedrockConfig, error) {
	modelID := os.Getenv("BEDROCK_MODEL_ID")
	if modelID == "" {
		return BedrockConfig{}, fmt.Errorf(
			"BEDROCK_MODEL_ID is required",
		)
	}

	promptFileName := os.Getenv("PROMPT_FILE_NAME")
	if promptFileName == "" {
		promptFileName = defaultPromptFileName
	}

	return BedrockConfig{
		ModelID:        modelID,
		PromptFileName: promptFileName,
	}, nil
}

func loadSlackConfig() (SlackConfig, error) {
	clientID := os.Getenv("SLACK_CLIENT_ID")
	if clientID == "" {
		return SlackConfig{}, fmt.Errorf(
			"SLACK_CLIENT_ID is required",
		)
	}

	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	clientSecretARN := os.Getenv("SLACK_CLIENT_SECRET_ARN")

	if clientSecret == "" && clientSecretARN == "" {
		return SlackConfig{}, fmt.Errorf(
			"SLACK_CLIENT_SECRET or SLACK_CLIENT_SECRET_ARN is required",
		)
	}

	redirectURI := os.Getenv("SLACK_REDIRECT_URI")
	if redirectURI == "" {
		return SlackConfig{}, fmt.Errorf(
			"SLACK_REDIRECT_URI is required",
		)
	}

	tokenTableName, err := loadSlackTokenTableName()
	if err != nil {
		return SlackConfig{}, err
	}

	return SlackConfig{
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		ClientSecretARN: clientSecretARN,
		RedirectURI:     redirectURI,
		TokenTableName:  tokenTableName,
	}, nil
}

func loadSlackTokenConfig() (SlackTokenConfig, error) {
	tokenTableName, err := loadSlackTokenTableName()
	if err != nil {
		return SlackTokenConfig{}, err
	}

	return SlackTokenConfig{
		TokenTableName: tokenTableName,
	}, nil
}

func loadSlackTokenTableName() (string, error) {
	tokenTableName := os.Getenv("SLACK_TOKEN_TABLE_NAME")
	if tokenTableName == "" {
		return "", fmt.Errorf(
			"SLACK_TOKEN_TABLE_NAME is required",
		)
	}

	return tokenTableName, nil
}
