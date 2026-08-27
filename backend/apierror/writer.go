package apierror

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// ErrorResponse 構造体を定義
type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, err error, logger *zap.Logger) {
	var apiErr *APIError

	// デフォルトを status 500 に設定
	status := http.StatusInternalServerError
	code := ErrorCodeInternal

	if errors.As(err, &apiErr) {
		status = apiErr.HTTPStatus
		code = apiErr.Code
	}

	fields := []zap.Field{
		zap.String("error_code", string(code)),
		zap.Int("http_status", status),
		zap.Error(err),
	}

	if status >= http.StatusInternalServerError {
		logger.Error("request failed", fields...)
	} else {
		logger.Warn("request rejected", fields...)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encodeErr := json.NewEncoder(w).Encode(ErrorResponse{
		Error: string(code),
	}); encodeErr != nil {
		logger.Error(
			"failed to write json response",
			zap.Error(encodeErr),
		)
	}
}
