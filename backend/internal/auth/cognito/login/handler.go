package login

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
)

type StateGenerator interface {
	Generate() (
		oauthstate.State,
		error,
	)
}

type StateStore interface {
	Save(
		ctx context.Context,
		state oauthstate.State,
	) error
}

type AuthorizationClient interface {
	AuthorizationURL(
		state oauthstate.State,
	) (string, error)
}

type Handler struct {
	stateGenerator      StateGenerator
	stateStore          StateStore
	authorizationClient AuthorizationClient
}

func NewHandler(
	stateGenerator StateGenerator,
	stateStore StateStore,
	authorizationClient AuthorizationClient,
) (*Handler, error) {
	if stateGenerator == nil {
		return nil, fmt.Errorf(
			"create Cognito login handler: state generator is nil",
		)
	}

	if stateStore == nil {
		return nil, fmt.Errorf(
			"create Cognito login handler: state store is nil",
		)
	}

	if authorizationClient == nil {
		return nil, fmt.Errorf(
			"create Cognito login handler: authorization client is nil",
		)
	}

	return &Handler{
		stateGenerator:      stateGenerator,
		stateStore:          stateStore,
		authorizationClient: authorizationClient,
	}, nil
}

func (
	handler *Handler,
) ServeHTTP(
	response http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		response.Header().Set(
			"Allow",
			http.MethodGet,
		)

		http.Error(
			response,
			http.StatusText(
				http.StatusMethodNotAllowed,
			),
			http.StatusMethodNotAllowed,
		)

		return
	}

	response.Header().Set(
		"Cache-Control",
		"no-store",
	)
	response.Header().Set(
		"Pragma",
		"no-cache",
	)

	state, err :=
		handler.stateGenerator.Generate()
	if err != nil {
		writeInternalServerError(response)

		return
	}

	authorizationURL, err :=
		handler.authorizationClient.
			AuthorizationURL(state)
	if err != nil {
		writeInternalServerError(response)

		return
	}

	if err := handler.stateStore.Save(
		request.Context(),
		state,
	); err != nil {
		writeInternalServerError(response)

		return
	}

	http.Redirect(
		response,
		request,
		authorizationURL,
		http.StatusFound,
	)
}

func writeInternalServerError(
	response http.ResponseWriter,
) {
	http.Error(
		response,
		http.StatusText(
			http.StatusInternalServerError,
		),
		http.StatusInternalServerError,
	)
}
