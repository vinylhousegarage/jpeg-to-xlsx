package apierror

import (
	"errors"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("base error")
	internalInfo := "debug info"

	tests := []struct {
		name         string
		code         ErrorCode
		status       int
		err          error
		internalArgs []string
		wantInternal string
	}{
		{
			name:         "all arguments provided",
			code:         ErrorCode("TEST_CODE"),
			status:       http.StatusInternalServerError,
			err:          originalErr,
			internalArgs: []string{internalInfo},
			wantInternal: internalInfo,
		},
		{
			name:         "no internal info",
			code:         ErrorCode("TEST_CODE_NO_INTERNAL"),
			status:       http.StatusBadRequest,
			err:          originalErr,
			internalArgs: nil,
			wantInternal: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := New(
				tt.code,
				tt.status,
				tt.err,
				tt.internalArgs...,
			)

			if got.Code != tt.code {
				t.Errorf("Code = %q, want %q", got.Code, tt.code)
			}

			if got.HTTPStatus != tt.status {
				t.Errorf(
					"HTTPStatus = %d, want %d",
					got.HTTPStatus,
					tt.status,
				)
			}

			if !errors.Is(got.Err, tt.err) {
				t.Errorf("Err = %v, want %v", got.Err, tt.err)
			}

			if got.Internal != tt.wantInternal {
				t.Errorf(
					"Internal = %q, want %q",
					got.Internal,
					tt.wantInternal,
				)
			}
		})
	}
}
