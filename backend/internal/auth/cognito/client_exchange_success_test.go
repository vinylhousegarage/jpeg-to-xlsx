package cognito

import (
	"context"
	"testing"
)

func TestClientExchange(
	t *testing.T,
) {
	t.Parallel()

	server := newTestOAuthServer(t)

	client, err := NewClient(
		newTestClientConfig(server),
	)
	if err != nil {
		t.Fatalf(
			"NewClient() error = %v",
			err,
		)
	}

	state := newTestOAuthState()

	idToken, err := client.Exchange(
		context.Background(),
		testCode,
		state.CodeVerifier,
	)
	if err != nil {
		t.Fatalf(
			"Exchange() error = %v",
			err,
		)
	}

	if idToken != testIDToken {
		t.Errorf(
			"Exchange() ID token = %q, want %q",
			idToken,
			testIDToken,
		)
	}

	request := server.tokenRequest()

	if request.Calls != 1 {
		t.Fatalf(
			"token endpoint calls = %d, want 1",
			request.Calls,
		)
	}

	if !request.HasBasicAuth {
		t.Error(
			"token request does not use HTTP Basic authentication",
		)
	}

	if request.ClientID != testClientID {
		t.Errorf(
			"token request Basic Auth client ID = %q, want %q",
			request.ClientID,
			testClientID,
		)
	}

	if request.ClientSecret != testClientSecret {
		t.Errorf(
			"token request Basic Auth client secret = %q, want configured secret",
			request.ClientSecret,
		)
	}

	assertTokenRequestFormValue(
		t,
		request,
		"grant_type",
		"authorization_code",
	)
	assertTokenRequestFormValue(
		t,
		request,
		"code",
		testCode,
	)
	assertTokenRequestFormValue(
		t,
		request,
		"redirect_uri",
		testRedirectURI,
	)
	assertTokenRequestFormValue(
		t,
		request,
		"code_verifier",
		state.CodeVerifier,
	)

	if request.Form.Has("client_secret") {
		t.Error(
			"token request form contains client_secret",
		)
	}

	if request.Form.Has("client_id") {
		t.Error(
			"token request form contains client_id",
		)
	}
}

func assertTokenRequestFormValue(
	t *testing.T,
	request testTokenRequest,
	name string,
	want string,
) {
	t.Helper()

	values, exists := request.Form[name]
	if !exists {
		t.Errorf(
			"token request form value %q is missing",
			name,
		)

		return
	}

	if len(values) != 1 {
		t.Errorf(
			"token request form value %q has %d values, want 1",
			name,
			len(values),
		)

		return
	}

	if values[0] != want {
		t.Errorf(
			"token request form value %q = %q, want %q",
			name,
			values[0],
			want,
		)
	}
}
