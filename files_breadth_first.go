package main

import (
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

func WalkFilesWithLayoutBreadthFirst(fsys fs.FS, ext string, layoutFilename string, dir string) map[string][]string {
	groups := make(map[string][]string)
	layouts := make([]string, 0)
	subdir, sublayout, subgroups := "", false, make(map[string][]string)

	walkDirBreadthFirst(fsys, ".", func(path string, d fs.DirEntry, err error) error {
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
		fileDir, filename := filepath.Split(pathWithoutExt)

		// commit last subgroup when switching dirs
		if fileDir != subdir {
			for k, v := range subgroups {
				if sublayout {
					v = append(v, filepath.Join(subdir, layoutFilename))
				}
				groups[k] = append(v, k)
			}
			if sublayout {
				layouts = append(layouts, subdir)
			}
			subdir, sublayout, subgroups = fileDir, false, make(map[string][]string)
		}

		if filename == layoutFilename {
			sublayout = true
		} else if dir == "." || strings.HasPrefix(pathWithoutExt, dir) {
			var files []string
			// add parent layout files that share common prefix
			for _, layout := range layouts {
				if strings.HasPrefix(fileDir, layout) {
					files = append(files, filepath.Join(layout, layoutFilename))
				}
				// no need to check deeper layout files
				if layout == fileDir {
					break
				}
			}
			subgroups[pathWithoutExt] = files
		}
		return err
	})
	// commit remaining subgroup
	if len(subgroups) != 0 {
		for k, v := range subgroups {
			if sublayout {
				v = append(v, filepath.Join(subdir, layoutFilename))
			}
			groups[k] = append(v, k)
		}
	}
	return groups
}

func walkDirBreadthFirst(fsys fs.FS, root string, fn fs.WalkDirFunc) error {
	info, err := fs.Stat(fsys, root)
	if err != nil {
		err = fn(root, nil, err)
	} else {
		d := fs.FileInfoToDirEntry(info)
		err = fn(root, d, nil)
		// Walk root if it is a directory and err is nil
		if err == nil && d.IsDir() {
			entry := namedEntry{root, d}
			err = walkDir(fsys, []namedEntry{entry}, fn)
		}
	}
	if err == fs.SkipDir || err == fs.SkipAll {
		return nil
	}
	return err
}

type namedEntry struct {
	name string
	d    fs.DirEntry
}

// walkDir recursively descends path breadth first, calling walkDirFn.
func walkDir(fsys fs.FS, queue []namedEntry, walkDirFn fs.WalkDirFunc) error {
	if len(queue) == 0 {
		return nil
	}
	name, d := queue[0].name, queue[0].d
	queue = queue[1:] // remove first entry

	dirs, err := fs.ReadDir(fsys, name)
	if err != nil {
		// Second call, to report ReadDir error.
		err = walkDirFn(name, d, err)
		if err != nil {
			if err == fs.SkipDir && d.IsDir() {
				err = nil
			}
			return err
		}
	}

	var subqueue []namedEntry
	for _, d1 := range dirs {
		name1 := path.Join(name, d1.Name())
		err := walkDirFn(name1, d1, nil)
		if err != nil {
			if err == fs.SkipAll {
				return err
			}
			if err == fs.SkipDir {
				if d1.IsDir() {
					continue // skip current directory
				} else {
					subqueue = nil
					break // skip parent directory
				}
			}
		}
		if d1.IsDir() {
			subqueue = append(subqueue, namedEntry{name1, d1})
		}
	}
	queue = append(queue, subqueue...)

	return walkDir(fsys, queue, walkDirFn)
}
