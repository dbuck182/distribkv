package resp

import (
	"bufio"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Value
	}{
		{
			name:     "Simple String",
			input:    "+OK\r\n",
			expected: Value{Typ: '+', Str: "OK"},
		},
		{
			name:     "Error",
			input:    "-Error message\r\n",
			expected: Value{Typ: '-', Str: "Error message"},
		},
		{
			name:     "Integer",
			input:    ":1000\r\n",
			expected: Value{Typ: ':', Num: 1000},
		},
		{
			name:     "Bulk String",
			input:    "$5\r\nhello\r\n",
			expected: Value{Typ: '$', Bulk: "hello"},
		},
		{
			name:     "Bulk String with newline inside",
			input:    "$11\r\nhello\nworld\r\n",
			expected: Value{Typ: '$', Bulk: "hello\nworld"},
		},
		{
			name:     "Null Bulk String",
			input:    "$-1\r\n",
			expected: Value{Typ: '$'}, // zero value Bulk is ""
		},
		{
			name:  "Array of Bulk Strings (SET command)",
			input: "*3\r\n$3\r\nSET\r\n$1\r\na\r\n$1\r\n1\r\n",
			expected: Value{
				Typ: '*',
				Array: []Value{
					{Typ: '$', Bulk: "SET"},
					{Typ: '$', Bulk: "a"},
					{Typ: '$', Bulk: "1"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(bufio.NewReader(strings.NewReader(tt.input)))
			val, err := r.Read()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if val.Typ != tt.expected.Typ {
				t.Errorf("expected type %c, got %c", tt.expected.Typ, val.Typ)
			}
			if val.Str != tt.expected.Str {
				t.Errorf("expected Str %q, got %q", tt.expected.Str, val.Str)
			}
			if val.Bulk != tt.expected.Bulk {
				t.Errorf("expected Bulk %q, got %q", tt.expected.Bulk, val.Bulk)
			}
			if len(val.Array) != len(tt.expected.Array) {
				t.Fatalf("expected array length %d, got %d", len(tt.expected.Array), len(val.Array))
			}
			for i := range val.Array {
				if val.Array[i].Bulk != tt.expected.Array[i].Bulk {
					t.Errorf("array element %d: expected Bulk %q, got %q", i, tt.expected.Array[i].Bulk, val.Array[i].Bulk)
				}
			}
		})
	}
}
