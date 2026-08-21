package output

import (
	"go-duplicate/internal/finder"
	"go-duplicate/internal/units"
	"time"
)

// Stats summarizes a single scan run. It is used both in [JSONOutput] and
// [PlainOutput], and its JSON shape is part of [Report].
type Stats struct {
	// ScannedDir is the directory path as passed to the scan (not
	// resolved to an absolute path).
	ScannedDir string `json:"scanned_dir"`
	// ScannedAt is the time the scan started.
	ScannedAt time.Time `json:"scanned_at"`
	// FilesScanned is the total number of files scanned, as returned by
	// [finder.Find] — regardless of whether duplicates were found.
	FilesScanned int64 `json:"files_scanned"`
	// DuplicateGroups is the number of duplicate groups found (len(groups)).
	DuplicateGroups int `json:"duplicate_groups"`
	// DuplicateFiles is the total number of redundant copies across all
	// duplicate groups — for each group, len(Files)-1, summed across
	// groups. The one file per group that would be kept is not counted.
	DuplicateFiles int `json:"duplicate_files"`
	// PotentialSavingBytes is the disk space that could be reclaimed by
	// keeping only one copy per group: for each group,
	// FileSize * (len(Files) - 1), summed across all groups.
	PotentialSavingBytes int64 `json:"potential_saving_bytes"`
}

// NewStats builds a [Stats] snapshot for a completed scan of dir. filesScanned
// and groups should be the values returned by [finder.Find]; ScannedAt is
// set to the current time, i.e. when NewStats is called, not when the scan
// began.
//
// DuplicateGroups is len(groups); DuplicateFiles and PotentialSavingBytes
// are derived by excluding one "kept" file per group (len(Files)-1), so
// both reflect only the redundant copies, not the full group size.
func NewStats(
	dir string,
	filesScanned int64,
	groups []finder.DuplicateGroup,
) Stats {
	duplicateFiles := 0
	for _, g := range groups {
		duplicateFiles += len(g.Files) - 1
	}

	return Stats{
		ScannedDir:           dir,
		ScannedAt:            time.Now(),
		FilesScanned:         filesScanned,
		DuplicateGroups:      len(groups),
		DuplicateFiles:       duplicateFiles,
		PotentialSavingBytes: units.CalculatePotentialSaveSpace(groups),
	}
}
