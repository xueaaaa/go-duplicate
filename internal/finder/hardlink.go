package finder

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// Hardlink replaces every file in each duplicate group, except the first
// (v.Files[0]), with a hard link to the first; the caller is responsible
// for ordering Files so the file to keep as the link target is first. Note
// that after Hardlink succeeds, removing the first file removes the only
// remaining copy of the data for the whole group.
//
// Files that are already hard-linked to the first file are left alone and
// do not count towards the returned count. Unlike [file.Scan], Hardlink
// does not skip symlinks silently: encountering one aborts the call with
// an error, since a hard link cannot target a symlink meaningfully.
//
// Hardlink stops at the first error and returns it along with the number
// of files successfully linked so far; the remaining files in the current
// group, and any subsequent groups, are left untouched. Linking across
// filesystem boundaries fails with a wrapped syscall.EXDEV error.
//
// Replacement of each target is done via a temporary link followed by a
// rename, so a target is never left partially modified on error.
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

// hardlinkOne replaces target with a hard link to origin by linking into a
// temporary path and renaming over target, so target is never left in a
// partially-modified state. It returns 1 if a link was created, or 0 if
// target was already the same file as origin (no-op) or an error occurred.
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
