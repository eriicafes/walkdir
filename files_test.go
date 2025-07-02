package main

import (
	"os"
	"testing"
)

func BenchmarkFilesWalk(b *testing.B) {
	for b.Loop() {
		fsys := os.DirFS("example")
		filesWithLayouts := WalkFilesWithLayout(fsys, "html", "layout", "app")
		if len(filesWithLayouts) != 5 {
			b.Errorf("wrong file length: expected %d got %d ", 5, len(filesWithLayouts))
		}
	}
}
