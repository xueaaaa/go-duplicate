package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("world"), 0644)
	os.Mkdir(filepath.Join(dir, "a"), 0755)
	os.WriteFile(filepath.Join(dir, "a", "c.txt"), []byte("hello"), 0644)

	tests := []struct {
		name     string
		expected int
	}{
		{
			name:     "primary",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileChan := make(chan File)

			go func() {
				defer close(fileChan)

				if err := Scan(context.Background(), dir, fileChan); err != nil {
					t.Error(err)
				}
			}()

			count := 0
			for file := range fileChan {
				t.Log(file)
				count++
			}

			if count != tt.expected {
				t.Fatalf("expected %d files, got %d", tt.expected, count)
			}
		})
	}
}
