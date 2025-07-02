package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type dataset struct {
	Name         string
	Dir          string
	NumDirs      int
	NestingDepth int
}

var fileTypes = []string{"html", "go", "ts", "css"}

const filesPerDir = 10

func main() {
	// Resolve script directory (not CWD)
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to get current file path")
	}
	baseDir := filepath.Dir(currentFile)

	// Define datasets
	datasets := []dataset{
		{Name: "data_lg", Dir: filepath.Join(baseDir, "..", "example", "data_lg"), NumDirs: 100, NestingDepth: 5},
		{Name: "data_sm", Dir: filepath.Join(baseDir, "..", "example", "data_sm"), NumDirs: 3, NestingDepth: 2},
	}

	for _, ds := range datasets {
		generateDataset(ds)
	}
}

func generateDataset(ds dataset) {
	log.Printf("Generating dataset %s at %s...\n", ds.Name, ds.Dir)

	// Cleanup
	if err := os.RemoveAll(ds.Dir); err != nil {
		log.Fatalf("failed to clean directory %s: %v", ds.Dir, err)
	}

	htmlCount := 0

	for i := range make([]struct{}, ds.NumDirs) {
		currentPath := ds.Dir

		for d := range make([]struct{}, ds.NestingDepth) {
			currentPath = filepath.Join(currentPath, fmt.Sprintf("dir%d_%d", i, d))

			if err := os.MkdirAll(currentPath, 0755); err != nil {
				log.Fatalf("failed to create dir: %v", err)
			}

			for f := range make([]struct{}, filesPerDir) {
				ext := fileTypes[f%len(fileTypes)]
				isLayout := f%5 == 0

				filename := fmt.Sprintf("file%d.%s", f, ext)
				if isLayout {
					filename = fmt.Sprintf("layout.%s", ext)
				}

				if ext == "html" && !isLayout {
					htmlCount++
				}

				filePath := filepath.Join(currentPath, filename)
				content := fmt.Sprintf("// dummy content for %s\n", filePath)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					log.Fatalf("write failed: %v", err)
				}
			}
		}
	}

	// Write count to example/data_<name>.txt
	countPath := filepath.Join(filepath.Dir(ds.Dir), fmt.Sprintf("%s.txt", ds.Name))
	if err := os.WriteFile(countPath, fmt.Appendf(nil, "%d", htmlCount), 0644); err != nil {
		log.Fatalf("failed to write count file: %v", err)
	}

	log.Printf("✅ %s: %d .html files (excluding layout.html)\n📄 Count written to: %s\n", ds.Name, htmlCount, countPath)
}
