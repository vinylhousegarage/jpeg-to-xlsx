package xlsx

import (
	"bytes"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{
		"name": "山田太郎",
		"age": 30,
		"member": true,
		"note": null,
		"tags": ["AWS", "Go"]
	}`)

	xlsxData, err := Convert(jsonData)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	file, err := excelize.OpenReader(bytes.NewReader(xlsxData))
	if err != nil {
		t.Fatalf("open generated xlsx: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	sheet := file.GetSheetName(0)

	tests := []struct {
		cell string
		want string
	}{
		{cell: "A1", want: "key"},
		{cell: "B1", want: "value"},
		{cell: "A2", want: "age"},
		{cell: "B2", want: "30"},
		{cell: "A3", want: "member"},
		{cell: "B3", want: "true"},
		{cell: "A4", want: "name"},
		{cell: "B4", want: "山田太郎"},
		{cell: "A5", want: "note"},
		{cell: "B5", want: ""},
		{cell: "A6", want: "tags"},
		{cell: "B6", want: `["AWS","Go"]`},
	}

	for _, tt := range tests {
		got, err := file.GetCellValue(sheet, tt.cell)
		if err != nil {
			t.Fatalf("get cell %s: %v", tt.cell, err)
		}
		if got != tt.want {
			t.Errorf("cell %s = %q, want %q", tt.cell, got, tt.want)
		}
	}
}

func TestConvert_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := Convert([]byte(`{"name":`))
	if err == nil {
		t.Fatal("Convert() error = nil, want error")
	}
}
