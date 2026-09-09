package cognito

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
)

const (
	testClientID          = "test-cognito-client-id"
	testClientSecret      = "test-cognito-client-secret"
	testRedirectURI       = "https://example.com/api/auth/callback"
	testLogoutRedirectURI = "https://example.com/"

	testCode         = "test-authorization-code"
	testStateValue   = "test-oauth-state"
	testNonce        = "test-oidc-nonce"
	testIDToken      = "test-id-token"
	testAccessToken  = "test-access-token"
	testRefreshToken = "test-refresh-token"
)

type testTokenRequest struct {
	Calls int
	Form  url.Values

	HasBasicAuth bool
	ClientID     string
	ClientSecret string
}

type testOAuthServer struct {
	server *httptest.Server

	mu sync.Mutex

	tokenRequestCalls int
	tokenRequestForm  url.Values

	tokenRequestHasBasicAuth bool
	tokenRequestClientID     string
	tokenRequestClientSecret string

	tokenResponseStatus int
	tokenResponseBody   string
}

func newTestOAuthServer(t *testing.T) *testOAuthServer {
	t.Helper()

	oauthServer := &testOAuthServer{
		tokenResponseStatus: http.StatusOK,
		tokenResponseBody: `{
			"access_token": "test-access-token",
			"id_token": "test-id-token",
			"refresh_token": "test-refresh-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`,
	}

	oauthServer.server = httptest.NewServer(http.HandlerFunc(
		func(
			response http.ResponseWriter,
			request *http.Request,
		) {
			if request.URL.Path != "/oauth2/token" {
				http.NotFound(
					response,
					request,
				)

				return
			}

			if err := request.ParseForm(); err != nil {
				http.Error(
					response,
					"invalid form",
					http.StatusBadRequest,
				)

				return
			}

			clientID,
				clientSecret,
				hasBasicAuth := request.BasicAuth()

			oauthServer.mu.Lock()

			oauthServer.tokenRequestCalls++
			oauthServer.tokenRequestForm = cloneURLValues(request.Form)
			oauthServer.tokenRequestHasBasicAuth = hasBasicAuth
			oauthServer.tokenRequestClientID = clientID
			oauthServer.tokenRequestClientSecret = clientSecret

			status := oauthServer.tokenResponseStatus
			body := oauthServer.tokenResponseBody

			oauthServer.mu.Unlock()

			response.Header().Set(
				"Content-Type",
				"application/json",
			)
			response.WriteHeader(status)

			_, _ = response.Write([]byte(body))
		}),
	)

	t.Cleanup(oauthServer.server.Close)

	return oauthServer
}

func (server *testOAuthServer) authorizationEndpoint() string {
	return server.server.URL + "/oauth2/authorize"
}

func (server *testOAuthServer) tokenEndpoint() string {
	return server.server.URL + "/oauth2/token"
}

func (server *testOAuthServer) logoutEndpoint() string {
	return server.server.URL + "/logout"
}

func (server *testOAuthServer) setTokenResponse(
	status int,
	body string,
) {
	server.mu.Lock()
	defer server.mu.Unlock()

	server.tokenResponseStatus = status
	server.tokenResponseBody = body
}

func (server *testOAuthServer) tokenRequest() testTokenRequest {
	server.mu.Lock()
	defer server.mu.Unlock()

	return testTokenRequest{
		Calls:        server.tokenRequestCalls,
		Form:         cloneURLValues(server.tokenRequestForm),
		HasBasicAuth: server.tokenRequestHasBasicAuth,
		ClientID:     server.tokenRequestClientID,
		ClientSecret: server.tokenRequestClientSecret,
	}
}

func newTestClientConfig(server *testOAuthServer) Config {
	return Config{
		ClientID:              testClientID,
		ClientSecret:          testClientSecret,
		AuthorizationEndpoint: server.authorizationEndpoint(),
		TokenEndpoint:         server.tokenEndpoint(),
		LogoutEndpoint:        server.logoutEndpoint(),
		RedirectURI:           testRedirectURI,
		LogoutRedirectURI:     testLogoutRedirectURI,
		Scopes: []string{
			"openid",
			"email",
			"profile",
		},
	}
}

func newTestOAuthState() oauthstate.State {
	createdAt := time.Date(
		2026,
		time.August,
		30,
		12,
		34,
		56,
		0,
		time.UTC,
	)

	return oauthstate.State{
		Value:        testStateValue,
		CodeVerifier: strings.Repeat("v", 43),
		Nonce:        testNonce,
		CreatedAt:    createdAt,
		ExpiresAt:    createdAt.Add(10 * time.Minute),
	}
}

func cloneURLValues(values url.Values) url.Values {
	if values == nil {
		return nil
	}

	cloned := make(url.Values, len(values))

	for key, entries := range values {
		cloned[key] = append(
			[]string(nil),
			entries...,
		)
	}

	return cloned
}
