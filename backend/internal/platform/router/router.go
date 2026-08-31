package router

import "net/http"

func SetupRoutes(
	mux *http.ServeMux,
	authLoginHandler http.Handler,
	authCallbackHandler http.Handler,
	slackLoginHandler http.Handler,
	slackCallbackHandler http.Handler,
	storageHandler http.Handler,
) {
	mux.Handle("HEAD /api/auth/login", methodNotAllowed(http.MethodGet))
	mux.Handle("HEAD /api/auth/callback", methodNotAllowed(http.MethodGet))
	mux.Handle("HEAD /api/oauth/slack/login", methodNotAllowed(http.MethodGet))
	mux.Handle("HEAD /api/oauth/slack/callback", methodNotAllowed(http.MethodGet))
	mux.Handle("HEAD /api/storage/upload", methodNotAllowed(http.MethodPost))

	mux.Handle("GET /api/auth/login", authLoginHandler)
	mux.Handle("GET /api/auth/callback", authCallbackHandler)
	mux.Handle("GET /api/oauth/slack/login", slackLoginHandler)
	mux.Handle("GET /api/oauth/slack/callback", slackCallbackHandler)
	mux.Handle("POST /api/storage/upload", storageHandler)
}

func methodNotAllowed(allowedMethod string) http.Handler {
	return http.HandlerFunc(
		func(response http.ResponseWriter, _ *http.Request) {
			response.Header().Set("Allow", allowedMethod)
			response.WriteHeader(http.StatusMethodNotAllowed)
		},
	)
}
