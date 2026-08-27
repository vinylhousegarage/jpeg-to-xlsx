package oauth

import (
	"net/url"
	"testing"
)

func TestBuildAuthURL(t *testing.T) {
	t.Parallel()

	const (
		clientID    = "test-client-id"
		redirectURI = "https://example.com/callback"
		state       = "test-state"
	)

	authURLString := BuildAuthURL(
		clientID,
		redirectURI,
		state,
	)

	parsedURL, err := url.Parse(authURLString)
	if err != nil {
		t.Fatalf(
			"failed to parse generated auth URL: %v",
			err,
		)
	}

	if parsedURL.Scheme != "https" {
		t.Errorf(
			"expected scheme %q, got %q",
			"https",
			parsedURL.Scheme,
		)
	}

	if parsedURL.Host != "slack.com" {
		t.Errorf(
			"expected host %q, got %q",
			"slack.com",
			parsedURL.Host,
		)
	}

	if parsedURL.Path != "/oauth/v2/authorize" {
		t.Errorf(
			"expected path %q, got %q",
			"/oauth/v2/authorize",
			parsedURL.Path,
		)
	}

	query := parsedURL.Query()

	if query.Get("client_id") != clientID {
		t.Errorf(
			"expected client_id %q, got %q",
			clientID,
			query.Get("client_id"),
		)
	}

	if query.Get("scope") != slackBotScopes {
		t.Errorf(
			"expected scope %q, got %q",
			slackBotScopes,
			query.Get("scope"),
		)
	}

	if query.Get("redirect_uri") != redirectURI {
		t.Errorf(
			"expected redirect_uri %q, got %q",
			redirectURI,
			query.Get("redirect_uri"),
		)
	}

	if query.Get("state") != state {
		t.Errorf(
			"expected state %q, got %q",
			state,
			query.Get("state"),
		)
	}
}
