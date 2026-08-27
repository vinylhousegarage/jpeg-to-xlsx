package apierror

import "fmt"

type ErrorCode string

const (
	// Common
	ErrorCodeInternal      ErrorCode = "internal_server_error"
	ErrorCodeInvalidJSON   ErrorCode = "invalid_json"
	ErrorCodeInvalidMethod ErrorCode = "invalid_method"

	// Upload
	ErrorCodeMissingShotNumber ErrorCode = "missing_shot_number"
	ErrorCodeS3SigningFailed   ErrorCode = "s3_signing_failed"

	// Slack OAuth
	ErrorCodeMissingState ErrorCode = "missing_state"
	ErrorCodeInvalidState ErrorCode = "invalid_state"
	ErrorCodeMissingCode  ErrorCode = "missing_code"
)

type APIError struct {
	Code       ErrorCode
	HTTPStatus int
	Err        error
	Internal   string
}

func (e *APIError) Error() string {
	if e.Err == nil {
		return string(e.Code)
	}

	return fmt.Sprintf("[%s] %v", e.Code, e.Err)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, status int, err error, internal ...string) *APIError {
	apiErr := &APIError{
		Code:       code,
		HTTPStatus: status,
		Err:        err,
	}
	if len(internal) > 0 {
		apiErr.Internal = internal[0]
	}
	return apiErr
}
