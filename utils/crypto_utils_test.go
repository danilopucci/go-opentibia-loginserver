package utils

import (
	"testing"
)

// Test for Sha1Hash function
func TestSha1Hash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},          // Known SHA-1 hash for "hello"
		{"world", "7c211433f02071597741e6ff5a8ea34789abbf43"},          // Known SHA-1 hash for "world"
		{"Go is awesome!", "9804fa3199b57f98bdb480fcbd94b52ae9f88642"}, // Another sample input
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},               // SHA-1 hash for an empty string
	}

	for _, test := range tests {
		result := Sha1Hash(test.input)
		if result != test.expected {
			t.Errorf("Sha1Hash(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}
