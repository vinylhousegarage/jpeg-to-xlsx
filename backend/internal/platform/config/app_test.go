package config

import "testing"

func TestLoadAppConfig(t *testing.T) {
	tests := []struct {
		name         string
		env          string
		expectedEnv  string
		cookieSecure bool
	}{
		{
			name:         "default local",
			env:          "",
			expectedEnv:  appEnvLocal,
			cookieSecure: false,
		},
		{
			name:         "local",
			env:          appEnvLocal,
			expectedEnv:  appEnvLocal,
			cookieSecure: false,
		},
		{
			name:         "staging",
			env:          appEnvStaging,
			expectedEnv:  appEnvStaging,
			cookieSecure: true,
		},
		{
			name:         "production",
			env:          appEnvProduction,
			expectedEnv:  appEnvProduction,
			cookieSecure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tt.env)

			cfg := loadAppConfig()

			if cfg.Env != tt.expectedEnv {
				t.Errorf(
					"expected Env %q, got %q",
					tt.expectedEnv,
					cfg.Env,
				)
			}

			if cfg.CookieSecure != tt.cookieSecure {
				t.Errorf(
					"expected CookieSecure %v, got %v",
					tt.cookieSecure,
					cfg.CookieSecure,
				)
			}
		})
	}
}
