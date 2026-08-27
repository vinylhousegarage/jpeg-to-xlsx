package bedrock

import (
	"testing"
)

func TestParseResponse(t *testing.T) {
	s := NewService(nil)

	tests := []struct {
		name    string
		input   string
		wantKey string
		wantErr bool
	}{
		{
			name:    "Valid JSON",
			input:   `{"key": "value"}`,
			wantKey: "value",
			wantErr: false,
		},
		{
			name:    "JSON with Markdown",
			input:   "Here is your json:\n```json\n{\"key\": \"value\"}\n```",
			wantKey: "value",
			wantErr: false,
		},
		{
			name:    "Invalid JSON Format",
			input:   "This is just text, not JSON",
			wantKey: "",
			wantErr: true,
		},
		{
			name:    "Empty JSON",
			input:   "{}",
			wantKey: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := s.parseResponse(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantKey != "" {
				if val, ok := res["key"].(string); !ok || val != tt.wantKey {
					t.Errorf("got %v, want %v", res["key"], tt.wantKey)
				}
			}
		})
	}
}
