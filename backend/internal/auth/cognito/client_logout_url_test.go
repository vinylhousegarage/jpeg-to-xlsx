package cognito

import (
	"net/url"
	"strings"
	"testing"
)

func TestClientLogoutURL(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOAuthServer(t)
	config := newTestClientConfig(server)

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	logoutURL, err := client.LogoutURL()
	if err != nil {
		t.Fatalf("LogoutURL() error = %v", err)
	}

	parsedURL, err := url.Parse(logoutURL)
	if err != nil {
		t.Fatalf("parse logout URL: %v", err)
	}

	expectedEndpoint, err := url.Parse(
		server.logoutEndpoint(),
	)
	if err != nil {
		t.Fatalf(
			"parse expected logout endpoint: %v",
			err,
		)
	}

	if parsedURL.Scheme != expectedEndpoint.Scheme {
		t.Errorf(
			"logout URL scheme = %q, want %q",
			parsedURL.Scheme,
			expectedEndpoint.Scheme,
		)
	}

	if parsedURL.Host != expectedEndpoint.Host {
		t.Errorf(
			"logout URL host = %q, want %q",
			parsedURL.Host,
			expectedEndpoint.Host,
		)
	}

	if parsedURL.Path != expectedEndpoint.Path {
		t.Errorf(
			"logout URL path = %q, want %q",
			parsedURL.Path,
			expectedEndpoint.Path,
		)
	}

	query := parsedURL.Query()

	assertLogoutQueryValue(
		t,
		query,
		"client_id",
		testClientID,
	)
	assertLogoutQueryValue(
		t,
		query,
		"logout_uri",
		testLogoutRedirectURI,
	)

	if query.Has("client_secret") {
		t.Error(
			"logout URL contains client_secret query parameter",
		)
	}

	if strings.Contains(
		logoutURL,
		testClientSecret,
	) {
		t.Error(
			"logout URL contains the client secret",
		)
	}
}

func TestClientLogoutURLRejectsNilClient(
	t *testing.T,
) {
	t.Parallel()

	var client *Client

	logoutURL, err := client.LogoutURL()
	if err == nil {
		t.Fatal(
			"LogoutURL() error = nil, want an error",
		)
	}

	if logoutURL != "" {
		t.Errorf(
			"LogoutURL() URL = %q, want empty",
			logoutURL,
		)
	}
}

func assertLogoutQueryValue(
	t *testing.T,
	query url.Values,
	name string,
	want string,
) {
	t.Helper()

	values, exists := query[name]
	if !exists {
		t.Errorf(
			"logout URL query %q is missing",
			name,
		)

		return
	}

	if len(values) != 1 {
		t.Errorf(
			"logout URL query %q has %d values, want 1",
			name,
			len(values),
		)

		return
	}

	if values[0] != want {
		t.Errorf(
			"logout URL query %q = %q, want %q",
			name,
			values[0],
			want,
		)
	}
}
