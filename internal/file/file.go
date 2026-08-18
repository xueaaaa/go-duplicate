package file

// File describes a single regular file found by [Scan].
type File struct {
	// Path is the file's path, as passed to Scan joined with the entry name
	// (absolute if dir was absolute, relative otherwise).
	Path string `json:"path"`
	// Size is the file size in bytes.
	Size int64 `json:"size_bytes"`
	// IsHardlink reports whether the file's inode has more than one hard
	// link (i.e. syscall.Stat_t.Nlink > 1).
	IsHardlink bool `json:"hardlink"`
}
