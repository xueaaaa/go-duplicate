package file

import (
	"context"
	"os"
	"path/filepath"
)

func Scan(ctx context.Context, dir string, out chan<- File) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
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
			info, err := entry.Info()
			if err != nil {
				return err
			}

			file := File{
				Path: entryPath,
				Size: info.Size(),
			}

			out <- file
		}
	}

	return nil
}
