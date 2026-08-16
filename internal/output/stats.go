package output

import (
	"go-duplicate/internal/finder"
	"go-duplicate/internal/units"
	"time"
)

type Stats struct {
	ScannedDir           string    `json:"scanned_dir"`
	ScannedAt            time.Time `json:"scanned_at"`
	FilesScanned         int64     `json:"files_scanned"`
	DuplicateGroups      int       `json:"duplicate_groups"`
	DuplicateFiles       int       `json:"duplicate_files"`
	PotentialSavingBytes int64     `json:"potential_saving_bytes"`
}

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
