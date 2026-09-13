package callback

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/apierror"
	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/slack/oauth"
)

const oauthStateCookieName = "oauth_state"

type codeExchanger interface {
	ExchangeCode(
		ctx context.Context,
		code string,
		redirectURI string,
	) (*oauth.Token, error)
}

type conversationOpener interface {
	OpenConversation(
		ctx context.Context,
		accessToken string,
		userID string,
	) (string, error)
}

type tokenStore interface {
	Save(
		ctx context.Context,
		cognitoSub string,
		token *oauth.Token,
	) error
}

type sessionResolver interface {
	Resolve(
		ctx context.Context,
		request *http.Request,
	) (
		authsession.Session,
		error,
	)
}

type Handler struct {
	redirectURI     string
	cookieSecure    bool
	sessionResolver sessionResolver
	exchanger       codeExchanger
	opener          conversationOpener
	store           tokenStore
	logger          *zap.Logger
}

func NewHandler(
	redirectURI string,
	cookieSecure bool,
	sessionResolver sessionResolver,
	exchanger codeExchanger,
	opener conversationOpener,
	store tokenStore,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		redirectURI:     redirectURI,
		cookieSecure:    cookieSecure,
		sessionResolver: sessionResolver,
		exchanger:       exchanger,
		opener:          opener,
		store:           store,
		logger:          logger,
	}
}

func (h *Handler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		apierror.WriteError(
			w,
			apierror.New(
				apierror.ErrorCodeInvalidMethod,
				http.StatusMethodNotAllowed,
				nil,
			),
			h.logger,
		)
		return
	}

	sessionValue, err := h.sessionResolver.Resolve(r.Context(), r)
	if err != nil {
		if errors.Is(err, authsession.ErrUnauthenticated) {
			apierror.WriteError(
				w,
				apierror.New(
					apierror.ErrorCodeUnauthorized,
					http.StatusUnauthorized,
					err,
				),
				h.logger,
			)
			return
		}

		apierror.WriteError(
			w,
			apierror.New(
				apierror.ErrorCodeInternal,
				http.StatusInternalServerError,
				err,
				"resolve session for Slack callback",
			),
			h.logger,
		)
		return
	}

	stateCookie, err := r.Cookie(oauthStateCookieName)
	if err != nil {
		apierror.WriteError(
			w,
			apierror.New(
				apierror.ErrorCodeMissingState,
				http.StatusBadRequest,
				err,
			),
			h.logger,
		)
		return
	}

	queryState := r.URL.Query().Get("state")
	if queryState == "" || queryState != stateCookie.Value {
		apierror.WriteError(
			w,
			apierror.New(
				apierror.ErrorCodeInvalidState,
				http.StatusBadRequest,
				nil,
			),
			h.logger,
		)
		return
	}

	http.SetCookie(w, oauth.BuildDeleteStateCookie(h.cookieSecure))

	code := r.URL.Query().Get("code")
	if code == "" {
		apierror.WriteError(
			w,
			apierror.New(
				apierror.ErrorCodeMissingCode,
				http.StatusBadRequest,
				nil,
			),
			h.logger,
		)
		return
	}

	token, err := h.exchanger.ExchangeCode(r.Context(), code, h.redirectURI)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	channelID, err := h.opener.OpenConversation(
		r.Context(),
		token.AccessToken,
		token.UserID,
	)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	token.ChannelID = channelID

	if err := h.store.Save(
		r.Context(),
		sessionValue.CognitoSub,
		token,
	); err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	http.Redirect(w, r, "/?slack=connected", http.StatusSeeOther)
}
