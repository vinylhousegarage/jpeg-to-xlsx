package cognito

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-jose/go-jose/v4"
)

const testOIDCKeyID = "test-cognito-key-id"

type testOIDCServer struct {
	server     *httptest.Server
	privateKey *rsa.PrivateKey

	mu sync.Mutex

	discoveryRequestCalls int
	jwksRequestCalls      int

	discoveryResponseStatus int
	discoveryResponseBody   string
	jwksResponseStatus      int
	jwksResponseBody        string
}

type testOIDCRequests struct {
	DiscoveryCalls int
	JWKSCalls      int
}

func newTestOIDCServer(
	t *testing.T,
) *testOIDCServer {
	t.Helper()

	privateKey, err := rsa.GenerateKey(
		rand.Reader,
		2048,
	)
	if err != nil {
		t.Fatalf(
			"generate test RSA key: %v",
			err,
		)
	}

	oidcServer := &testOIDCServer{
		privateKey:              privateKey,
		discoveryResponseStatus: http.StatusOK,
		jwksResponseStatus:      http.StatusOK,
	}

	oidcServer.server = httptest.NewServer(
		http.HandlerFunc(
			func(
				response http.ResponseWriter,
				request *http.Request,
			) {
				switch request.URL.Path {
				case "/.well-known/openid-configuration":
					oidcServer.handleDiscovery(
						response,
					)

				case "/.well-known/jwks.json":
					oidcServer.handleJWKS(
						response,
					)

				default:
					http.NotFound(
						response,
						request,
					)
				}
			},
		),
	)

	t.Cleanup(
		oidcServer.server.Close,
	)

	return oidcServer
}

func (
	server *testOIDCServer,
) issuer() string {
	return server.server.URL
}

func (
	server *testOIDCServer,
) jwksEndpoint() string {
	return server.issuer() +
		"/.well-known/jwks.json"
}

func (
	server *testOIDCServer,
) setDiscoveryResponse(
	status int,
	body string,
) {
	server.mu.Lock()
	defer server.mu.Unlock()

	server.discoveryResponseStatus = status
	server.discoveryResponseBody = body
}

func (
	server *testOIDCServer,
) setJWKSResponse(
	status int,
	body string,
) {
	server.mu.Lock()
	defer server.mu.Unlock()

	server.jwksResponseStatus = status
	server.jwksResponseBody = body
}

func (
	server *testOIDCServer,
) requests() testOIDCRequests {
	server.mu.Lock()
	defer server.mu.Unlock()

	return testOIDCRequests{
		DiscoveryCalls: server.discoveryRequestCalls,
		JWKSCalls:      server.jwksRequestCalls,
	}
}

func (
	server *testOIDCServer,
) handleDiscovery(
	response http.ResponseWriter,
) {
	server.mu.Lock()

	server.discoveryRequestCalls++

	status := server.discoveryResponseStatus
	body := server.discoveryResponseBody

	server.mu.Unlock()

	response.Header().Set(
		"Content-Type",
		"application/json",
	)
	response.WriteHeader(status)

	if body != "" {
		_, _ = response.Write(
			[]byte(body),
		)

		return
	}

	discovery := map[string]any{
		"issuer": server.issuer(),
		"authorization_endpoint": server.issuer() +
			"/oauth2/authorize",
		"token_endpoint": server.issuer() +
			"/oauth2/token",
		"jwks_uri": server.jwksEndpoint(),
		"response_types_supported": []string{
			"code",
		},
		"subject_types_supported": []string{
			"public",
		},
		"id_token_signing_alg_values_supported": []string{
			"RS256",
		},
	}

	_ = json.NewEncoder(
		response,
	).Encode(discovery)
}

func (
	server *testOIDCServer,
) handleJWKS(
	response http.ResponseWriter,
) {
	server.mu.Lock()

	server.jwksRequestCalls++

	status := server.jwksResponseStatus
	body := server.jwksResponseBody

	server.mu.Unlock()

	response.Header().Set(
		"Content-Type",
		"application/json",
	)
	response.WriteHeader(status)

	if body != "" {
		_, _ = response.Write(
			[]byte(body),
		)

		return
	}

	keySet := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:       &server.privateKey.PublicKey,
				KeyID:     testOIDCKeyID,
				Algorithm: string(jose.RS256),
				Use:       "sig",
			},
		},
	}

	_ = json.NewEncoder(
		response,
	).Encode(keySet)
}
