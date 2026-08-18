package output

import (
	"encoding/hex"
	"encoding/json"
	"go-duplicate/internal/file"
	"go-duplicate/internal/finder"
	"io"
)

const SCHEMA_VERSION int = 1

type Report struct {
	SchemaVersion int         `json:"schema_version"`
	Stats         Stats       `json:"stats"`
	Groups        []GroupJSON `json:"groups"`
}

type GroupJSON struct {
	Hash           string      `json:"hash"`
	FileSizeBytes  int64       `json:"file_size_bytes"`
	FileCount      int         `json:"file_count"`
	TotalSizeBytes int64       `json:"total_size_bytes"`
	Files          []file.File `json:"files"`
}

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
		SchemaVersion: SCHEMA_VERSION,
		Stats:         stats,
		Groups:        groupsJson,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "	")
	return encoder.Encode(report)
}
