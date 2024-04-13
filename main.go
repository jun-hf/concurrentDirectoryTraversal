package main

// import (
// 	"flag"
// 	"fmt"
// 	"io/fs"
// 	"os"
// 	"path/filepath"
// 	"sync"
// 	"time"
// )

// func main() {
// 	flag.Parse()
// 	roots := flag.Args()
// 	if len(roots) == 0 {
// 		roots = []string{"."}
// 	}

// 	fileSize := make(chan int64)
// 	var wg sync.WaitGroup
// 	go func() {
// 		for _, r := range roots {
// 			wg.Add(1)
// 			go traverseDirectory(r, fileSize, &wg)
// 		}
// 	}()

// 	go func() {
// 		wg.Wait()
// 		fmt.Println("Closing fileSize channel")
// 		close(fileSize)
// 	}()

// 	var bytesNum int64
// 	var fileNum int
// 	ticker := time.NewTicker(500 *time.Millisecond)
// 	loop:
// 	for {
// 		select{
// 		case size, ok := <- fileSize:
// 			if !ok {
// 				break loop
// 			}
// 			bytesNum += size
// 			fileNum++
// 		case <-ticker.C:
// 			printFileStorage(fileNum, bytesNum)
// 		}
// 	}
// 	printFileStorage(fileNum, bytesNum)
// }

// func listDirContent(dir string) []fs.DirEntry {
// 	dirEntries, err := os.ReadDir(dir)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return dirEntries
// }

// func traverseDirectory(dir string, fileSize chan int64, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for _, f := range listDirContent(dir) {
// 		if f.IsDir() {
// 			dirPath := filepath.Join(dir, f.Name())
// 			wg.Add(1)
// 			traverseDirectory(dirPath, fileSize, wg)
// 		} else {
// 			fInfo, err := f.Info()
// 			if err != nil {
// 				fmt.Printf("for %v unable to get info: %v",f.Name(), err)
// 			}
// 			fileSize <- fInfo.Size()	
// 		}
// 	}
// }

// func printFileStorage(fileNum int, bytesNum int64) {
// 	fmt.Printf("number of files: %v, number of bytes: %v\n", fileNum, bytesNum)
// }