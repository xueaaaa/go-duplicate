package file

type File struct {
	Path       string `json:"path"`
	Size       int64  `json:"size_bytes"`
	IsHardlink bool   `json:"hardlink"`
}
