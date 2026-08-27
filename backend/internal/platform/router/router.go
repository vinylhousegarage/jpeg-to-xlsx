package router

import "net/http"

func SetupRoutes(
	mux *http.ServeMux,
	presignHandler http.Handler,
	slackLoginHandler http.Handler,
	slackCallbackHandler http.Handler,
) {
	mux.Handle("/api/storage/upload", presignHandler)
	mux.Handle("/api/oauth/slack/login", slackLoginHandler)
	mux.Handle("/api/oauth/slack/callback", slackCallbackHandler)
}
