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

var done = make(chan bool)
func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

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
		case <-done:
			// drain the chanel, before exiting the loop
			for range rootCh {
			}
			return
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
	if cancelled() {
		return
	}
	for _, f := range ListDirectory(dir) {
		if f.IsDir() {
			dirPath := filepath.Join(dir, f.Name())
			wg.Add(1)
			go TraverseDirectory(dirPath, wg, rootResponse, currentRoot)
		} else {
			fileInfo, err := f.Info()
			if err != nil {
				fmt.Println(err)
			}
			rootResponse <- rootSize{currentRoot, fileInfo.Size()}
		}
	}
}

// sema is a semaphore that limits ListDirectory to 20 
var sema = make(chan struct{}, 20)

func ListDirectory(dir string) []fs.DirEntry {
	select {
	case sema <- struct{}{}:
	case <-done:
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	return entries
}