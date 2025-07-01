package main

import (
	"fmt"
	"os"
)

func main() {
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
