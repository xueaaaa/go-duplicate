package finder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHardlinkOne(t *testing.T) {
	tests := []struct {
		name          string
		isSameFile    bool
		isSymlink     bool
		errorExpected bool
		expected      int
	}{
		{
			name:          "one duplicate",
			isSameFile:    false,
			isSymlink:     false,
			errorExpected: false,
			expected:      1,
		},
		{
			name:          "same file",
			isSameFile:    true,
			isSymlink:     false,
			errorExpected: false,
			expected:      0,
		},
		{
			name:          "target is symlink",
			isSameFile:    false,
			isSymlink:     true,
			errorExpected: true,
			expected:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			origin := filepath.Join(dir, "origin.bin")
			if err := os.WriteFile(origin, []byte("hello world"), 0644); err != nil {
				t.Fatal(err)
			}

			originInfo, err := os.Lstat(origin)
			if err != nil {
				t.Fatal(err)
			}

			target := filepath.Join(dir, "target.bin")
			switch {
			case tt.isSameFile:
				if err := os.Link(origin, target); err != nil {
					t.Fatal(err)
				}
			case tt.isSymlink:
				if err := os.Symlink(origin, target); err != nil {
					t.Fatal(err)
				}
			default:
				if err := os.WriteFile(target, []byte("hello world"), 0644); err != nil {
					t.Fatal(err)
				}
			}

			got, err := hardlinkOne(origin, originInfo, target)
			if tt.errorExpected && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.errorExpected && err != nil {
				t.Fatal("unexpected error", err)
			}
			if got != tt.expected {
				t.Fatalf("expected %d got %d", tt.expected, got)
			}
		})
	}
}
