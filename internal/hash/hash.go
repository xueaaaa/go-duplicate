package hash

import (
	"crypto/sha256"
	"io"
	"os"

	"github.com/xueaaaa/go-duplicate/internal/file"
)

// Hash returns the SHA-256 digest of the full contents of file, read
// directly from file.Path. Unlike [Fingerprint], it ignores file.Size and
// streams the entire file through the hash rather than sampling parts of
// it, so it is more expensive but never misses a difference outside a
// sampled region.
//
// Errors from opening or reading the file (including a file that changes
// or disappears between calls) are returned as-is.
func Hash(file file.File) ([]byte, error) {
	f, err := os.Open(file.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
