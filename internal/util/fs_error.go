package util

import (
	"errors"
	"io/fs"
	"syscall"
)

// IsSkippableFSError reports whether err represents a filesystem condition
// that scanning/hashing code in this module should silently skip rather
// than treat as fatal — used by [file.Scan], [finder.Find], and
// [finder.Hardlink] to decide whether to continue past an inaccessible or
// unusual file/directory.
//
// It returns true for:
//   - fs.ErrPermission — the file or directory is not accessible.
//   - fs.ErrNotExist — the path no longer exists (e.g. removed after
//     being listed, but before being opened).
//   - syscall.ENOTDIR — part of the path stopped being a directory
//     between being listed and being accessed (a race with concurrent
//     filesystem changes).
//   - syscall.ENXIO, syscall.ENODEV — the path refers to a special file
//     (e.g. a device node) that cannot be read as a regular file.
//
// Any other error, including wrapped variants that don't match one of the
// above via errors.Is/errors.As, returns false.
func IsSkippableFSError(err error) bool {
	if errors.Is(err, fs.ErrPermission) || errors.Is(err, fs.ErrNotExist) {
		return true
	}
	if errno, ok := errors.AsType[syscall.Errno](err); ok {
		switch errno {
		case syscall.ENXIO,
			syscall.ENODEV,
			syscall.ENOTDIR:
			return true
		}
	}
	return false
}
