package callback

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func validateRedirectURL(
	value string,
) (
	string,
	error,
) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf(
			"redirect URL is empty",
		)
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		return "", fmt.Errorf(
			"redirect URL is invalid: %w",
			err,
		)
	}

	if parsedURL.Scheme != "https" {
		return "", fmt.Errorf(
			"redirect URL scheme = %q, want %q",
			parsedURL.Scheme,
			"https",
		)
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf(
			"redirect URL host is empty",
		)
	}

	if parsedURL.User != nil {
		return "", fmt.Errorf(
			"redirect URL must not contain user information",
		)
	}

	return parsedURL.String(), nil
}

func writeRedirect(
	response http.ResponseWriter,
	redirectURL string,
) {
	response.Header().Set(
		"Location",
		redirectURL,
	)

	response.WriteHeader(
		http.StatusSeeOther,
	)
}

func writeError(
	response http.ResponseWriter,
	status int,
) {
	http.Error(
		response,
		http.StatusText(status),
		status,
	)
}
