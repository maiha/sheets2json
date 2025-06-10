package main

import (
	"reflect"
	"testing"
)

func TestConvertToJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]interface{}
		expected []map[string]interface{}
	}{
		{
			name:     "empty input",
			input:    [][]interface{}{},
			expected: []map[string]interface{}{},
		},
		{
			name: "simple data",
			input: [][]interface{}{
				{"Name", "Age", "City"},
				{"Alice", 30, "Tokyo"},
				{"Bob", 25, "Osaka"},
			},
			expected: []map[string]interface{}{
				{"Name": "Alice", "Age": 30, "City": "Tokyo"},
				{"Name": "Bob", "Age": 25, "City": "Osaka"},
			},
		},
		{
			name: "missing cells",
			input: [][]interface{}{
				{"Name", "Age", "City"},
				{"Alice", 30},
				{"Bob"},
			},
			expected: []map[string]interface{}{
				{"Name": "Alice", "Age": 30, "City": ""},
				{"Name": "Bob", "Age": "", "City": ""},
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
