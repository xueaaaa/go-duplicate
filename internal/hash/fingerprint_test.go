package hash

import (
	"bytes"
	"crypto/sha256"
	"go-duplicate/internal/file"
	"os"
	"path/filepath"
	"testing"
)

func TestFingerprint(t *testing.T) {
	dir := t.TempDir()
	content := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	path := filepath.Join(dir, "test.bin")
	os.WriteFile(path, content, 0644)

	const sampleSize = int64(4)

	h := sha256.New()
	h.Write(content[:sampleSize])
	h.Write(content[int64(len(content))-sampleSize:])
	expected := h.Sum(nil)

	tests := []struct {
		name     string
		file     file.File
		expected []byte
	}{
		{
			name:     "large_file",
			file:     file.File{Path: path, Size: int64(len(content))},
			expected: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Fingerprint(tt.file, sampleSize)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, tt.expected) {
				t.Fatalf("expected %x, got %x", tt.expected, got)
			}
		})
	}
}
