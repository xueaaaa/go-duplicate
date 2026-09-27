package output

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xueaaaa/go-duplicate/internal/finder"
	"github.com/xueaaaa/go-duplicate/internal/units"
)

// PlainOutput writes a human-readable summary of groups and stats to w:
// scan metadata, aggregate stats, potential space savings, and then a
// numbered listing of each duplicate group with its files. This format is
// meant for terminal/log output, not for parsing — for a stable,
// versioned machine-readable format use [JSONOutput] instead.
//
// If groups is empty, PlainOutput prints the stats followed by
// "No duplicate files found" and returns, without a groups section.
//
// Files within each group are listed in the order they appear in
// group.Files (see [finder.DuplicateGroup] for ordering caveats).
// "Summary size" for a group is FileSize * len(Files) — the combined size
// of all copies, not just the reclaimable space.
//
// PlainOutput stops at the first write error and returns it; on error, w
// may contain a partially written report.
func PlainOutput(w io.Writer, groups []finder.DuplicateGroup, stats Stats) error {
	saveSpaceStr := units.PrintableSize(stats.PotentialSavingBytes)

	if _, err := fmt.Fprintf(w, "\nScanning complete (directory: %s)\n", stats.ScannedDir); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "\nScanning complete at: %s\n", stats.ScannedAt); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "\nFiles scanned: %d\n", stats.FilesScanned); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Duplicate groups: %d\n", stats.DuplicateGroups); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Duplicate files: %d\n", stats.DuplicateFiles); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "Potential space savings:", saveSpaceStr); err != nil {
		return err
	}

	if len(groups) == 0 {
		if _, err := fmt.Fprintln(w, "\nNo duplicate files found"); err != nil {
			return err
		}
		return nil
	}

	for i, group := range groups {
		if _, err := fmt.Fprintf(w, "\n[%d]\n", i+1); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "File size:", units.PrintableSize(group.FileSize)); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "Summary size:",
			units.PrintableSize(int64(len(group.Files))*group.FileSize)); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "Hash:", hex.EncodeToString(group.Hash[:])); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "Files:"); err != nil {
			return err
		}

		for j, file := range group.Files {
			if _, err := fmt.Fprintf(w, "\t[%d] - %s\n", j+1, file.Path); err != nil {
				return err
			}
		}
	}

	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}

	return nil
}
