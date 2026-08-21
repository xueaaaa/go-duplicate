package hash

import (
	"crypto/sha256"
	"go-duplicate/internal/file"
	"io"
	"os"
)

// Fingerprint returns a cheap, partial-content digest of file, suitable for
// quickly narrowing down duplicate candidates before a full [Hash].
//
// If file.Size <= 2*sampleSize, Fingerprint reads and hashes the entire
// file via [Hash] — for small files there's no cheaper partial read to do.
// Otherwise it reads the first and last sampleSize bytes of the file and
// returns sha256(head || tail); everything in between is not read or
// hashed, so two files can share a fingerprint while differing in the
// middle.
//
// sampleSize must be non-negative; a sampleSize of 0 causes Fingerprint to
// hash an empty head and tail for any file larger than 0 bytes, meaning
// all such files of a given size class will produce the same fingerprint
// regardless of content. Negative values are not supported and will cause
// a read error.
//
// Fingerprint reads file.Path directly and uses file.Size only to decide
// which strategy to use; if the file has changed size since file.Size was
// recorded, reading may fail with an error from the underlying I/O.
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
