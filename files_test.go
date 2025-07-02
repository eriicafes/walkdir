package main

import (
	"os"
	"slices"
	"strconv"
	"testing"

	"io/fs"
	"testing/fstest"
)

func BenchmarkWalkSmall(b *testing.B) {
	fsys := os.DirFS("example")
	visited := 0

	for b.Loop() {
		err := fs.WalkDir(fsys, "data_sm", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			visited++
			return nil
		})
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
	if visited == 0 {
		b.Fatal("no files visited")
	}
}

func BenchmarkWalkBreadthFirstSmall(b *testing.B) {
	fsys := os.DirFS("example")
	visited := 0

	for b.Loop() {
		err := walkDirBreadthFirst(fsys, "data_sm", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			visited++
			return nil
		})
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
	if visited == 0 {
		b.Fatal("no files visited")
	}
}

func BenchmarkWalkLarge(b *testing.B) {
	fsys := os.DirFS("example")
	visited := 0

	for b.Loop() {
		err := fs.WalkDir(fsys, "data_lg", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			visited++
			return nil
		})
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
	if visited == 0 {
		b.Fatal("no files visited")
	}
}

func BenchmarkWalkBreadthFirstLarge(b *testing.B) {
	fsys := os.DirFS("example")
	visited := 0

	for b.Loop() {
		err := walkDirBreadthFirst(fsys, "data_lg", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			visited++
			return nil
		})
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
	if visited == 0 {
		b.Fatal("no files visited")
	}
}

func BenchmarkWalkLayoutSmall(b *testing.B) {
	fsys := os.DirFS("example")
	countBytes, _ := fs.ReadFile(fsys, "data_sm.txt")
	count, err := strconv.Atoi(string(countBytes))
	if err != nil {
		b.Fatalf("failed to get count: %v", err)
	}

	for b.Loop() {
		filesWithLayouts := WalkFilesWithLayout(fsys, "html", "layout", "data_sm")
		if len(filesWithLayouts) != count {
			b.Errorf("wrong file length: expected %d got %d ", count, len(filesWithLayouts))
		}
	}
}

func BenchmarkWalkLayoutBreadthFirstSmall(b *testing.B) {
	fsys := os.DirFS("example")
	countBytes, _ := fs.ReadFile(fsys, "data_sm.txt")
	count, err := strconv.Atoi(string(countBytes))
	if err != nil {
		b.Fatalf("failed to get count: %v", err)
	}

	for b.Loop() {
		filesWithLayouts := WalkFilesWithLayoutBreadthFirst(fsys, "html", "layout", "data_sm")
		if len(filesWithLayouts) != count {
			b.Errorf("wrong file length: expected %d got %d ", count, len(filesWithLayouts))
		}
	}
}

func BenchmarkWalkLayoutLarge(b *testing.B) {
	fsys := os.DirFS("example")
	countBytes, _ := fs.ReadFile(fsys, "data_lg.txt")
	count, err := strconv.Atoi(string(countBytes))
	if err != nil {
		b.Fatalf("failed to get count: %v", err)
	}

	for b.Loop() {
		filesWithLayouts := WalkFilesWithLayout(fsys, "html", "layout", "data_lg")
		if len(filesWithLayouts) != count {
			b.Errorf("wrong file length: expected %d got %d ", count, len(filesWithLayouts))
		}
	}
}

func BenchmarkWalkLayoutBreadthFirstLarge(b *testing.B) {
	fsys := os.DirFS("example")
	countBytes, _ := fs.ReadFile(fsys, "data_lg.txt")
	count, err := strconv.Atoi(string(countBytes))
	if err != nil {
		b.Fatalf("failed to get count: %v", err)
	}

	for b.Loop() {
		filesWithLayouts := WalkFilesWithLayoutBreadthFirst(fsys, "html", "layout", "data_lg")
		if len(filesWithLayouts) != count {
			b.Errorf("wrong file length: expected %d got %d ", count, len(filesWithLayouts))
		}
	}
}

func TestWalkDirBreadthFirst(t *testing.T) {
	memFS := fstest.MapFS{
		"root/file1.txt":          {Data: []byte("")},
		"root/dirA/file1.txt":     {Data: []byte("")},
		"root/dirB/file1.txt":     {Data: []byte("")},
		"root/dirB/sub/file1.txt": {Data: []byte("")},
	}

	var visited []string
	err := walkDirBreadthFirst(memFS, "root", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		visited = append(visited, path)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"root",
		"root/dirA",
		"root/dirB",
		"root/file1.txt",
		"root/dirA/file1.txt",
		"root/dirB/file1.txt",
		"root/dirB/sub",
		"root/dirB/sub/file1.txt",
	}
	if !slices.Equal(visited, expected) {
		t.Errorf("expected:\n  %v\ngot\n: %v", expected, visited)
	}
}

func TestWalkDirBreadthFirst_SkipAll(t *testing.T) {
	memFS := fstest.MapFS{
		"root/file1.txt":      {Data: []byte("")},
		"root/dirA/file1.txt": {Data: []byte("")},
	}

	var visited []string
	err := walkDirBreadthFirst(memFS, "root", func(path string, d fs.DirEntry, err error) error {
		visited = append(visited, path)
		return fs.SkipAll
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"root",
	}
	if !slices.Equal(visited, expected) {
		t.Errorf("expected:\n  %v\ngot\n: %v", expected, visited)
	}
}

func TestWalkDirBreadthFirst_SkipDir(t *testing.T) {
	memFS := fstest.MapFS{
		"root/file1.txt":          {Data: []byte("")},
		"root/dirA/file1.txt":     {Data: []byte("")},
		"root/dirB/afile1.txt":    {Data: []byte("")},
		"root/dirB/file1.txt":     {Data: []byte("")},
		"root/dirB/sub/file1.txt": {Data: []byte("")},
		"root/dirB/zfile1.txt":    {Data: []byte("")},
	}

	// skipped entire directory
	var visited []string
	err := walkDirBreadthFirst(memFS, "root", func(path string, d fs.DirEntry, err error) error {
		visited = append(visited, path)
		if path == "root/dirB" {
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"root",
		"root/dirA",
		"root/dirB",
		"root/file1.txt",
		"root/dirA/file1.txt",
	}
	if !slices.Equal(visited, expected) {
		t.Errorf("expected:\n  %v\ngot\n: %v", expected, visited)
	}

	// skipped path appears before some
	visited = nil
	err = walkDirBreadthFirst(memFS, "root", func(path string, d fs.DirEntry, err error) error {
		visited = append(visited, path)
		if path == "root/dirB/afile1.txt" {
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected = []string{
		"root",
		"root/dirA",
		"root/dirB",
		"root/file1.txt",
		"root/dirA/file1.txt",
		"root/dirB/afile1.txt",
	}
	if !slices.Equal(visited, expected) {
		t.Errorf("expected:\n  %v\ngot\n: %v", expected, visited)
	}

	// skipped path appears after some
	visited = nil
	err = walkDirBreadthFirst(memFS, "root", func(path string, d fs.DirEntry, err error) error {
		visited = append(visited, path)
		if path == "root/dirB/zfile1.txt" {
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected = []string{
		"root",
		"root/dirA",
		"root/dirB",
		"root/file1.txt",
		"root/dirA/file1.txt",
		"root/dirB/afile1.txt",
		"root/dirB/file1.txt",
		"root/dirB/sub",
		"root/dirB/zfile1.txt",
	}
	if !slices.Equal(visited, expected) {
		t.Errorf("expected:\n  %v\ngot\n: %v", expected, visited)
	}
}
