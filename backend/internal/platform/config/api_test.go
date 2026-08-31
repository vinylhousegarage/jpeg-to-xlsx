package config

import "testing"

func TestLoadAPI(
	t *testing.T,
) {
	setValidAPIEnvironment(t)

	config, err := LoadAPI()
	if err != nil {
		t.Fatalf(
			"LoadAPI() error = %v",
			err,
		)
	}

	if config == nil {
		t.Fatal(
			"LoadAPI() config = nil",
		)
	}

	if config.App.Env != appEnvLocal {
		t.Errorf(
			"App.Env = %q, want %q",
			config.App.Env,
			appEnvLocal,
		)
	}

	if config.AWS.Region !=
		defaultAWSRegion {
		t.Errorf(
			"AWS.Region = %q, want %q",
			config.AWS.Region,
			defaultAWSRegion,
		)
	}

	if config.Auth.CognitoClientID !=
		testCognitoClientID {
		t.Errorf(
			"Auth.CognitoClientID = %q, want %q",
			config.Auth.CognitoClientID,
			testCognitoClientID,
		)
	}

	if config.Slack.ClientID !=
		testSlackClientID {
		t.Errorf(
			"Slack.ClientID = %q, want %q",
			config.Slack.ClientID,
			testSlackClientID,
		)
	}

	if config.Storage.InputBucketName !=
		testInputBucketName {
		t.Errorf(
			"Storage.InputBucketName = %q, want %q",
			config.Storage.InputBucketName,
			testInputBucketName,
		)
	}
}
