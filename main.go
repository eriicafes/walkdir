package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

func main() {
	err := walkDirBreadthFirst(os.DirFS("example"), ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println("path:", path)
		defer fmt.Println("")

		if d.Name() == "test.txt" {
			fmt.Println("skip dir")
			return fs.SkipDir
		}

		if d.Name() == "node_modules" {
			fmt.Println("skip dir")
			return fs.SkipDir
		}

		if err != nil {
			fmt.Println("error:", err)
			return err
		}
		if d.IsDir() {
			fmt.Println("is dir:", d.IsDir())
		}
		return nil
	})
	log.Fatal("done:", err)

	fmt.Println("start go")
	fsys := os.DirFS("example")

	filesWithLayouts := WalkFilesWithLayout(
		fsys,
		"html",
		"layout",
		"app",
	)

	for file, paths := range filesWithLayouts {
		fmt.Println(file, "->", paths)
	}
	fmt.Println("end go")
}
