package login

import (
	"testing"
)

func TestNewHandler(
	t *testing.T,
) {
	t.Parallel()

	dependencies :=
		newTestHandlerDependencies()

	handler, err := NewHandler(
		dependencies.stateGenerator,
		dependencies.stateStore,
		dependencies.authorizationClient,
	)
	if err != nil {
		t.Fatalf(
			"NewHandler() error = %v",
			err,
		)
	}

	if handler == nil {
		t.Fatal(
			"NewHandler() handler = nil, want non-nil",
		)
	}
}

func TestNewHandlerRejectsNilDependency(
	t *testing.T,
) {
	t.Parallel()

	tests := []struct {
		name       string
		newHandler func(
			dependencies testHandlerDependencies,
		) (*Handler, error)
	}{
		{
			name: "nil state generator",
			newHandler: func(
				dependencies testHandlerDependencies,
			) (*Handler, error) {
				return NewHandler(
					nil,
					dependencies.stateStore,
					dependencies.
						authorizationClient,
				)
			},
		},
		{
			name: "nil state store",
			newHandler: func(
				dependencies testHandlerDependencies,
			) (*Handler, error) {
				return NewHandler(
					dependencies.stateGenerator,
					nil,
					dependencies.
						authorizationClient,
				)
			},
		},
		{
			name: "nil authorization client",
			newHandler: func(
				dependencies testHandlerDependencies,
			) (*Handler, error) {
				return NewHandler(
					dependencies.stateGenerator,
					dependencies.stateStore,
					nil,
				)
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				dependencies :=
					newTestHandlerDependencies()

				handler, err :=
					test.newHandler(
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
