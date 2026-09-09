package cognito

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/oauth2"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/cognito/oauthstate"
)

const (
	minCodeVerifierLength = 43
	maxCodeVerifierLength = 128
)

type Config struct {
	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	RedirectURI           string
	Scopes                []string
}

type Client struct {
	oauthConfig oauth2.Config
}

func NewClient(
	config Config,
) (*Client, error) {
	if strings.TrimSpace(config.ClientID) == "" {
		return nil, fmt.Errorf(
			"create Cognito client: client ID is empty",
		)
	}

	if strings.TrimSpace(config.ClientSecret) == "" {
		return nil, fmt.Errorf(
			"create Cognito client: client secret is empty",
		)
	}

	if err := validateEndpoint(
		"authorization endpoint",
		config.AuthorizationEndpoint,
	); err != nil {
		return nil, fmt.Errorf(
			"create Cognito client: %w",
			err,
		)
	}

	if err := validateEndpoint(
		"token endpoint",
		config.TokenEndpoint,
	); err != nil {
		return nil, fmt.Errorf(
			"create Cognito client: %w",
			err,
		)
	}

	if err := validateEndpoint(
		"redirect URI",
		config.RedirectURI,
	); err != nil {
		return nil, fmt.Errorf(
			"create Cognito client: %w",
			err,
		)
	}

	if len(config.Scopes) == 0 {
		return nil, fmt.Errorf(
			"create Cognito client: scopes are empty",
		)
	}

	scopes := make(
		[]string,
		len(config.Scopes),
	)

	for index, scope := range config.Scopes {
		if strings.TrimSpace(scope) == "" {
			return nil, fmt.Errorf(
				"create Cognito client: scope at index %d is empty",
				index,
			)
		}

		scopes[index] = scope
	}

	return &Client{
		oauthConfig: oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL: config.
					AuthorizationEndpoint,
				TokenURL: config.TokenEndpoint,
				AuthStyle: oauth2.
					AuthStyleInHeader,
			},
			RedirectURL: config.RedirectURI,
			Scopes:      scopes,
		},
	}, nil
}

func (c *Client) AuthorizationURL(
	state oauthstate.State,
) (string, error) {
	if c == nil {
		return "", fmt.Errorf(
			"create authorization URL: client is nil",
		)
	}

	if strings.TrimSpace(state.Value) == "" {
		return "", fmt.Errorf(
			"create authorization URL: state is empty",
		)
	}

	if err := validateCodeVerifier(
		state.CodeVerifier,
	); err != nil {
		return "", fmt.Errorf(
			"create authorization URL: %w",
			err,
		)
	}

	if strings.TrimSpace(state.Nonce) == "" {
		return "", fmt.Errorf(
			"create authorization URL: nonce is empty",
		)
	}

	return c.oauthConfig.AuthCodeURL(
		state.Value,
		oauth2.S256ChallengeOption(
			state.CodeVerifier,
		),
		oauth2.SetAuthURLParam(
			"nonce",
			state.Nonce,
		),
		oauth2.SetAuthURLParam(
			"identity_provider",
			"Google",
		),
	), nil
}

func (c *Client) Exchange(
	ctx context.Context,
	code string,
	codeVerifier string,
) (string, error) {
	if c == nil {
		return "", fmt.Errorf(
			"exchange authorization code: client is nil",
		)
	}

	if ctx == nil {
		return "", fmt.Errorf(
			"exchange authorization code: context is nil",
		)
	}

	if strings.TrimSpace(code) == "" {
		return "", fmt.Errorf(
			"exchange authorization code: code is empty",
		)
	}

	if err := validateCodeVerifier(
		codeVerifier,
	); err != nil {
		return "", fmt.Errorf(
			"exchange authorization code: %w",
			err,
		)
	}

	token, err := c.oauthConfig.Exchange(
		ctx,
		code,
		oauth2.VerifierOption(
			codeVerifier,
		),
	)
	if err != nil {
		return "", fmt.Errorf(
			"exchange authorization code: token endpoint: %w",
			err,
		)
	}

	rawIDToken := token.Extra("id_token")

	idToken, ok := rawIDToken.(string)
	if !ok {
		return "", fmt.Errorf(
			"exchange authorization code: ID token is missing or invalid",
		)
	}

	if strings.TrimSpace(idToken) == "" {
		return "", fmt.Errorf(
			"exchange authorization code: ID token is empty",
		)
	}

	return idToken, nil
}

func validateEndpoint(
	name string,
	value string,
) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf(
			"%s is empty",
			name,
		)
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf(
			"%s is invalid: %w",
			name,
			err,
		)
	}

	if parsedURL.Scheme != "http" &&
		parsedURL.Scheme != "https" {
		return fmt.Errorf(
			"%s has unsupported scheme %q",
			name,
			parsedURL.Scheme,
		)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf(
			"%s host is empty",
			name,
		)
	}

	return nil
}

func validateCodeVerifier(
	value string,
) error {
	length := len(value)

	if length < minCodeVerifierLength ||
		length > maxCodeVerifierLength {
		return fmt.Errorf(
			"code verifier length is %d, want between %d and %d",
			length,
			minCodeVerifierLength,
			maxCodeVerifierLength,
		)
	}

	for _, character := range value {
		if !isCodeVerifierCharacter(
			character,
		) {
			return fmt.Errorf(
				"code verifier contains invalid character %q",
				character,
			)
		}
	}

	return nil
}

func isCodeVerifierCharacter(
	character rune,
) bool {
	switch {
	case character >= 'A' &&
		character <= 'Z':
		return true

	case character >= 'a' &&
		character <= 'z':
		return true

	case character >= '0' &&
		character <= '9':
		return true

	case character == '-',
		character == '.',
		character == '_',
		character == '~':
		return true

	default:
		return false
	}
}
