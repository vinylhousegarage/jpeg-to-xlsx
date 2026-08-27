package logger

import "testing"

func TestNewLogger(t *testing.T) {
	t.Parallel()

	t.Run("Production mode", func(t *testing.T) {
		t.Parallel()

		l, err := NewLogger("production")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if l == nil {
			t.Error("Expected logger instance, got nil")
		}
	})

	t.Run("Development mode", func(t *testing.T) {
		t.Parallel()

		l, err := NewLogger("development")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if l == nil {
			t.Error("Expected logger instance, got nil")
		}
	})
}
