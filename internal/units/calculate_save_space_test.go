package units

import (
	"testing"

	"github.com/xueaaaa/go-duplicate/internal/file"
	"github.com/xueaaaa/go-duplicate/internal/finder"
)

func TestCalculatePotentialSaveSpace(t *testing.T) {
	tests := []struct {
		name     string
		groups   []finder.DuplicateGroup
		expected int64
	}{
		{
			name: "primary",
			groups: []finder.DuplicateGroup{
				{
					FileSize: 10 * Byte,
					Files:    make([]file.File, 10),
				},
				{
					FileSize: 10 * KiB,
					Files:    make([]file.File, 9),
				},
				{
					FileSize: 10 * MiB,
					Files:    make([]file.File, 8),
				},
				{
					FileSize: 10 * GiB,
					Files:    make([]file.File, 7),
				},
				{
					FileSize: 10 * TiB,
					Files:    make([]file.File, 6),
				},
			},
			expected: 9*10*Byte + 8*10*KiB + 7*10*MiB + 6*10*GiB + 5*10*TiB,
		},
		{
			name:     "empty group",
			groups:   make([]finder.DuplicateGroup, 0),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := CalculatePotentialSaveSpace(tt.groups)

			if size != tt.expected {
				t.Fatalf("Expected %d got %d", tt.expected, size)
			}
		})
	}
}
