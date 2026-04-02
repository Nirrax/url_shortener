package utils

import (
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIsUniqueConstraintViolation(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "unique constraint violation",
			err:      &pq.Error{Code: "23505"},
			expected: true,
		},
		{
			name:     "different postgres error",
			err:      &pq.Error{Code: "23503"}, // foreign key violation
			expected: false,
		},
		{
			name:     "generic error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsUniqueConstraintViolation(tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
