package finder

import "go-duplicate/internal/file"

// DuplicateGroup is a set of files with identical content, as found by
// [Find].
type DuplicateGroup struct {
	// Hash is the full-content hash ([hash2.Hash]) shared by every file
	// in Files.
	Hash [32]byte
	// FileSize is the size in bytes shared by every file in Files.
	FileSize int64
	// Files are the files that make up this duplicate group. There are
	// always at least two.
	//
	// The order of Files is not guaranteed to be stable across calls to
	// Find; if a specific file must be preserved by [Delete], the caller
	// is responsible for reordering Files before calling Delete.
	Files []file.File
}
