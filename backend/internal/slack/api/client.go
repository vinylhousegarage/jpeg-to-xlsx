package api

import "net/http"

const (
	slackPostMessageURL      = "https://slack.com/api/chat.postMessage"
	slackOpenConversationURL = "https://slack.com/api/conversations.open"
)

type Client struct {
	httpClient          *http.Client
	postMessageURL      string
	openConversationURL string
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		httpClient:          httpClient,
		postMessageURL:      slackPostMessageURL,
		openConversationURL: slackOpenConversationURL,
	}
}
