package hash_test

import (
	"bytes"
	"crypto/sha256"
	"go-duplicate/internal/file"
	"go-duplicate/internal/hash"
	"os"
	"path/filepath"
	"testing"
)

func TestHash(t *testing.T) {
	dir := t.TempDir()
	content := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	path := filepath.Join(dir, "test.bin")
	os.WriteFile(path, content, 0644)

	h := sha256.New()
	h.Write(content)
	expected := h.Sum(nil)

	tests := []struct {
		name     string
		file     file.File
		expected []byte
	}{
		{
			name:     "primary",
			file:     file.File{Path: path, Size: int64(len(content))},
			expected: expected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hash.Hash(tt.file)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(got, tt.expected) {
				t.Fatalf("expected %x, got %x", tt.expected, got)
			}
		})
	}
}
