package login

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

type Handler struct {
	clientID     string
	redirectURI  string
	cookieSecure bool
	logger       *zap.Logger
}

func NewHandler(
	clientID string,
	redirectURI string,
	cookieSecure bool,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		clientID:     clientID,
		redirectURI:  redirectURI,
		cookieSecure: cookieSecure,
		logger:       logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	state, err := oauth.GenerateState()
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	http.SetCookie(w, oauth.BuildStateCookie(state, h.cookieSecure))

	authURL := oauth.BuildAuthURL(h.clientID, h.redirectURI, state)
	http.Redirect(w, r, authURL, http.StatusFound)
}
