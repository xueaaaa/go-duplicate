package hash_test

import (
	"bytes"
	"crypto/sha256"
	"go-duplicate/internal/file"
	"go-duplicate/internal/hash"
	"go-duplicate/internal/units"
	"os"
	"path/filepath"
	"testing"
)

func TestFingerprint(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name          string
		fileSize      int64
		sampleSize    int64
		errorExpected bool
		expected      []byte
	}{
		{
			name:          "small_file",
			fileSize:      4 * units.Byte,
			sampleSize:    4 * units.Byte,
			errorExpected: false,
		},
		{
			name:          "large_sample_size",
			fileSize:      4 * units.Byte,
			sampleSize:    10 * units.Byte,
			errorExpected: false,
		},
		{
			name:          "large_file",
			fileSize:      12 * units.KiB,
			sampleSize:    4 * units.KiB,
			errorExpected: false,
		},
		{
			name:          "empty_file",
			fileSize:      0 * units.Byte,
			sampleSize:    4 * units.Byte,
			errorExpected: false,
		},
		{
			name:          "missing_file",
			fileSize:      0 * units.Byte,
			sampleSize:    0 * units.Byte,
			errorExpected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := make([]byte, 0)
			for i := int64(0); i < tt.fileSize; i++ {
				content = append(content, byte('A'))
			}

			path := filepath.Join(dir, tt.name+".bin")
			if tt.fileSize > 0*units.Byte || (tt.fileSize == 0*units.Byte && !tt.errorExpected) {
				if err := os.WriteFile(path, content, 0644); err != nil {
					t.Fatal(err)
				}
			}

			f := file.File{
				Path: path,
				Size: tt.fileSize,
			}

			h := sha256.New()
			if int64(len(content)) <= tt.sampleSize {
				h.Write(content)
			} else {
				h.Write(content[:tt.sampleSize])
				h.Write(content[int64(len(content))-tt.sampleSize:])
			}
			tt.expected = h.Sum(nil)

			got, err := hash.Fingerprint(f, tt.sampleSize)
			if tt.errorExpected {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected nil-error got %s", err)
			}
			if !bytes.Equal(got, tt.expected) {
				t.Fatalf("expected %x, got %x", tt.expected, got)
			}
		})
	}
}
