package config

func LoadAPI() (*APIConfig, error) {
	appConfig := loadAppConfig()
	awsConfig := loadAWSConfig()

	authConfig, err := loadBFFAuthConfig()
	if err != nil {
		return nil, err
	}

	slackConfig, err := loadSlackConfig()
	if err != nil {
		return nil, err
	}

	storageConfig, err := loadAPIStorageConfig()
	if err != nil {
		return nil, err
	}

	return &APIConfig{
		App:     appConfig,
		AWS:     awsConfig,
		Auth:    authConfig,
		Slack:   slackConfig,
		Storage: storageConfig,
	}, nil
}

func LoadProcessor() (*ProcessorConfig, error) {
	appConfig := loadAppConfig()
	awsConfig := loadAWSConfig()

	storageConfig, err :=
		loadProcessorStorageConfig()
	if err != nil {
		return nil, err
	}

	bedrockConfig, err := loadBedrockConfig()
	if err != nil {
		return nil, err
	}

	slackTokenConfig, err :=
		loadSlackTokenConfig()
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

func LoadBedrockPlayground() (
	*BedrockPlaygroundConfig,
	error,
) {
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
