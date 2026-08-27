package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const slackTokenURL = "https://slack.com/api/oauth.v2.access"

type Token struct {
	AccessToken string
	BotUserID   string
	TeamID      string
	UserID      string
	ChannelID   string
}

type Client struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
	tokenURL     string
}

func NewClient(
	httpClient *http.Client,
	clientID string,
	clientSecret string,
) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		httpClient:   httpClient,
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenURL:     slackTokenURL,
	}
}

type tokenResponse struct {
	OK          bool   `json:"ok"`
	AccessToken string `json:"access_token"`
	BotUserID   string `json:"bot_user_id"`
	Error       string `json:"error"`

	Team struct {
		ID string `json:"id"`
	} `json:"team"`

	AuthedUser struct {
		ID string `json:"id"`
	} `json:"authed_user"`
}

func (c *Client) ExchangeCode(
	ctx context.Context,
	code string,
	redirectURI string,
) (*Token, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.tokenURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create Slack OAuth token request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	req.SetBasicAuth(
		c.clientID,
		c.clientSecret,
	)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"send Slack OAuth token request: %w",
			err,
		)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"slack OAuth token request returned status %d",
			res.StatusCode,
		)
	}

	var payload tokenResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf(
			"decode Slack OAuth token response: %w",
			err,
		)
	}

	if !payload.OK {
		return nil, fmt.Errorf(
			"slack OAuth token exchange failed: %s",
			payload.Error,
		)
	}

	if payload.AccessToken == "" {
		return nil, fmt.Errorf(
			"slack OAuth token response is missing access_token",
		)
	}

	if payload.AuthedUser.ID == "" {
		return nil, fmt.Errorf(
			"slack OAuth token response is missing authed_user.id",
		)
	}

	return &Token{
		AccessToken: payload.AccessToken,
		BotUserID:   payload.BotUserID,
		TeamID:      payload.Team.ID,
		UserID:      payload.AuthedUser.ID,
	}, nil
}
