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

	buffer := make([]byte, file.Size)
	if _, err = io.ReadFull(f, buffer); err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(buffer)

	return h.Sum(nil), nil
}
