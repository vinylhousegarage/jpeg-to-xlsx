package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	OAuthStateBytes = 32
	OAuthStateTTL   = 10 * time.Minute

	oauthStateCookieName = "oauth_state"
	oauthStateCookiePath = "/api/oauth/slack"

	slackAuthorizeURL = "https://slack.com/oauth/v2/authorize"
	slackBotScopes    = "chat:write,im:write"
)

func GenerateState() (string, error) {
	stateBytes := make([]byte, OAuthStateBytes)

	if _, err := rand.Read(stateBytes); err != nil {
		return "", fmt.Errorf("generate OAuth state: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(stateBytes), nil
}

func BuildStateCookie(state string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		Path:     oauthStateCookiePath,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(OAuthStateTTL.Seconds()),
		Expires:  time.Now().Add(OAuthStateTTL),
	}
}

func BuildDeleteStateCookie(secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    "",
		Path:     oauthStateCookiePath,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
}

func BuildAuthURL(clientID, redirectURI, state string) string {
	query := url.Values{}
	query.Set("client_id", clientID)
	query.Set("scope", slackBotScopes)
	query.Set("redirect_uri", redirectURI)
	query.Set("state", state)

	return slackAuthorizeURL + "?" + query.Encode()
}
