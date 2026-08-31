package callback

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
)

type callbackError struct {
	status int
	err    error
}

func (
	callbackErr *callbackError,
) Error() string {
	if callbackErr == nil ||
		callbackErr.err == nil {
		return "Cognito callback failed"
	}

	return callbackErr.err.Error()
}

func (
	callbackErr *callbackError,
) Unwrap() error {
	if callbackErr == nil {
		return nil
	}

	return callbackErr.err
}

func newCallbackError(
	status int,
	err error,
) *callbackError {
	return &callbackError{
		status: status,
		err:    err,
	}
}

func (
	handler *Handler,
) completeAuthentication(
	ctx context.Context,
	response http.ResponseWriter,
	code string,
	stateValue string,
) *callbackError {
	state, err := handler.oauthStateStore.Consume(
		ctx,
		stateValue,
	)
	if err != nil {
		if errors.Is(
			err,
			oauthstate.ErrNotFound,
		) {
			return newCallbackError(
				http.StatusBadRequest,
				fmt.Errorf(
					"complete Cognito authentication: consume OAuth state: %w",
					err,
				),
			)
		}

		return newCallbackError(
			http.StatusInternalServerError,
			fmt.Errorf(
				"complete Cognito authentication: consume OAuth state: %w",
				err,
			),
		)
	}

	if state.Value != stateValue {
		return newCallbackError(
			http.StatusBadRequest,
			fmt.Errorf(
				"complete Cognito authentication: OAuth state does not match",
			),
		)
	}

	rawIDToken, err :=
		handler.cognitoClient.Exchange(
			ctx,
			code,
			state.CodeVerifier,
		)
	if err != nil {
		return newCallbackError(
			http.StatusBadGateway,
			fmt.Errorf(
				"complete Cognito authentication: exchange authorization code: %w",
				err,
			),
		)
	}

	identity, err := handler.verifier.Verify(
		ctx,
		rawIDToken,
		state.Nonce,
	)
	if err != nil {
		return newCallbackError(
			http.StatusUnauthorized,
			fmt.Errorf(
				"complete Cognito authentication: verify ID token: %w",
				err,
			),
		)
	}

	if strings.TrimSpace(
		identity.Subject,
	) == "" {
		return newCallbackError(
			http.StatusUnauthorized,
			fmt.Errorf(
				"complete Cognito authentication: Cognito subject is empty",
			),
		)
	}

	rawSessionID, err :=
		handler.idGenerator.Generate()
	if err != nil {
		return newCallbackError(
			http.StatusInternalServerError,
			fmt.Errorf(
				"complete Cognito authentication: generate session ID: %w",
				err,
			),
		)
	}

	idHash, err := session.HashID(
		rawSessionID,
	)
	if err != nil {
		return newCallbackError(
			http.StatusInternalServerError,
			fmt.Errorf(
				"complete Cognito authentication: hash session ID: %w",
				err,
			),
		)
	}

	now := handler.now().UTC()

	sessionValue := session.Session{
		IDHash:     idHash,
		CognitoSub: identity.Subject,
		CreatedAt:  now,
		ExpiresAt: now.Add(
			handler.sessionLifetime,
		),
	}

	if err := handler.sessionStore.Save(
		ctx,
		sessionValue,
	); err != nil {
		return newCallbackError(
			http.StatusInternalServerError,
			fmt.Errorf(
				"complete Cognito authentication: save session: %w",
				err,
			),
		)
	}

	if err := handler.cookieWriter.Write(
		response,
		rawSessionID,
	); err != nil {
		deleteErr := handler.sessionStore.Delete(
			ctx,
			idHash,
		)

		if deleteErr != nil {
			return newCallbackError(
				http.StatusInternalServerError,
				errors.Join(
					fmt.Errorf(
						"complete Cognito authentication: write session cookie: %w",
						err,
					),
					fmt.Errorf(
						"complete Cognito authentication: roll back session: %w",
						deleteErr,
					),
				),
			)
		}

		return newCallbackError(
			http.StatusInternalServerError,
			fmt.Errorf(
				"complete Cognito authentication: write session cookie: %w",
				err,
			),
		)
	}

	return nil
}
