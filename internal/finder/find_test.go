package finder

import (
	"context"
	"go-duplicate/internal/params"
	"os"
	"path/filepath"
	"testing"
)

func TestFind(t *testing.T) {
	dir := t.TempDir()
	sameContent1 := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	unique1 := []byte("ABC")
	sameContent2 := []byte("12345678")
	unique2 := []byte("ABCDEFGHIJKLMNOPQRSTUVWXZY")

	path1 := filepath.Join(dir, "first")
	path2 := filepath.Join(dir, "second")
	os.MkdirAll(path1, 0755)
	os.MkdirAll(path2, 0755)
	os.WriteFile(filepath.Join(dir, "unique1.bin"), unique1, 0644)
	os.WriteFile(filepath.Join(path2, "unique2.bin"), unique2, 0644)
	os.WriteFile(filepath.Join(dir, "copy1.bin"), sameContent1, 0644)
	os.WriteFile(filepath.Join(path1, "copy2.bin"), sameContent1, 0644)
	os.WriteFile(filepath.Join(dir, "copy3.bin"), sameContent2, 0644)
	os.WriteFile(filepath.Join(path2, "copy4.bin"), sameContent2, 0644)

	tests := []struct {
		name     string
		expected int
	}{
		{
			name:     "primary",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, groups, err := Find(context.Background(), dir, params.Params{SampleSize: 4})
			if err != nil {
				t.Fatal(err)
			}

			if len(groups) != tt.expected {
				t.Fatalf("expected %x groups, got %x", tt.expected, len(groups))
			}
		})
	}
}
