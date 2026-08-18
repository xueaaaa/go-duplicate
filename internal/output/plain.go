package output

import (
	"encoding/hex"
	"fmt"
	"go-duplicate/internal/finder"
	"go-duplicate/internal/units"
	"io"
)

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
