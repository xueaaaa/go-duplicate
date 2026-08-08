package hash

import (
	"crypto/sha256"
	"go-duplicate/internal/file"
	"io"
	"os"
)

func Fingerprint(file file.File, sampleSize int64) ([]byte, error) {
	if file.Size <= 2*sampleSize {
		return Hash(file)
	}

	f, err := os.Open(file.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	head := make([]byte, sampleSize)
	if _, err := io.ReadFull(f, head); err != nil {
		return nil, err
	}

	tail := make([]byte, sampleSize)
	if _, err = f.Seek(-sampleSize, io.SeekEnd); err != nil {
		return nil, err
	}
	if _, err = io.ReadFull(f, tail); err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(head)
	h.Write(tail)

	return h.Sum(nil), nil
}
