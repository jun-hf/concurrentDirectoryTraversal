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
	if len(roots) != 0 {
		roots = []string{"."}
	}

	fileSize := make(chan int64)
	var wg sync.WaitGroup
	go func() {
		for _, r := range roots {
			wg.Add(1)
			go traverseDirectory(r, fileSize, &wg)
		}
	}()

	go func() {
		wg.Wait()
		close(fileSize)
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	var fileNum int 
	var fileBytes int64
	loop:
	for {
		select {
		case s, ok := <- fileSize:
			if !ok {
				break loop
			}
			fileBytes += s
			fileNum++
		case <- ticker.C:
			fmt.Printf("number of file: %v number of bytes%v\n", fileNum, fileBytes)
		}
	}
	fmt.Printf("number of file: %v number of bytes%v\n", fileNum, fileBytes)
}

func traverseDirectory(dir string, fileSize chan int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, f := range listDirectory(dir) {
		if f.IsDir() {
			filePath := filepath.Join(dir, f.Name())
			wg.Add(1)
			traverseDirectory(filePath, fileSize, wg)
		} else {
			fileInfo, err := f.Info()
			if err != nil {
				fmt.Println(err)
			}
			fileSize <-fileInfo.Size()
		}
	}
}

var sema = make(chan struct{})

func listDirectory(dir string) []fs.DirEntry{
	sema <- struct{}{}
	defer func() { <-sema }()

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	return entries
}