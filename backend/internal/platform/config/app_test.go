package config

import "testing"

func TestLoadAppConfig(
	t *testing.T,
) {
	tests := []struct {
		name             string
		environment      string
		wantEnvironment  string
		wantCookieSecure bool
	}{
		{
			name:             "default local",
			environment:      "",
			wantEnvironment:  appEnvLocal,
			wantCookieSecure: false,
		},
		{
			name:             "local",
			environment:      appEnvLocal,
			wantEnvironment:  appEnvLocal,
			wantCookieSecure: false,
		},
		{
			name:             "staging",
			environment:      appEnvStaging,
			wantEnvironment:  appEnvStaging,
			wantCookieSecure: true,
		},
		{
			name:             "production",
			environment:      appEnvProduction,
			wantEnvironment:  appEnvProduction,
			wantCookieSecure: true,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Setenv(
					"APP_ENV",
					test.environment,
				)

				config := loadAppConfig()

				if config.Env !=
					test.wantEnvironment {
					t.Errorf(
						"Env = %q, want %q",
						config.Env,
						test.wantEnvironment,
					)
				}

				if config.CookieSecure !=
					test.wantCookieSecure {
					t.Errorf(
						"CookieSecure = %v, want %v",
						config.CookieSecure,
						test.wantCookieSecure,
					)
				}
			},
		)
	}
}
