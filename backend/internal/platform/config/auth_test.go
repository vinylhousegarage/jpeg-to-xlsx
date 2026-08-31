package config

import (
	"strings"
	"testing"
)

func TestLoadBFFAuthConfig(
	t *testing.T,
) {
	setValidBFFAuthEnvironment(t)

	config, err := loadBFFAuthConfig()
	if err != nil {
		t.Fatalf(
			"loadBFFAuthConfig() error = %v",
			err,
		)
	}

	if config.CognitoClientID !=
		testCognitoClientID {
		t.Errorf(
			"CognitoClientID = %q, want %q",
			config.CognitoClientID,
			testCognitoClientID,
		)
	}

	if config.CognitoClientSecretARN !=
		testCognitoClientSecretARN {
		t.Errorf(
			"CognitoClientSecretARN = %q, want %q",
			config.CognitoClientSecretARN,
			testCognitoClientSecretARN,
		)
	}

	if config.CognitoIssuer !=
		testCognitoIssuer {
		t.Errorf(
			"CognitoIssuer = %q, want %q",
			config.CognitoIssuer,
			testCognitoIssuer,
		)
	}

	if config.CognitoAuthorizationEndpoint !=
		testCognitoAuthorizationEndpoint {
		t.Errorf(
			"CognitoAuthorizationEndpoint = %q, want %q",
			config.CognitoAuthorizationEndpoint,
			testCognitoAuthorizationEndpoint,
		)
	}

	if config.CognitoTokenEndpoint !=
		testCognitoTokenEndpoint {
		t.Errorf(
			"CognitoTokenEndpoint = %q, want %q",
			config.CognitoTokenEndpoint,
			testCognitoTokenEndpoint,
		)
	}

	if config.CognitoRedirectURI !=
		testCognitoRedirectURI {
		t.Errorf(
			"CognitoRedirectURI = %q, want %q",
			config.CognitoRedirectURI,
			testCognitoRedirectURI,
		)
	}

	if config.PostLoginRedirectURL !=
		testPostLoginRedirectURL {
		t.Errorf(
			"PostLoginRedirectURL = %q, want %q",
			config.PostLoginRedirectURL,
			testPostLoginRedirectURL,
		)
	}

	if config.OAuthStateTableName !=
		testOAuthStateTableName {
		t.Errorf(
			"OAuthStateTableName = %q, want %q",
			config.OAuthStateTableName,
			testOAuthStateTableName,
		)
	}

	if config.OAuthStateTTL !=
		defaultOAuthStateTTL {
		t.Errorf(
			"OAuthStateTTL = %v, want %v",
			config.OAuthStateTTL,
			defaultOAuthStateTTL,
		)
	}

	if config.SessionTableName !=
		testAuthSessionTableName {
		t.Errorf(
			"SessionTableName = %q, want %q",
			config.SessionTableName,
			testAuthSessionTableName,
		)
	}

	if config.SessionLifetime !=
		defaultAuthSessionLifetime {
		t.Errorf(
			"SessionLifetime = %v, want %v",
			config.SessionLifetime,
			defaultAuthSessionLifetime,
		)
	}
}

func TestLoadBFFAuthConfigRejectsMissingEnvironment(
	t *testing.T,
) {
	requiredEnvironmentNames := []string{
		"COGNITO_CLIENT_ID",
		"COGNITO_CLIENT_SECRET_ARN",
		"COGNITO_ISSUER",
		"COGNITO_AUTHORIZATION_ENDPOINT",
		"COGNITO_TOKEN_ENDPOINT",
		"COGNITO_REDIRECT_URI",
		"AUTH_REDIRECT_URL",
		"COGNITO_OAUTH_STATE_TABLE_NAME",
		"AUTH_SESSION_TABLE_NAME",
	}

	for _, environmentName := range requiredEnvironmentNames {
		t.Run(
			environmentName,
			func(t *testing.T) {
				setValidBFFAuthEnvironment(t)

				t.Setenv(
					environmentName,
					"",
				)

				config, err :=
					loadBFFAuthConfig()
				if err == nil {
					t.Fatal(
						"loadBFFAuthConfig() error = nil, want an error",
					)
				}

				if config !=
					(BFFAuthConfig{}) {
					t.Errorf(
						"loadBFFAuthConfig() config = %+v, want zero value",
						config,
					)
				}

				wantError :=
					environmentName +
						" is required"

				if !strings.Contains(
					err.Error(),
					wantError,
				) {
					t.Errorf(
						"loadBFFAuthConfig() error = %q, want to contain %q",
						err,
						wantError,
					)
				}
			},
		)
	}
}
