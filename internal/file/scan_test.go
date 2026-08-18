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

	expected := 3

	t.Run("primary", func(t *testing.T) {
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

		if count != expected {
			t.Fatalf("expected %d files, got %d", expected, count)
		}
	})
}

func TestScan_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	expected := 0

	t.Run("empty directory", func(t *testing.T) {
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

		if count != expected {
			t.Fatalf("expected %d files, got %d", expected, count)
		}
	})
}

func TestScan_DeepRecursion(t *testing.T) {
	t.Run("deep recursion", func(t *testing.T) {
		dir := t.TempDir()

		os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0644)
		os.Mkdir(filepath.Join(dir, "a"), 0755)
		os.WriteFile(filepath.Join(dir, "a", "c.txt"), []byte("world"), 0644)
		os.Mkdir(filepath.Join(dir, "a", "b"), 0755)
		os.WriteFile(filepath.Join(dir, "a", "b", "c.txt"), []byte("hello"), 0644)
		os.Mkdir(filepath.Join(dir, "a", "b", "c"), 0755)
		os.Mkdir(filepath.Join(dir, "a", "b", "c", "d"), 0755)
		os.WriteFile(filepath.Join(dir, "a", "b", "c", "d", "e.txt"), []byte("world"), 0644)

		expected := 4

		t.Run("deep recursion", func(t *testing.T) {
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

			if count != expected {
				t.Fatalf("expected %d files, got %d", expected, count)
			}
		})
	})
}
