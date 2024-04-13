package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type rootSize struct {
	root string
	byteSize int64
}

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	rootCh := make(chan rootSize)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range roots {
			wg.Add(1)
			go TraverseDirectory(r, &wg, rootCh, r)
		}
	}()

	go func() {
		wg.Wait()
		close(rootCh)
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	var rootContent = []rootSize{}
	var totalFiles int
loop:
	for {
		select {
		case rootResult, ok := <-rootCh:
			if !ok {
				break loop
			}
			rootContent = append(rootContent, rootResult)
			totalFiles++
		case <-ticker.C:
			printRootContent(rootContent)
		}
	}

	printRootContent(rootContent)
}

func printRootContent(rootList []rootSize) {
	rootMap := make(map[string]int64)
	for _, r := range rootList {
		rootMap[r.root] += r.byteSize
	}
	for key, value:= range rootMap {
		fmt.Printf("root: %q size: %v bytes\n", key, value)
	}
}

func TraverseDirectory(dir string, wg *sync.WaitGroup, rootResponse chan<- rootSize, currentRoot string) {
	defer wg.Done()
	for _, f := range ListDirectory(dir) {
		if f.IsDir() {
			dirPath := filepath.Join(dir, f.Name())
			wg.Add(1)
			TraverseDirectory(dirPath, wg, rootResponse, currentRoot)
		} else {
			fileInfo, err := f.Info()
			if err != nil {
				fmt.Println(err)
			}
			rootResponse <- rootSize{currentRoot, fileInfo.Size()}
		}
	}
}

func ListDirectory(dir string) []fs.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	return entries
}
