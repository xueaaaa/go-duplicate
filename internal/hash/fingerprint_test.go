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
	contentSF := []byte("ABC")
	pathSF := filepath.Join(dir, "test_small_file.bin")
	os.WriteFile(pathSF, contentSF, 0644)
	contentLF := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	pathLF := filepath.Join(dir, "test_large_file.bin")
	os.WriteFile(pathLF, contentLF, 0644)

	const sampleSize = int64(4)

	h := sha256.New()
	h.Write(contentSF)
	expectedSF := h.Sum(nil)
	h.Reset()
	h.Write(contentLF[:sampleSize])
	h.Write(contentLF[int64(len(contentLF))-sampleSize:])
	expectedLF := h.Sum(nil)

	tests := []struct {
		name     string
		file     file.File
		expected []byte
	}{
		{
			name:     "small_file",
			file:     file.File{Path: pathSF, Size: int64(len(contentSF))},
			expected: expectedSF,
		},
		{
			name:     "large_file",
			file:     file.File{Path: pathLF, Size: int64(len(contentLF))},
			expected: expectedLF,
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
