package config

import "time"

const (
	appEnvLocal      = "local"
	appEnvStaging    = "staging"
	appEnvProduction = "production"

	defaultAWSRegion = "ap-northeast-1"

	defaultPromptFileName = "extractor.txt"

	defaultOAuthStateTTL = 10 * time.Minute

	defaultAuthSessionLifetime = 8 * time.Hour
)
