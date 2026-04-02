package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase62Encoder_Encode(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"zero", 0, "0"},
		{"one", 1, "1"},
		{"nine", 9, "9"},
		{"ten maps to A", 10, "A"},
		{"thirty five maps to Z", 35, "Z"},
		{"thirty six maps to a", 36, "a"},
		{"sixty one maps to z", 61, "z"},
		{"sixty two maps to 10", 62, "10"},
		{"large number", 3844, "100"},
		{"max uint64", ^uint64(0), "LygHa16AHYF"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := NewBase62Encoder()
			assert.Equal(t, tc.expected, e.Encode(tc.input))
		})
	}
}

func TestBase62Encoder_Decode(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      uint64
		expectError   bool
		errorContains string
	}{
		{"zero", "0", 0, false, ""},
		{"one", "1", 1, false, ""},
		{"nine", "9", 9, false, ""},
		{"uppercase A", "A", 10, false, ""},
		{"uppercase Z", "Z", 35, false, ""},
		{"lowercase a", "a", 36, false, ""},
		{"lowercase z", "z", 61, false, ""},
		{"two digits", "10", 62, false, ""},
		{"three digits", "100", 3844, false, ""},
		{"max uint64", "LygHa16AHYF", ^uint64(0), false, ""},
		{"invalid character dash", "a-b", 0, true, "invalid character: -"},
		{"invalid character space", "a b", 0, true, "invalid character:  "},
		{"invalid character symbol", "abc!", 0, true, "invalid character: !"},
		{"overflow", "zzzzzzzzzzz", 0, true, "overflows uint64"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := NewBase62Encoder()
			result, err := e.Decode(tc.input)
			if tc.expectError {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestBase62Encoder_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input uint64
	}{
		{"zero", 0},
		{"one", 1},
		{"sixty one", 61},
		{"sixty two", 62},
		{"large number", 123456789},
		{"max uint64", ^uint64(0)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := NewBase62Encoder()
			encoded := e.Encode(tc.input)
			decoded, err := e.Decode(encoded)
			assert.NoError(t, err)
			assert.Equal(t, tc.input, decoded)
		})
	}
}
