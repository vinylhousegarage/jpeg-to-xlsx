package config

import "testing"

func TestLoadSlackConfigRejectsMissingEnvironment(
	t *testing.T,
) {
	tests := []struct {
		name      string
		configure func(*testing.T)
		wantError string
	}{
		{
			name: "missing client ID",
			configure: func(t *testing.T) {
				t.Helper()

				t.Setenv(
					"SLACK_CLIENT_ID",
					"",
				)
			},
			wantError: "SLACK_CLIENT_ID is required",
		},
		{
			name: "missing client secret and ARN",
			configure: func(t *testing.T) {
				t.Helper()

				t.Setenv(
					"SLACK_CLIENT_SECRET",
					"",
				)
				t.Setenv(
					"SLACK_CLIENT_SECRET_ARN",
					"",
				)
			},
			wantError: "SLACK_CLIENT_SECRET or " +
				"SLACK_CLIENT_SECRET_ARN is required",
		},
		{
			name: "missing redirect URI",
			configure: func(t *testing.T) {
				t.Helper()

				t.Setenv(
					"SLACK_REDIRECT_URI",
					"",
				)
			},
			wantError: "SLACK_REDIRECT_URI is required",
		},
		{
			name: "missing token table name",
			configure: func(t *testing.T) {
				t.Helper()

				t.Setenv(
					"SLACK_TOKEN_TABLE_NAME",
					"",
				)
			},
			wantError: "SLACK_TOKEN_TABLE_NAME is required",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				setValidSlackEnvironment(t)
				test.configure(t)

				config, err :=
					loadSlackConfig()
				if err == nil {
					t.Fatal(
						"loadSlackConfig() error = nil, want an error",
					)
				}

				if config !=
					(SlackConfig{}) {
					t.Errorf(
						"loadSlackConfig() config = %+v, want zero value",
						config,
					)
				}

				if err.Error() !=
					test.wantError {
					t.Errorf(
						"loadSlackConfig() error = %q, want %q",
						err,
						test.wantError,
					)
				}
			},
		)
	}
}
