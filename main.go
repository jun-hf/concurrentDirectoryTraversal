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

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	fileSize := make(chan int64)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range roots {
			wg.Add(1)
			go TraverseDirectory(r, &wg, fileSize)
		}
	}()

	go func() {
		wg.Wait()
		close(fileSize)
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	var totalBytes int64
	var totalFiles int
loop:
	for {
		select {
		case size, ok := <-fileSize:
			if !ok {
				break loop
			}
			totalBytes += size
			totalFiles++
		case <-ticker.C:
			fmt.Printf("number of files: %v, number of bytes: %v\n", totalFiles, totalBytes)
		}
	}

	fmt.Printf("number of files: %v, number of bytes: %v\n", totalFiles, totalBytes)
}

func TraverseDirectory(dir string, wg *sync.WaitGroup, fileSize chan<- int64) {
	defer wg.Done()
	for _, f := range ListDirectory(dir) {
		if f.IsDir() {
			dirPath := filepath.Join(dir, f.Name())
			wg.Add(1)
			TraverseDirectory(dirPath, wg, fileSize)
		} else {
			fileInfo, err := f.Info()
			if err != nil {
				fmt.Println(err)
			}
			fileSize <- fileInfo.Size()
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
