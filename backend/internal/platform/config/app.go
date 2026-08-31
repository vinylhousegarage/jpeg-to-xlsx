package config

import "os"

func loadAppConfig() AppConfig {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = appEnvLocal
	}

	return AppConfig{
		Env: env,
		CookieSecure: env == appEnvProduction ||
			env == appEnvStaging,
	}
}
