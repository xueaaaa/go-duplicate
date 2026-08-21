package finder_test

import (
	"fmt"
	"go-duplicate/internal/file"
	"go-duplicate/internal/finder"
	"os"
	"path/filepath"
	"testing"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name          string
		filesPerGroup []int
		expectedCount int
		expectError   bool
		preDelete     *[2]int
	}{
		{
			name:          "one group with two elements",
			filesPerGroup: []int{2},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "two groups with three elements",
			filesPerGroup: []int{3, 3},
			expectedCount: 4,
			expectError:   false,
		},
		{
			name:          "one group with one file",
			filesPerGroup: []int{1},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "no groups",
			filesPerGroup: []int{},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "error (pre-deleted)",
			filesPerGroup: []int{3, 4},
			expectedCount: 3,
			expectError:   true,
			preDelete:     &[2]int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			groups := make([]finder.DuplicateGroup, 0)

			for i, v := range tt.filesPerGroup {
				files := make([]file.File, 0)

				for k := 0; k < v; k++ {
					path := filepath.Join(dir, fmt.Sprintf("%d_%d.bin", i, k))

					if err := os.WriteFile(path, []byte("hello world"), 0644); err != nil {
						t.Fatal(err)
					}

					files = append(files, file.File{
						Path: path,
					})
				}

				groups = append(groups, finder.DuplicateGroup{
					Files: files,
				})
			}

			if tt.preDelete != nil {
				i, k := tt.preDelete[0], tt.preDelete[1]
				if err := os.Remove(groups[i].Files[k].Path); err != nil {
					t.Fatal(err)
				}
			}

			count, err := finder.Delete(groups)
			if err != nil {
				if !tt.expectError {
					t.Fatal(err)
				}

				if count != tt.expectedCount {
					t.Fatalf("expected %d got %d (with error)", tt.expectedCount, count)
				}
			}

			if tt.expectedCount != count {
				t.Fatalf("expected %d got %d", tt.expectedCount, count)
			}
		})
	}
}
