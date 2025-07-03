package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing/fstest"
)

func generateFS(root string, numDirs, nestingDepth int) fs.FS {
	var (
		filesPerDir = 10
		fileTypes   = []string{"html", "go", "ts", "css"}
		fsys        = fstest.MapFS{}
	)

	for i := range make([]struct{}, numDirs) {
		currentPath := root

		for d := range make([]struct{}, nestingDepth) {
			currentPath = filepath.Join(currentPath, fmt.Sprintf("dir%d_%d", i, d))

			for f := range make([]struct{}, filesPerDir) {
				ext := fileTypes[f%len(fileTypes)]
				isLayout := f%5 == 0

				filename := fmt.Sprintf("file%d.%s", f, ext)
				if isLayout {
					filename = fmt.Sprintf("layout.%s", ext)
				}

				filePath := filepath.Join(currentPath, filename)
				content := fmt.Appendf(nil, "// dummy content for %s\n", filePath)
				fsys[filePath] = &fstest.MapFile{Data: content, Mode: 0644}
			}
		}
	}
	return fsys
}
