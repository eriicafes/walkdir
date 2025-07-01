package main

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

// WalkFilesWithLayout walks a directory and for each filename that matches the file extension ext,
// returns a slice of all layout filenames (without extension) in parent directories and the matched filename (without extension).
func WalkFilesWithLayout(fsys fs.FS, ext string, layoutFilename string, dir string) map[string][]string {
	groups := make(map[string][]string)
	layouts := make([]string, 0)
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return err
		}
		pathWithoutExt := strings.TrimSuffix(path, "."+ext)
		// if pathWithoutExt remained unchanged then path does not have ext
		if ext != "" && path == pathWithoutExt {
			return err
		}
		_, filename := filepath.Split(pathWithoutExt)
		if filename == layoutFilename {
			layouts = append(layouts, pathWithoutExt)
		} else if dir == "." || strings.HasPrefix(pathWithoutExt, dir) {
			groups[pathWithoutExt] = []string{pathWithoutExt}
		}
		return err
	})

	if len(layouts) < 1 {
		return groups
	}
	slices.SortFunc(layouts, func(a, b string) int {
		return len(a) - len(b)
	})
	for name := range groups {
		files := []string{}
		fileDir, _ := filepath.Split(name)
		for _, layout := range layouts {
			layoutDir, _ := filepath.Split(layout)
			if strings.HasPrefix(fileDir, layoutDir) {
				files = append(files, layout)
			}
			// no need to check deeper layout files
			if layoutDir == fileDir {
				break
			}
		}
		groups[name] = append(files, groups[name]...)
	}
	return groups
}
