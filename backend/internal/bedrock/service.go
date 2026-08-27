package bedrock

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// インターフェースを定義
type ImageAnalyzer interface {
	Analyze(ctx context.Context, imgData []byte) (string, error)
}

// 構造体を定義
type Service struct {
	analyzer ImageAnalyzer
}

// 構造体を初期化
func NewService(analyzer ImageAnalyzer) *Service {
	return &Service{analyzer: analyzer}
}

// メソッドを定義
func (s *Service) ProcessImage(ctx context.Context, rawImage []byte) (map[string]interface{}, error) {
	// タイムアウトを付与
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// タイムアウト付きのコンテキストを client に渡す
	text, err := s.analyzer.Analyze(timeoutCtx, rawImage)
	if err != nil {
		return nil, fmt.Errorf("analyze error: %w", err)
	}

	// パース
	return s.parseResponse(text)
}

// テキストを受け取り、JSONを抽出してパース
func (s *Service) parseResponse(text string) (map[string]any, error) {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")

	if start == -1 || end == -1 || start >= end {
		return nil, fmt.Errorf("invalid json format: %s", text)
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(text[start:end+1]), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}
