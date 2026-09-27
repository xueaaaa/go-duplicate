package finder

import "github.com/xueaaaa/go-duplicate/internal/file"

// DuplicateGroup is a set of files with identical content, as found by
// [Find].
type DuplicateGroup struct {
	// Hash is the full-content hash ([hash2.Hash]) shared by every file
	// in Files.
	Hash [32]byte
	// FileSize is the size in bytes shared by every file in Files.
	FileSize int64
	// Files are the files that make up this duplicate group, sorted by Path
	// ascending. There are always at least two.
	Files []file.File
}
