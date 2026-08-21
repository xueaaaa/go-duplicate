package params

// Format identifies how scan results should be rendered.
type Format string

const (
	// Plain renders results as human-readable plain text
	Plain Format = "plain"
	// JSON renders results as machine-readable JSON.
	JSON Format = "json"
)

var validFormats = map[Format]bool{
	Plain: true,
	JSON:  true,
}

// IsValid reports whether f is one of the recognized output formats.
func (f Format) IsValid() bool {
	return validFormats[f]
}
