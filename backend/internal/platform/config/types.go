package config

import "time"

type APIConfig struct {
	App     AppConfig
	AWS     AWSConfig
	Auth    BFFAuthConfig
	Slack   SlackConfig
	Storage StorageConfig
}

type ProcessorConfig struct {
	App        AppConfig
	AWS        AWSConfig
	Storage    StorageConfig
	Bedrock    BedrockConfig
	SlackToken SlackTokenConfig
}

type AppConfig struct {
	Env          string
	CookieSecure bool
}

type AWSConfig struct {
	IsLambda bool
	Region   string
}

type StorageConfig struct {
	InputBucketName  string
	OutputBucketName string
}

type BedrockConfig struct {
	ModelID        string
	PromptFileName string
}

type BedrockPlaygroundConfig struct {
	App     AppConfig
	AWS     AWSConfig
	Bedrock BedrockConfig
}

type SlackConfig struct {
	ClientID        string
	ClientSecret    string
	ClientSecretARN string
	RedirectURI     string
	TokenTableName  string
}

type SlackTokenConfig struct {
	TokenTableName string
}

type BFFAuthConfig struct {
	CognitoClientID        string
	CognitoClientSecretARN string

	CognitoIssuer                string
	CognitoAuthorizationEndpoint string
	CognitoTokenEndpoint         string
	CognitoLogoutEndpoint        string
	CognitoRedirectURI           string

	PostLoginRedirectURL  string
	PostLogoutRedirectURL string

	OAuthStateTableName string
	OAuthStateTTL       time.Duration

	SessionTableName string
	SessionLifetime  time.Duration
}
