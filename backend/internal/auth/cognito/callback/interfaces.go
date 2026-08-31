package callback

import (
	"context"
	"net/http"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
)

type OAuthStateStore interface {
	Consume(
		ctx context.Context,
		value string,
	) (
		oauthstate.State,
		error,
	)
}

type CognitoClient interface {
	Exchange(
		ctx context.Context,
		code string,
		codeVerifier string,
	) (
		string,
		error,
	)
}

type IDTokenVerifier interface {
	Verify(
		ctx context.Context,
		rawIDToken string,
		expectedNonce string,
	) (
		cognito.Identity,
		error,
	)
}

type SessionIDGenerator interface {
	Generate() (
		string,
		error,
	)
}

type SessionStore interface {
	Save(
		ctx context.Context,
		value session.Session,
	) error

	Delete(
		ctx context.Context,
		idHash string,
	) error
}

type SessionCookieWriter interface {
	Write(
		response http.ResponseWriter,
		rawSessionID string,
	) error
}
