package util

import (
	"errors"
	"io/fs"
	"syscall"
)

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
