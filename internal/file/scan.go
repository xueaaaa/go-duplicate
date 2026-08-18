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

// Scan recursively walks dir and sends every regular file it finds to out.
// It blocks on each send, so out must be drained concurrently (e.g. by a
// goroutine reading from it) or Scan will deadlock. Scan does not close out.
//
// Non-existent directories and directories/files with denied access are
// silently skipped (see util.IsSkippableFSError). Symlinks, sockets, devices,
// and other non-regular files are skipped as well. Any other filesystem
// error aborts the scan and is returned.
//
// Scan checks ctx before each directory read and before each send; if ctx
// is cancelled, it stops and returns ctx.Err().
//
// # Usage
//
//	fileChan := make(chan File)
//	go func() {
//	    defer close(fileChan)
//	    if err := Scan(ctx, dir, fileChan); err != nil {
//	        log.Println(err)
//	    }
//	}()
//	for f := range fileChan {
//	    // process f
//	}
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
