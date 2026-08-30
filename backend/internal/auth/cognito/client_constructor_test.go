package cognito

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	server := newTestOAuthServer(t)
	config := newTestClientConfig(server)

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf(
			"NewClient() error = %v",
			err,
		)
	}

	if client == nil {
		t.Fatal(
			"NewClient() client = nil, want non-nil",
		)
	}
}

func TestNewClientRejectsInvalidConfig(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*Config)
	}{
		{
			name: "empty client ID",
			mutate: func(config *Config) {
				config.ClientID = ""
			},
		},
		{
			name: "whitespace client ID",
			mutate: func(config *Config) {
				config.ClientID = "   "
			},
		},
		{
			name: "empty client secret",
			mutate: func(config *Config) {
				config.ClientSecret = ""
			},
		},
		{
			name: "whitespace client secret",
			mutate: func(config *Config) {
				config.ClientSecret = "   "
			},
		},
		{
			name: "empty authorization endpoint",
			mutate: func(config *Config) {
				config.AuthorizationEndpoint = ""
			},
		},
		{
			name: "invalid authorization endpoint",
			mutate: func(config *Config) {
				config.AuthorizationEndpoint =
					"://invalid"
			},
		},
		{
			name: "authorization endpoint without scheme",
			mutate: func(config *Config) {
				config.AuthorizationEndpoint =
					"example.com/oauth2/authorize"
			},
		},
		{
			name: "authorization endpoint without host",
			mutate: func(config *Config) {
				config.AuthorizationEndpoint =
					"https:///oauth2/authorize"
			},
		},
		{
			name: "empty token endpoint",
			mutate: func(config *Config) {
				config.TokenEndpoint = ""
			},
		},
		{
			name: "invalid token endpoint",
			mutate: func(config *Config) {
				config.TokenEndpoint =
					"://invalid"
			},
		},
		{
			name: "token endpoint without scheme",
			mutate: func(config *Config) {
				config.TokenEndpoint =
					"example.com/oauth2/token"
			},
		},
		{
			name: "token endpoint without host",
			mutate: func(config *Config) {
				config.TokenEndpoint =
					"https:///oauth2/token"
			},
		},
		{
			name: "empty redirect URI",
			mutate: func(config *Config) {
				config.RedirectURI = ""
			},
		},
		{
			name: "invalid redirect URI",
			mutate: func(config *Config) {
				config.RedirectURI =
					"://invalid"
			},
		},
		{
			name: "redirect URI without scheme",
			mutate: func(config *Config) {
				config.RedirectURI =
					"example.com/api/auth/callback"
			},
		},
		{
			name: "redirect URI without host",
			mutate: func(config *Config) {
				config.RedirectURI =
					"https:///api/auth/callback"
			},
		},
		{
			name: "nil scopes",
			mutate: func(config *Config) {
				config.Scopes = nil
			},
		},
		{
			name: "empty scopes",
			mutate: func(config *Config) {
				config.Scopes = []string{}
			},
		},
		{
			name: "scope contains empty value",
			mutate: func(config *Config) {
				config.Scopes = []string{
					"openid",
					"",
				}
			},
		},
		{
			name: "scope contains whitespace value",
			mutate: func(config *Config) {
				config.Scopes = []string{
					"openid",
					"   ",
				}
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				server := newTestOAuthServer(t)
				config := newTestClientConfig(
					server,
				)

				test.mutate(&config)

				client, err := NewClient(config)
				if err == nil {
					t.Fatal(
						"NewClient() error = nil, want an error",
					)
				}

				if client != nil {
					t.Errorf(
						"NewClient() client = %v, want nil",
						client,
					)
				}
			},
		)
	}
}
