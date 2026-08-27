package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type openConversationRequest struct {
	Users string `json:"users"`
}

type openConversationResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`

	Channel struct {
		ID string `json:"id"`
	} `json:"channel"`
}

func (c *Client) OpenConversation(
	ctx context.Context,
	accessToken string,
	userID string,
) (string, error) {
	if accessToken == "" {
		return "", fmt.Errorf(
			"open slack conversation: access token is empty",
		)
	}

	if userID == "" {
		return "", fmt.Errorf(
			"open slack conversation: user ID is empty",
		)
	}

	body, err := json.Marshal(
		openConversationRequest{
			Users: userID,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"marshal slack open conversation request: %w",
			err,
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.openConversationURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf(
			"create slack open conversation request: %w",
			err,
		)
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+accessToken,
	)
	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf(
			"send slack open conversation request: %w",
			err,
		)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"slack open conversation request returned status %d",
			res.StatusCode,
		)
	}

	var payload openConversationResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf(
			"decode slack open conversation response: %w",
			err,
		)
	}

	if !payload.OK {
		return "", fmt.Errorf(
			"slack open conversation failed: %s",
			payload.Error,
		)
	}

	if payload.Channel.ID == "" {
		return "", fmt.Errorf(
			"slack open conversation response is missing channel.id",
		)
	}

	return payload.Channel.ID, nil
}
