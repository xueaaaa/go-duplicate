package file

import (
	"context"
	"fmt"
	"go-duplicate/internal/util"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

func Scan(ctx context.Context, dir string, out chan<- File) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if util.IsSkippableFSError(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		entryPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			err = Scan(ctx, entryPath, out)
			if err != nil {
				return err
			}
		} else {
			if entry.Type()&fs.ModeType != 0 {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				if util.IsSkippableFSError(err) {
					continue
				}
				return err
			}

			stat, ok := info.Sys().(*syscall.Stat_t)
			if !ok {
				return fmt.Errorf("failed to get syscall.Stat_t for %s", entryPath)
			}

			file := File{
				Path:       entryPath,
				Size:       info.Size(),
				IsHardlink: stat.Nlink > 1,
			}

			select {
			case out <- file:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}
