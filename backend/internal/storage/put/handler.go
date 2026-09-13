package put

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/apierror"
	authsession "github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/auth/session"
	"github.com/vinylhousegarage/jpeg-to-xlsx/backend/internal/storage"
)

type sessionResolver interface {
	Resolve(
		ctx context.Context,
		request *http.Request,
	) (
		authsession.Session,
		error,
	)
}

type Handler struct {
	service         *Service
	sessionResolver sessionResolver
	logger          *zap.Logger
}

func NewHandler(
	service *Service,
	sessionResolver sessionResolver,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		service:         service,
		sessionResolver: sessionResolver,
		logger:          logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := h.validateRequest(r)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	sessionValue, err := h.sessionResolver.Resolve(r.Context(), r)
	if err != nil {
		if errors.Is(err, authsession.ErrUnauthenticated) {
			apierror.WriteError(
				w,
				apierror.New(
					apierror.ErrorCodeUnauthorized,
					http.StatusUnauthorized,
					err,
				),
				h.logger,
			)
			return
		}

		apierror.WriteError(
			w,
			apierror.New(
				apierror.ErrorCodeInternal,
				http.StatusInternalServerError,
				err,
				"resolve session for upload",
			),
			h.logger,
		)
		return
	}

	url, expiresAt, err := h.service.GeneratePresignURL(
		r.Context(),
		sessionValue.CognitoSub,
		req.ShotNumber,
	)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	h.writeUploadResponse(w, url, expiresAt)
}

func (h *Handler) validateRequest(
	r *http.Request,
) (
	*storage.PutPresignRequest,
	error,
) {
	if r.Method != http.MethodPost {
		return nil, apierror.New(
			apierror.ErrorCodeInvalidMethod,
			http.StatusMethodNotAllowed,
			nil,
		)
	}

	var req storage.PutPresignRequest
	if err := json.NewDecoder(
		r.Body,
	).Decode(&req); err != nil {
		return nil, apierror.New(
			apierror.ErrorCodeInvalidJSON,
			http.StatusBadRequest,
			err,
		)
	}

	if req.ShotNumber == "" {
		return nil, apierror.New(
			apierror.ErrorCodeMissingShotNumber,
			http.StatusBadRequest,
			nil,
		)
	}

	return &req, nil
}

func (h *Handler) writeUploadResponse(
	w http.ResponseWriter,
	url string,
	expiresAt time.Time,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := storage.PutPresignResponse{
		ExpiresAt: expiresAt,
		UploadURL: url,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}
}
