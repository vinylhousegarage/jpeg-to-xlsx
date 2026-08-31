package router

import "net/http"

const (
	testAuthLoginResponse = "auth login"

	testAuthCallbackResponse = "auth callback"

	testSlackLoginResponse = "slack login"

	testSlackCallbackResponse = "slack callback"

	testStorageResponse = "storage"
)

func newTestRouter() *http.ServeMux {
	mux := http.NewServeMux()

	authLoginHandler := newTestHandler(
		testAuthLoginResponse,
	)

	authCallbackHandler := newTestHandler(
		testAuthCallbackResponse,
	)

	slackLoginHandler := newTestHandler(
		testSlackLoginResponse,
	)

	slackCallbackHandler := newTestHandler(
		testSlackCallbackResponse,
	)

	storageHandler := newTestHandler(
		testStorageResponse,
	)

	SetupRoutes(
		mux,
		authLoginHandler,
		authCallbackHandler,
		slackLoginHandler,
		slackCallbackHandler,
		storageHandler,
	)

	return mux
}

func newTestHandler(
	body string,
) http.Handler {
	return http.HandlerFunc(
		func(
			response http.ResponseWriter,
			_ *http.Request,
		) {
			_, _ = response.Write(
				[]byte(body),
			)
		},
	)
}
