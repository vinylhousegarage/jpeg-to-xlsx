package callback

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	RedirectURL     string
	SessionLifetime time.Duration
}

type Handler struct {
	oauthStateStore OAuthStateStore
	cognitoClient   CognitoClient
	verifier        IDTokenVerifier
	idGenerator     SessionIDGenerator
	sessionStore    SessionStore
	cookieWriter    SessionCookieWriter

	redirectURL     string
	sessionLifetime time.Duration
	now             func() time.Time
}

func NewHandler(
	oauthStateStore OAuthStateStore,
	cognitoClient CognitoClient,
	verifier IDTokenVerifier,
	idGenerator SessionIDGenerator,
	sessionStore SessionStore,
	cookieWriter SessionCookieWriter,
	config Config,
) (
	*Handler,
	error,
) {
	if oauthStateStore == nil {
		return nil, fmt.Errorf(
			"create Cognito callback handler: OAuth state store is nil",
		)
	}

	if cognitoClient == nil {
		return nil, fmt.Errorf(
			"create Cognito callback handler: Cognito client is nil",
		)
	}

	if verifier == nil {
		return nil, fmt.Errorf(
			"create Cognito callback handler: ID token verifier is nil",
		)
	}

	if idGenerator == nil {
		return nil, fmt.Errorf(
			"create Cognito callback handler: session ID generator is nil",
		)
	}

	if sessionStore == nil {
		return nil, fmt.Errorf(
			"create Cognito callback handler: session store is nil",
		)
	}

	if cookieWriter == nil {
		return nil, fmt.Errorf(
			"create Cognito callback handler: cookie writer is nil",
		)
	}

	redirectURL, err := validateRedirectURL(
		config.RedirectURL,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create Cognito callback handler: %w",
			err,
		)
	}

	if config.SessionLifetime < time.Second {
		return nil, fmt.Errorf(
			"create Cognito callback handler: session lifetime must be at least one second",
		)
	}

	return &Handler{
		oauthStateStore: oauthStateStore,
		cognitoClient:   cognitoClient,
		verifier:        verifier,
		idGenerator:     idGenerator,
		sessionStore:    sessionStore,
		cookieWriter:    cookieWriter,
		redirectURL:     redirectURL,
		sessionLifetime: config.SessionLifetime,
		now:             time.Now,
	}, nil
}

func (
	handler *Handler,
) ServeHTTP(
	response http.ResponseWriter,
	request *http.Request,
) {
	if response == nil {
		return
	}

	if handler == nil {
		writeError(
			response,
			http.StatusInternalServerError,
		)

		return
	}

	if request == nil ||
		request.URL == nil {
		writeError(
			response,
			http.StatusBadRequest,
		)

		return
	}

	if request.Method != http.MethodGet {
		response.Header().Set(
			"Allow",
			http.MethodGet,
		)

		writeError(
			response,
			http.StatusMethodNotAllowed,
		)

		return
	}

	query := request.URL.Query()

	if strings.TrimSpace(
		query.Get("error"),
	) != "" {
		writeError(
			response,
			http.StatusBadRequest,
		)

		return
	}

	code := query.Get("code")
	if strings.TrimSpace(code) == "" {
		writeError(
			response,
			http.StatusBadRequest,
		)

		return
	}

	stateValue := query.Get("state")
	if strings.TrimSpace(stateValue) == "" {
		writeError(
			response,
			http.StatusBadRequest,
		)

		return
	}

	if handler.now == nil {
		writeError(
			response,
			http.StatusInternalServerError,
		)

		return
	}

	callbackErr :=
		handler.completeAuthentication(
			request.Context(),
			response,
			code,
			stateValue,
		)
	if callbackErr != nil {
		writeError(
			response,
			callbackErr.status,
		)

		return
	}

	writeRedirect(
		response,
		handler.redirectURL,
	)
}
