package xlsx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"
)

func Convert(jsonData []byte) ([]byte, error) {
	var values map[string]json.RawMessage
	if err := json.Unmarshal(jsonData, &values); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	file := excelize.NewFile()
	defer func() {
		_ = file.Close()
	}()

	sheet := file.GetSheetName(0)

	for index, key := range keys {
		row := index + 1

		value, err := valueToCell(values[key])
		if err != nil {
			return nil, fmt.Errorf("convert value for key %q: %w", key, err)
		}

		if err := file.SetCellValue(sheet, fmt.Sprintf("A%d", row), key); err != nil {
			return nil, fmt.Errorf("set key cell for %q: %w", key, err)
		}
		if err := file.SetCellValue(sheet, fmt.Sprintf("B%d", row), value); err != nil {
			return nil, fmt.Errorf("set value cell for %q: %w", key, err)
		}
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write xlsx: %w", err)
	}

	return buffer.Bytes(), nil
}

func valueToCell(raw json.RawMessage) (string, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return "", fmt.Errorf("decode JSON value: %w", err)
	}

	switch typedValue := value.(type) {
	case nil:
		return "", nil
	case string:
		return typedValue, nil
	case json.Number:
		return typedValue.String(), nil
	case bool:
		return fmt.Sprintf("%t", typedValue), nil
	default:
		jsonData, err := json.Marshal(typedValue)
		if err != nil {
			return "", fmt.Errorf("marshal structured JSON value: %w", err)
		}

		return string(jsonData), nil
	}
}
