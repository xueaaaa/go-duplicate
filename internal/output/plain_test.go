package output

import (
	"bytes"
	"go-duplicate/internal/file"
	"go-duplicate/internal/finder"
	"strings"
	"testing"
)

func TestPlainOutput(t *testing.T) {
	var buf bytes.Buffer

	groups := []finder.DuplicateGroup{
		{
			FileSize: 10,
			Files: []file.File{
				{Path: "a.txt", Size: 10},
				{Path: "b.txt", Size: 10},
			},
		},
	}

	stats := NewStats("/tmp", 2, groups)
	err := PlainOutput(&buf, groups, stats)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	t.Log(out)

	if !strings.Contains(out, "Duplicate groups: 1") {
		t.Fatal("unexpected output")
	}
}
