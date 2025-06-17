package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestConvertToJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]interface{}
		expected []OrderedRow
	}{
		{
			name:     "empty input",
			input:    [][]interface{}{},
			expected: []OrderedRow{},
		},
		{
			name: "simple data",
			input: [][]interface{}{
				{"Name", "Age", "City"},
				{"Alice", 30, "Tokyo"},
				{"Bob", 25, "Osaka"},
			},
			expected: []OrderedRow{
				{
					Keys: []string{"Name", "Age", "City"},
					Values: map[string]interface{}{
						"Name": "Alice",
						"Age":  30,
						"City": "Tokyo",
					},
				},
				{
					Keys: []string{"Name", "Age", "City"},
					Values: map[string]interface{}{
						"Name": "Bob",
						"Age":  25,
						"City": "Osaka",
					},
				},
			},
		},
		{
			name: "missing cells",
			input: [][]interface{}{
				{"Name", "Age", "City"},
				{"Alice", 30},
				{"Bob"},
			},
			expected: []OrderedRow{
				{
					Keys: []string{"Name", "Age", "City"},
					Values: map[string]interface{}{
						"Name": "Alice",
						"Age":  30,
						"City": "",
					},
				},
				{
					Keys: []string{"Name", "Age", "City"},
					Values: map[string]interface{}{
						"Name": "Bob",
						"Age":  "",
						"City": "",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToJSON(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("convertToJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOrderedRowMarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		row          OrderedRow
		wantKeyOrder []string
	}{
		{
			name: "preserve header order",
			row: OrderedRow{
				Keys: []string{"Name", "Age", "City"},
				Values: map[string]interface{}{
					"Name": "Alice",
					"Age":  30,
					"City": "Tokyo",
				},
			},
			wantKeyOrder: []string{"Name", "Age", "City"},
		},
		{
			name: "different order than alphabetical",
			row: OrderedRow{
				Keys: []string{"Z", "A", "M"},
				Values: map[string]interface{}{
					"Z": "last",
					"A": "first",
					"M": "middle",
				},
			},
			wantKeyOrder: []string{"Z", "A", "M"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.row)
			if err != nil {
				t.Fatalf("Failed to marshal OrderedRow: %v", err)
			}

			// Check key order in JSON string
			jsonStr := string(jsonBytes)
			lastPos := -1
			for _, key := range tt.wantKeyOrder {
				pos := strings.Index(jsonStr, `"`+key+`":`)
				if pos == -1 {
					t.Errorf("Key %s not found in JSON", key)
					continue
				}
				if pos < lastPos {
					t.Errorf("Key %s appears before previous key (pos: %d < %d)", key, pos, lastPos)
				}
				lastPos = pos
			}
		})
	}
}
