package session

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const sessionCookieName = "__Host-session"

type CookieManager struct {
	lifetime time.Duration
	now      func() time.Time
}

func NewCookieManager(
	lifetime time.Duration,
) (*CookieManager, error) {
	if lifetime < time.Second {
		return nil, fmt.Errorf(
			"create cookie manager: lifetime must be at least one second",
		)
	}

	return &CookieManager{
		lifetime: lifetime,
		now:      time.Now,
	}, nil
}

func (
	manager *CookieManager,
) Write(
	response http.ResponseWriter,
	rawSessionID string,
) error {
	if manager == nil {
		return fmt.Errorf(
			"write session cookie: cookie manager is nil",
		)
	}

	if response == nil {
		return fmt.Errorf(
			"write session cookie: response writer is nil",
		)
	}

	if manager.now == nil {
		return fmt.Errorf(
			"write session cookie: clock is nil",
		)
	}

	if strings.TrimSpace(rawSessionID) == "" {
		return fmt.Errorf(
			"write session cookie: session ID is empty",
		)
	}

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    rawSessionID,
		Path:     "/",
		Expires:  manager.now().Add(manager.lifetime),
		MaxAge:   int(manager.lifetime / time.Second),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	if err := cookie.Valid(); err != nil {
		return fmt.Errorf(
			"write session cookie: validate cookie: %w",
			err,
		)
	}

	http.SetCookie(
		response,
		cookie,
	)

	return nil
}

func (
	manager *CookieManager,
) Read(
	request *http.Request,
) (string, error) {
	if manager == nil {
		return "", fmt.Errorf(
			"read session cookie: cookie manager is nil",
		)
	}

	if request == nil {
		return "", fmt.Errorf(
			"read session cookie: request is nil",
		)
	}

	cookie, err := request.Cookie(
		sessionCookieName,
	)
	if err != nil {
		return "", fmt.Errorf(
			"read session cookie: %w",
			err,
		)
	}

	if strings.TrimSpace(cookie.Value) == "" {
		return "", fmt.Errorf(
			"read session cookie: session ID is empty",
		)
	}

	return cookie.Value, nil
}

func (
	manager *CookieManager,
) Delete(
	response http.ResponseWriter,
) {
	if manager == nil || response == nil {
		return
	}

	http.SetCookie(
		response,
		&http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(1, 0).UTC(),
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)
}
