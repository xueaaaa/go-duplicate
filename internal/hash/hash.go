package hash

import (
	"crypto/sha256"
	"go-duplicate/internal/file"
	"io"
	"os"
)

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
