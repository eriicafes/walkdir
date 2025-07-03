package main

import (
	"maps"
	"slices"
	"testing"

	"io/fs"
	"testing/fstest"
)

// func BenchmarkWalk(b *testing.B) {
// 	smFsys := generateFS("data", 3, 2)
// 	lgFsys := generateFS("data", 100, 5)

// 	cases := []struct {
// 		name     string
// 		fsys     fs.FS
// 		walkFunc func(fsys fs.FS, root string, fn fs.WalkDirFunc) error
// 	}{
// 		{"FS Small", smFsys, fs.WalkDir},
// 		{"BreadthFirst Small", smFsys, walkDirBreadthFirst},
// 		{"FS Large", lgFsys, fs.WalkDir},
// 		{"BreadthFirst Large", lgFsys, walkDirBreadthFirst},
// 	}

// 	for _, c := range cases {
// 		b.Run(c.name, func(b *testing.B) {
// 			b.ReportAllocs()
// 			for b.Loop() {
// 				visited := 0
// 				c.walkFunc(c.fsys, "data", func(path string, d fs.DirEntry, err error) error {
// 					if err != nil {
// 						return err
// 					}
// 					visited++
// 					return nil
// 				})
// 			}
// 		})
// 	}
// }

func BenchmarkWalkLayout(b *testing.B) {
	smFsys := generateFS("data", 3, 2)
	mdFsys := generateFS("data", 30, 4)
	lgFsys := generateFS("data", 100, 5)

	cases := []struct {
		name     string
		fsys     fs.FS
		walkFunc func(fsys fs.FS, ext string, layoutFilename string, dir string) map[string][]string
	}{
		{"Initial Small", smFsys, WalkFilesWithLayout},
		{"BreadthFirst Small", smFsys, WalkFilesWithLayoutBreadthFirst},
		{"Trie Small", smFsys, WalkFilesWithLayoutTrie},
		{"Initial Medium", mdFsys, WalkFilesWithLayout},
		{"BreadthFirst Medium", mdFsys, WalkFilesWithLayoutBreadthFirst},
		{"Trie Medium", mdFsys, WalkFilesWithLayoutTrie},
		{"Initial Large", lgFsys, WalkFilesWithLayout},
		{"BreadthFirst Large", lgFsys, WalkFilesWithLayoutBreadthFirst},
		{"Trie Large", lgFsys, WalkFilesWithLayoutTrie},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				c.walkFunc(c.fsys, "html", "layout", "data")
			}
		})
	}
}

func TestWalkFilesWithLayout(t *testing.T) {
	fsys := fstest.MapFS{
		"index.html":              {},
		"index.tmpl":              {},
		"index":                   {},
		"layout.html":             {},
		"test/layout.html":        {},
		"test/index.tmpl":         {},
		"test/index":              {},
		"app/layout.html":         {},
		"app/index.html":          {},
		"app/dashboard.html":      {},
		"app/dashboard.tmpl":      {},
		"app/dashboard":           {},
		"app/account/layout.html": {},
		"app/account/index.html":  {},
		"auth/index.tmpl":         {},
		"auth/index":              {},
		"auth/login.html":         {},
		"auth/register.html":      {},
	}

	// root
	expectedRoot := map[string][]string{
		"index":             {"layout", "index"},
		"app/index":         {"layout", "app/layout", "app/index"},
		"app/dashboard":     {"layout", "app/layout", "app/dashboard"},
		"app/account/index": {"layout", "app/layout", "app/account/layout", "app/account/index"},
		"auth/login":        {"layout", "auth/login"},
		"auth/register":     {"layout", "auth/register"},
	}
	// app sub dir
	expectedSub := map[string][]string{
		"app/index":         {"layout", "app/layout", "app/index"},
		"app/dashboard":     {"layout", "app/layout", "app/dashboard"},
		"app/account/index": {"layout", "app/layout", "app/account/layout", "app/account/index"},
	}

	cases := []struct {
		name     string
		walkFunc func(fsys fs.FS, ext string, layoutFilename string, dir string) map[string][]string
	}{
		{"Initial", WalkFilesWithLayout},
		{"BreadthFirst", WalkFilesWithLayoutBreadthFirst},
		{"rie", WalkFilesWithLayoutTrie},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// walk root
			got := c.walkFunc(fsys, "html", "layout", ".")
			if !maps.EqualFunc(expectedRoot, got, slices.Equal) {
				t.Errorf("root expected: %v, got: %v", expectedRoot, got)
			}

			// walk sub dir
			got = c.walkFunc(fsys, "html", "layout", "app")
			if !maps.EqualFunc(expectedSub, got, slices.Equal) {
				t.Errorf("sub dir expected: %v, got: %v", expectedSub, got)
			}
		})
	}
}
