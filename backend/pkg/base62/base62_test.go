package base62

import (
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
		wantErr  bool
	}{
		{"zero", 0, "0", false},
		{"small number", 10, "a", false},
		{"medium number", 62, "10", false},
		{"large number", 1000, "g8", false},
		{"very large", 1000000, "4c92", false},
		{"negative", -1, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("Encode(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"zero", "0", 0, false},
		{"small", "a", 10, false},
		{"medium", "10", 62, false},
		{"large", "g8", 1000, false},
		{"very large", "4c92", 1000000, false},
		{"empty", "", 0, true},
		{"invalid char", "abc!", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("Decode(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Проверяем, что Encode и Decode работают корректно вместе
	testValues := []int64{0, 1, 10, 62, 100, 1000, 10000, 100000, 1000000}

	for _, val := range testValues {
		encoded, err := Encode(val)
		if err != nil {
			t.Errorf("Encode(%d) failed: %v", val, err)
			continue
		}

		decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf("Decode(%q) failed: %v", encoded, err)
			continue
		}

		if decoded != val {
			t.Errorf("Round trip failed: %d -> %q -> %d", val, encoded, decoded)
		}
	}
}
