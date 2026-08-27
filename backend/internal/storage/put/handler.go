package put

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/storage"
)

// 構造体を定義
type Handler struct {
	service *Service
	logger  *zap.Logger
}

// 構造体を初期化
func NewHandler(service *Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// リクエストを検証
	req, err := h.validateRequest(r)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	// 署名付きURLの生成
	url, expiresAt, err := h.service.GeneratePresignURL(r.Context(), req.ShotNumber)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	// JSONで返却
	h.writeUploadResponse(w, url, expiresAt)
}

// 検証メソッド
func (h *Handler) validateRequest(r *http.Request) (*storage.PutPresignRequest, error) {
	if r.Method != http.MethodPost {
		return nil, apierror.New(apierror.ErrorCodeInvalidMethod, http.StatusMethodNotAllowed, nil)
	}

	var req storage.PutPresignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierror.New(apierror.ErrorCodeInvalidJSON, http.StatusBadRequest, err)
	}

	if req.ShotNumber == "" {
		return nil, apierror.New(apierror.ErrorCodeMissingShotNumber, http.StatusBadRequest, nil)
	}

	return &req, nil
}

// JSONレスポンス書き込みメソッド
func (h *Handler) writeUploadResponse(w http.ResponseWriter, url string, expiresAt time.Time) {
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
