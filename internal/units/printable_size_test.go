package units

import "testing"

func TestPrintableSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "< 0",
			size:     -1025,
			expected: "-1025 B",
		},
		{
			name:     "< KiB",
			size:     1023,
			expected: "1023 B",
		},
		{
			name:     "[KiB; MiB)",
			size:     15 * KiB,
			expected: "15.00 KiB",
		},
		{
			name:     "[MiB; GiB)",
			size:     999 * MiB,
			expected: "999.00 MiB",
		},
		{
			name:     "[GiB; TiB)",
			size:     1 * GiB,
			expected: "1.00 GiB",
		},
		{
			name:     ">= TiB",
			size:     1 * TiB,
			expected: "1.00 TiB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := PrintableSize(tt.size)

			if str != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, str)
			}
		})
	}
}
