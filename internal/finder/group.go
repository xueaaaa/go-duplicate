package finder

import "go-duplicate/internal/file"

type DuplicateGroup struct {
	Hash  [32]byte
	Size  int64
	Files []file.File
}
