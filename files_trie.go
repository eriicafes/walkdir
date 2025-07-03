package main

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type TrieNode struct {
	Children     map[string]*TrieNode
	ContentFiles []string
	LayoutPath   string
}

// NewTrieNode creates a new TrieNode.
func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children: make(map[string]*TrieNode),
	}
}

func WalkFilesWithLayoutTrie(fsys fs.FS, ext string, layoutFilename string, dirFilter string) map[string][]string {
	root := NewTrieNode()
	groups := make(map[string][]string)

	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		pathWithoutExt := strings.TrimSuffix(path, "."+ext)
		if ext != "" && path == pathWithoutExt {
			return nil
		}

		pathNormalized := filepath.ToSlash(pathWithoutExt)
		dir, filename := filepath.Split(pathNormalized)
		dir = strings.TrimSuffix(dir, "/")

		parts := []string{}
		if dir != "" {
			parts = strings.Split(dir, "/")
		}

		curr := root
		for _, part := range parts {
			if _, ok := curr.Children[part]; !ok {
				curr.Children[part] = NewTrieNode()
			}
			curr = curr.Children[part]
		}

		if filename == layoutFilename {
			curr.LayoutPath = pathNormalized
		} else {
			curr.ContentFiles = append(curr.ContentFiles, pathNormalized)
		}
		return nil
	})

	// 2. Traverse the Trie to collect files and their associated layouts.
	var traverse func(*TrieNode, []string)
	traverse = func(node *TrieNode, layoutStack []string) {
		currentLayouts := layoutStack
		if node.LayoutPath != "" {
			currentLayouts = append(currentLayouts, node.LayoutPath)
		}

		for _, filePath := range node.ContentFiles {
			if dirFilter == "." || strings.HasPrefix(filePath, dirFilter) {
				resultFiles := make([]string, 0, len(currentLayouts)+1)
				resultFiles = append(resultFiles, currentLayouts...)
				resultFiles = append(resultFiles, filePath)
				groups[filePath] = resultFiles
			}
		}

		// Recurse to children directories.
		for _, child := range node.Children {
			traverse(child, currentLayouts)
		}
	}

	traverse(root, []string{})

	return groups
}
