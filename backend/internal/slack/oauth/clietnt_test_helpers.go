package oauth

import "net/http/httptest"

const (
	testClientID     = "client-id"
	testClientSecret = "client-secret"
	testCode         = "code123"
	testRedirectURI  = "https://example.com/callback"
)

func newTestClient(
	server *httptest.Server,
) *Client {
	client := NewClient(
		server.Client(),
		testClientID,
		testClientSecret,
	)

	client.tokenURL = server.URL

	return client
}
