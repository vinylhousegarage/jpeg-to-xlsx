package callback

import "testing"

func TestNewHandlerRejectsNilDependency(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name  string
		build func(
			callbackTestDependencies,
		) (
			*Handler,
			error,
		)
	}{
		{
			name: "OAuth state store",
			build: func(
				dependencies callbackTestDependencies,
			) (
				*Handler,
				error,
			) {
				return NewHandler(
					nil,
					dependencies.cognitoClient,
					dependencies.verifier,
					dependencies.idGenerator,
					dependencies.sessionStore,
					dependencies.cookieWriter,
					validTestHandlerConfig(),
				)
			},
		},
		{
			name: "Cognito client",
			build: func(
				dependencies callbackTestDependencies,
			) (
				*Handler,
				error,
			) {
				return NewHandler(
					dependencies.oauthStateStore,
					nil,
					dependencies.verifier,
					dependencies.idGenerator,
					dependencies.sessionStore,
					dependencies.cookieWriter,
					validTestHandlerConfig(),
				)
			},
		},
		{
			name: "ID token verifier",
			build: func(
				dependencies callbackTestDependencies,
			) (
				*Handler,
				error,
			) {
				return NewHandler(
					dependencies.oauthStateStore,
					dependencies.cognitoClient,
					nil,
					dependencies.idGenerator,
					dependencies.sessionStore,
					dependencies.cookieWriter,
					validTestHandlerConfig(),
				)
			},
		},
		{
			name: "session ID generator",
			build: func(
				dependencies callbackTestDependencies,
			) (
				*Handler,
				error,
			) {
				return NewHandler(
					dependencies.oauthStateStore,
					dependencies.cognitoClient,
					dependencies.verifier,
					nil,
					dependencies.sessionStore,
					dependencies.cookieWriter,
					validTestHandlerConfig(),
				)
			},
		},
		{
			name: "session store",
			build: func(
				dependencies callbackTestDependencies,
			) (
				*Handler,
				error,
			) {
				return NewHandler(
					dependencies.oauthStateStore,
					dependencies.cognitoClient,
					dependencies.verifier,
					dependencies.idGenerator,
					nil,
					dependencies.cookieWriter,
					validTestHandlerConfig(),
				)
			},
		},
		{
			name: "cookie writer",
			build: func(
				dependencies callbackTestDependencies,
			) (
				*Handler,
				error,
			) {
				return NewHandler(
					dependencies.oauthStateStore,
					dependencies.cognitoClient,
					dependencies.verifier,
					dependencies.idGenerator,
					dependencies.sessionStore,
					nil,
					validTestHandlerConfig(),
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				dependencies :=
					newCallbackTestDependencies()

				handler, err := test.build(
					dependencies,
				)
				if err == nil {
					t.Fatal(
						"NewHandler() error = nil, want an error",
					)
				}

				if handler != nil {
					t.Errorf(
						"NewHandler() handler = %v, want nil",
						handler,
					)
				}
			},
		)
	}
}
