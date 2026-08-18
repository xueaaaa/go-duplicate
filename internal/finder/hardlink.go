package finder

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

func Hardlink(groups []DuplicateGroup) (int, error) {
	count := 0

	for _, v := range groups {
		origin := v.Files[0].Path

		originInfo, err := os.Lstat(origin)
		if err != nil {
			return count, err
		}

		for _, file := range v.Files[1:] {
			n, err := hardlinkOne(origin, originInfo, file.Path)
			count += n
			if err != nil {
				return count, err
			}
		}
	}

	return count, nil
}

func hardlinkOne(origin string, originInfo os.FileInfo, target string) (int, error) {
	targetInfo, err := os.Lstat(target)
	if err != nil {
		return 0, err
	}

	if os.SameFile(originInfo, targetInfo) {
		return 0, nil
	}

	if targetInfo.Mode()&os.ModeSymlink != 0 {
		return 0, fmt.Errorf("skip symlink %q", target)
	}

	tmp := target + ".hardlink.tmp"
	_ = os.Remove(tmp)

	if err := os.Link(origin, tmp); err != nil {
		if errors.Is(err, syscall.EXDEV) {
			return 0, fmt.Errorf("cross-device link %q -> %q: %w", origin, target, err)
		}
		return 0, fmt.Errorf("link %q -> %q: %w", origin, tmp, err)
	}

	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return 0, fmt.Errorf("rename %q -> %q: %w", tmp, target, err)
	}

	return 1, nil
}
