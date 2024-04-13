package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	
	var wg sync.WaitGroup
	fileSize := make(chan int64)

	for _, r := range roots{
		wg.Add(1)
		go TraverseDirectory(r, &wg, fileSize)
	}

	go func() {
		wg.Wait()
		close(fileSize)
	}()
	

	for a := range fileSize {
		fmt.Println(a)
	}
}

func TraverseDirectory(dir string, wg *sync.WaitGroup, fileSize chan <-int64) {
	defer wg.Done()
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range entries {
		if !f.IsDir() {
			fileInfo, err := f.Info()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("adding file: %v\n", f.Name())
			fileSize <- fileInfo.Size()
		}
	}
}