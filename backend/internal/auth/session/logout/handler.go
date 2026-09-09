package logout

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

type SessionStore interface {
	Delete(
		ctx context.Context,
		sessionIDHash string,
	) error
}

type CookieDeleter interface {
	Delete(
		writer http.ResponseWriter,
	)
}

type LogoutURLProvider interface {
	LogoutURL() (
		string,
		error,
	)
}

type Handler struct {
	resolver          SessionResolver
	store             SessionStore
	cookieDeleter     CookieDeleter
	logoutURLProvider LogoutURLProvider
}

type logoutResponse struct {
	LogoutURL string `json:"logout_url"`
}

func NewHandler(
	resolver SessionResolver,
	store SessionStore,
	cookieDeleter CookieDeleter,
	logoutURLProvider LogoutURLProvider,
) (
	*Handler,
	error,
) {
	if resolver == nil {
		return nil, fmt.Errorf("create logout handler: resolver is nil")
	}

	if store == nil {
		return nil, fmt.Errorf("create logout handler: session store is nil")
	}

	if cookieDeleter == nil {
		return nil, fmt.Errorf("create logout handler: cookie deleter is nil")
	}

	if logoutURLProvider == nil {
		return nil, fmt.Errorf("create logout handler: logout URL provider is nil")
	}

	return &Handler{
		resolver:          resolver,
		store:             store,
		cookieDeleter:     cookieDeleter,
		logoutURLProvider: logoutURLProvider,
	}, nil
}

func (handler *Handler) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if writer == nil {
		return
	}

	writer.Header().Set("Cache-Control", "no-store")

	if handler == nil ||
		handler.resolver == nil ||
		handler.store == nil ||
		handler.cookieDeleter == nil ||
		handler.logoutURLProvider == nil ||
		request == nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	logoutURL, err := handler.logoutURLProvider.LogoutURL()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	sessionValue, err := handler.resolver.Resolve(
		request.Context(),
		request,
	)
	if err != nil {
		if errors.Is(err, authsession.ErrUnauthenticated) {
			handler.cookieDeleter.Delete(writer)

			writeLogoutResponse(writer, logoutURL)

			return
		}

		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if err := handler.store.Delete(request.Context(), sessionValue.IDHash); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	handler.cookieDeleter.Delete(writer)

	writeLogoutResponse(writer, logoutURL)
}

func writeLogoutResponse(writer http.ResponseWriter, logoutURL string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(writer).Encode(
		logoutResponse{
			LogoutURL: logoutURL,
		},
	)
}
