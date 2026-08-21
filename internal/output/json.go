package output

import (
	"encoding/hex"
	"encoding/json"
	"go-duplicate/internal/file"
	"go-duplicate/internal/finder"
	"io"
)

// SchemaVersion is the current version of the [Report] JSON schema.
// It must be incremented whenever a breaking change is made to the shape
// of Report or its nested types (removing/renaming a field, changing a
// field's type or meaning). Adding a new optional field is not considered
// breaking and does not require a bump.
const SchemaVersion int = 1

// Report is the top-level structure of the tool's JSON output, produced
// by [JSONOutput]. SchemaVersion identifies the shape of this structure
// for consumers that need to handle multiple versions.
type Report struct {
	SchemaVersion int         `json:"schema_version"`
	Stats         Stats       `json:"stats"`
	Groups        []GroupJSON `json:"groups"`
}

// GroupJSON is the JSON representation of a single [finder.DuplicateGroup].
type GroupJSON struct {
	// Hash is the group's content hash, hex-encoded.
	Hash string `json:"hash"`
	// FileSizeBytes is the size of each individual file in the group, in bytes.
	FileSizeBytes int64 `json:"file_size_bytes"`
	// FileCount is the number of files in the group.
	FileCount int `json:"file_count"`
	// TotalSizeBytes is FileSizeBytes * FileCount — the space these
	// duplicates would free up if all but one copy were removed.
	TotalSizeBytes int64 `json:"total_size_bytes"`
	// Files lists the individual duplicate files, in the same order as
	// finder.DuplicateGroup.Files (see that type for ordering caveats).
	Files []file.File `json:"files"`
}

// JSONOutput writes groups and stats to w as a single [Report], encoded
// as tab-indented JSON. Groups are converted to [GroupJSON] in the same
// order they appear in groups; an empty (non-nil) array is written if
// groups is empty.
//
// The written report's SchemaVersion is always [SchemaVersion]. Any error
// from the underlying JSON encoder or from writing to w is returned as-is.
func JSONOutput(w io.Writer, groups []finder.DuplicateGroup, stats Stats) error {
	groupsJson := make([]GroupJSON, 0)

	for _, g := range groups {
		group := GroupJSON{
			Hash:           hex.EncodeToString(g.Hash[:]),
			FileSizeBytes:  g.FileSize,
			FileCount:      len(g.Files),
			TotalSizeBytes: g.FileSize * int64(len(g.Files)),
			Files:          g.Files,
		}
		groupsJson = append(groupsJson, group)
	}

	report := Report{
		SchemaVersion: SchemaVersion,
		Stats:         stats,
		Groups:        groupsJson,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "	")
	return encoder.Encode(report)
}
