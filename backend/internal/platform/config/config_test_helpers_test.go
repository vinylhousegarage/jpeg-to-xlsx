package config

import "testing"

const (
	testInputBucketName  = "my-test-bucket"
	testOutputBucketName = "my-output-bucket"

	testBedrockModelID = "jp.anthropic.claude-sonnet-4-6"

	testSlackClientID       = "test-client-id"
	testSlackClientSecret   = "test-client-secret"
	testSlackRedirectURI    = "https://example.com/api/oauth/slack/callback"
	testSlackTokenTableName = "slack-tokens"

	testCognitoClientID = "test-cognito-client-id"

	testCognitoClientSecretARN = "arn:aws:secretsmanager:" +
		"ap-northeast-1:123456789012:" +
		"secret:test-cognito-client-secret"

	testCognitoIssuer = "https://cognito-idp.ap-northeast-1." +
		"amazonaws.com/ap-northeast-1_test"

	testCognitoAuthorizationEndpoint = "https://test.auth.ap-northeast-1." +
		"amazoncognito.com/oauth2/authorize"

	testCognitoTokenEndpoint = "https://test.auth.ap-northeast-1." +
		"amazoncognito.com/oauth2/token"

	testCognitoLogoutEndpoint = "https://test.auth.ap-northeast-1." +
		"amazoncognito.com/logout"

	testCognitoRedirectURI = "https://example.com/api/auth/callback"

	testPostLoginRedirectURL  = "https://example.com/"
	testPostLogoutRedirectURL = "https://example.com/"

	testOAuthStateTableName  = "test-cognito-oauth-states"
	testAuthSessionTableName = "test-auth-sessions"
)

func setValidAPIEnvironment(t *testing.T) {
	t.Helper()

	t.Setenv("APP_ENV", appEnvLocal)
	t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "")
	t.Setenv("AWS_REGION", defaultAWSRegion)
	t.Setenv("INPUT_BUCKET_NAME", testInputBucketName)

	setValidBFFAuthEnvironment(t)
	setValidSlackEnvironment(t)
}

func setValidBFFAuthEnvironment(t *testing.T) {
	t.Helper()

	t.Setenv("COGNITO_CLIENT_ID", testCognitoClientID)
	t.Setenv("COGNITO_CLIENT_SECRET_ARN", testCognitoClientSecretARN)
	t.Setenv("COGNITO_ISSUER", testCognitoIssuer)
	t.Setenv("COGNITO_AUTHORIZATION_ENDPOINT", testCognitoAuthorizationEndpoint)
	t.Setenv("COGNITO_TOKEN_ENDPOINT", testCognitoTokenEndpoint)
	t.Setenv("COGNITO_LOGOUT_ENDPOINT", testCognitoLogoutEndpoint)
	t.Setenv("COGNITO_REDIRECT_URI", testCognitoRedirectURI)
	t.Setenv("AUTH_REDIRECT_URL", testPostLoginRedirectURL)
	t.Setenv("AUTH_LOGOUT_REDIRECT_URL", testPostLogoutRedirectURL)
	t.Setenv("COGNITO_OAUTH_STATE_TABLE_NAME", testOAuthStateTableName)
	t.Setenv("AUTH_SESSION_TABLE_NAME", testAuthSessionTableName)
}

func setValidSlackEnvironment(t *testing.T) {
	t.Helper()

	t.Setenv("SLACK_CLIENT_ID", testSlackClientID)
	t.Setenv("SLACK_CLIENT_SECRET", testSlackClientSecret)
	t.Setenv("SLACK_CLIENT_SECRET_ARN", "")
	t.Setenv("SLACK_REDIRECT_URI", testSlackRedirectURI)
	t.Setenv("SLACK_TOKEN_TABLE_NAME", testSlackTokenTableName)
}
