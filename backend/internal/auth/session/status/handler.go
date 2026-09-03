package status

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
)

type SessionResolver interface {
	Resolve(
		ctx context.Context,
		request *http.Request,
	) (
		authsession.Session,
		error,
	)
}

type Handler struct {
	resolver SessionResolver
}

type response struct {
	Authenticated bool `json:"authenticated"`
}

func NewHandler(
	resolver SessionResolver,
) (
	*Handler,
	error,
) {
	if resolver == nil {
		return nil, fmt.Errorf(
			"create session status handler: resolver is nil",
		)
	}

	return &Handler{
		resolver: resolver,
	}, nil
}

func (
	handler *Handler,
) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if writer == nil {
		return
	}

	writer.Header().Set(
		"Cache-Control",
		"no-store",
	)

	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	if handler == nil ||
		handler.resolver == nil ||
		request == nil {
		writeResponse(
			writer,
			http.StatusInternalServerError,
			false,
		)

		return
	}

	_, err := handler.resolver.Resolve(
		request.Context(),
		request,
	)
	if err != nil {
		if errors.Is(
			err,
			authsession.ErrUnauthenticated,
		) {
			writeResponse(
				writer,
				http.StatusUnauthorized,
				false,
			)

			return
		}

		writeResponse(
			writer,
			http.StatusInternalServerError,
			false,
		)

		return
	}

	writeResponse(
		writer,
		http.StatusOK,
		true,
	)
}

func writeResponse(
	writer http.ResponseWriter,
	statusCode int,
	authenticated bool,
) {
	writer.WriteHeader(
		statusCode,
	)

	_ = json.NewEncoder(
		writer,
	).Encode(
		response{
			Authenticated: authenticated,
		},
	)
}
