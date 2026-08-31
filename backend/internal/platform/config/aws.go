package config

import "os"

func loadAWSConfig() AWSConfig {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultAWSRegion
	}

	return AWSConfig{
		IsLambda: os.Getenv(
			"AWS_LAMBDA_FUNCTION_NAME",
		) != "",
		Region: region,
	}
}
