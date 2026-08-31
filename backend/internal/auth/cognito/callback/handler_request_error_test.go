package callback

import (
	"net/http"
	"testing"
)

func TestHandlerRejectsInvalidCallbackQuery(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name       string
		code       string
		stateValue string
	}{
		{
			name:       "missing code",
			code:       "",
			stateValue: testOAuthStateValue,
		},
		{
			name:       "whitespace code",
			code:       "   ",
			stateValue: testOAuthStateValue,
		},
		{
			name:       "missing state",
			code:       testAuthorizationCode,
			stateValue: "",
		},
		{
			name:       "whitespace state",
			code:       testAuthorizationCode,
			stateValue: "   ",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				dependencies :=
					newCallbackTestDependencies()

				handler := newTestHandler(
					t,
					dependencies,
				)

				request :=
					newCallbackRequestWithQuery(
						test.code,
						test.stateValue,
					)
				response :=
					newCallbackResponseRecorder()

				handler.ServeHTTP(
					response,
					request,
				)

				if response.Code !=
					http.StatusBadRequest {
					t.Errorf(
						"response status = %d, want %d",
						response.Code,
						http.StatusBadRequest,
					)
				}

				if dependencies.oauthStateStore.
					consumeCalls != 0 {
					t.Errorf(
						"Consume() calls = %d, want 0",
						dependencies.oauthStateStore.
							consumeCalls,
					)
				}

				if dependencies.cognitoClient.
					exchangeCalls != 0 {
					t.Errorf(
						"Exchange() calls = %d, want 0",
						dependencies.cognitoClient.
							exchangeCalls,
					)
				}

				if dependencies.verifier.
					verifyCalls != 0 {
					t.Errorf(
						"Verify() calls = %d, want 0",
						dependencies.verifier.
							verifyCalls,
					)
				}

				if dependencies.idGenerator.
					generateCalls != 0 {
					t.Errorf(
						"Generate() calls = %d, want 0",
						dependencies.idGenerator.
							generateCalls,
					)
				}

				assertNoSessionWasEstablished(
					t,
					dependencies,
				)
			},
		)
	}
}
