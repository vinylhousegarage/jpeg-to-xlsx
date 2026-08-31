package config

import "testing"

func TestLoadSlackConfigWithClientSecret(
	t *testing.T,
) {
	setValidSlackEnvironment(t)

	config, err := loadSlackConfig()
	if err != nil {
		t.Fatalf(
			"loadSlackConfig() error = %v",
			err,
		)
	}

	if config.ClientID !=
		testSlackClientID {
		t.Errorf(
			"ClientID = %q, want %q",
			config.ClientID,
			testSlackClientID,
		)
	}

	if config.ClientSecret !=
		testSlackClientSecret {
		t.Errorf(
			"ClientSecret = %q, want %q",
			config.ClientSecret,
			testSlackClientSecret,
		)
	}

	if config.ClientSecretARN != "" {
		t.Errorf(
			"ClientSecretARN = %q, want empty",
			config.ClientSecretARN,
		)
	}

	if config.RedirectURI !=
		testSlackRedirectURI {
		t.Errorf(
			"RedirectURI = %q, want %q",
			config.RedirectURI,
			testSlackRedirectURI,
		)
	}

	if config.TokenTableName !=
		testSlackTokenTableName {
		t.Errorf(
			"TokenTableName = %q, want %q",
			config.TokenTableName,
			testSlackTokenTableName,
		)
	}
}

func TestLoadSlackConfigWithClientSecretARN(
	t *testing.T,
) {
	const wantARN = "arn:aws:secretsmanager:" +
		"ap-northeast-1:123456789012:" +
		"secret:test"

	setValidSlackEnvironment(t)

	t.Setenv(
		"SLACK_CLIENT_SECRET",
		"",
	)
	t.Setenv(
		"SLACK_CLIENT_SECRET_ARN",
		wantARN,
	)

	config, err := loadSlackConfig()
	if err != nil {
		t.Fatalf(
			"loadSlackConfig() error = %v",
			err,
		)
	}

	if config.ClientSecret != "" {
		t.Errorf(
			"ClientSecret = %q, want empty",
			config.ClientSecret,
		)
	}

	if config.ClientSecretARN != wantARN {
		t.Errorf(
			"ClientSecretARN = %q, want %q",
			config.ClientSecretARN,
			wantARN,
		)
	}

	if config.ClientID !=
		testSlackClientID {
		t.Errorf(
			"ClientID = %q, want %q",
			config.ClientID,
			testSlackClientID,
		)
	}

	if config.RedirectURI !=
		testSlackRedirectURI {
		t.Errorf(
			"RedirectURI = %q, want %q",
			config.RedirectURI,
			testSlackRedirectURI,
		)
	}

	if config.TokenTableName !=
		testSlackTokenTableName {
		t.Errorf(
			"TokenTableName = %q, want %q",
			config.TokenTableName,
			testSlackTokenTableName,
		)
	}
}
