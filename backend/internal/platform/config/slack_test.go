package config

import (
	"strings"
	"testing"
)

func TestLoadSlackConfig(t *testing.T) {
	t.Run("success with client secret", func(t *testing.T) {
		t.Setenv("SLACK_CLIENT_ID", "test-client-id")
		t.Setenv("SLACK_CLIENT_SECRET", "test-client-secret")
		t.Setenv("SLACK_CLIENT_SECRET_ARN", "")
		t.Setenv(
			"SLACK_REDIRECT_URI",
			"https://example.com/oauth/slack/callback",
		)
		t.Setenv(
			"SLACK_TOKEN_TABLE_NAME",
			"slack-tokens",
		)

		cfg, err := loadSlackConfig()
		if err != nil {
			t.Fatalf("loadSlackConfig() error = %v", err)
		}

		if cfg.ClientID != "test-client-id" {
			t.Errorf(
				"expected ClientID %q, got %q",
				"test-client-id",
				cfg.ClientID,
			)
		}

		if cfg.ClientSecret != "test-client-secret" {
			t.Errorf(
				"expected ClientSecret %q, got %q",
				"test-client-secret",
				cfg.ClientSecret,
			)
		}

		if cfg.ClientSecretARN != "" {
			t.Errorf(
				"expected ClientSecretARN to be empty, got %q",
				cfg.ClientSecretARN,
			)
		}

		expectedRedirectURI :=
			"https://example.com/oauth/slack/callback"

		if cfg.RedirectURI != expectedRedirectURI {
			t.Errorf(
				"expected RedirectURI %q, got %q",
				expectedRedirectURI,
				cfg.RedirectURI,
			)
		}

		if cfg.TokenTableName != "slack-tokens" {
			t.Errorf(
				"expected TokenTableName %q, got %q",
				"slack-tokens",
				cfg.TokenTableName,
			)
		}
	})

	tests := []struct {
		name            string
		clientID        string
		clientSecret    string
		clientSecretARN string
		redirectURI     string
		tokenTableName  string
		expectedErr     string
	}{
		{
			name:            "missing client ID",
			clientID:        "",
			clientSecret:    "test-client-secret",
			clientSecretARN: "",
			redirectURI:     "https://example.com/oauth/slack/callback",
			tokenTableName:  "slack-tokens",
			expectedErr:     "SLACK_CLIENT_ID is required",
		},
		{
			name:            "missing client secret and ARN",
			clientID:        "test-client-id",
			clientSecret:    "",
			clientSecretARN: "",
			redirectURI:     "https://example.com/oauth/slack/callback",
			tokenTableName:  "slack-tokens",
			expectedErr:     "SLACK_CLIENT_SECRET or SLACK_CLIENT_SECRET_ARN is required",
		},
		{
			name:            "missing redirect URI",
			clientID:        "test-client-id",
			clientSecret:    "test-client-secret",
			clientSecretARN: "",
			redirectURI:     "",
			tokenTableName:  "slack-tokens",
			expectedErr:     "SLACK_REDIRECT_URI is required",
		},
		{
			name:            "missing token table name",
			clientID:        "test-client-id",
			clientSecret:    "test-client-secret",
			clientSecretARN: "",
			redirectURI:     "https://example.com/oauth/slack/callback",
			tokenTableName:  "",
			expectedErr:     "SLACK_TOKEN_TABLE_NAME is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SLACK_CLIENT_ID", tt.clientID)
			t.Setenv(
				"SLACK_CLIENT_SECRET",
				tt.clientSecret,
			)
			t.Setenv(
				"SLACK_CLIENT_SECRET_ARN",
				tt.clientSecretARN,
			)
			t.Setenv(
				"SLACK_REDIRECT_URI",
				tt.redirectURI,
			)
			t.Setenv(
				"SLACK_TOKEN_TABLE_NAME",
				tt.tokenTableName,
			)

			_, err := loadSlackConfig()
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(
				err.Error(),
				tt.expectedErr,
			) {
				t.Errorf(
					"expected error containing %q, got %q",
					tt.expectedErr,
					err.Error(),
				)
			}
		})
	}

	t.Run("success with client secret ARN", func(t *testing.T) {
		const wantARN = "arn:aws:secretsmanager:ap-northeast-1:123456789012:secret:test"

		t.Setenv("SLACK_CLIENT_ID", "test-client-id")
		t.Setenv("SLACK_CLIENT_SECRET", "")
		t.Setenv("SLACK_CLIENT_SECRET_ARN", wantARN)
		t.Setenv(
			"SLACK_REDIRECT_URI",
			"https://example.com/api/oauth/slack/callback",
		)
		t.Setenv("SLACK_TOKEN_TABLE_NAME", "slack-tokens")

		cfg, err := loadSlackConfig()
		if err != nil {
			t.Fatalf("loadSlackConfig() error = %v", err)
		}

		if cfg.ClientSecret != "" {
			t.Errorf(
				"expected ClientSecret to be empty, got %q",
				cfg.ClientSecret,
			)
		}

		if cfg.ClientSecretARN != wantARN {
			t.Errorf(
				"expected ClientSecretARN %q, got %q",
				wantARN,
				cfg.ClientSecretARN,
			)
		}
	})
}
