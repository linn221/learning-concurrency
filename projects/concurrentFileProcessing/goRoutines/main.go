package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type pair struct {
	hash string
	path string
}
type filePath []string
type Results map[string]filePath

func HashFile(path string) *pair {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatal(err)
	}

	checksum := hash.Sum(nil)
	return &pair{
		hash: fmt.Sprintf("%x", checksum),
		path: path,
	}
}

func DoWork(path string, results *MyResult) {
	p := HashFile(path)
	results.mu.Lock()
	defer results.mu.Unlock()
	results.Results[p.hash] = append(results.Results[p.hash], p.path)
}

type MyResult struct {
	Results
	mu sync.Mutex
}

func SearchDir(dir string) Results {
	var wg sync.WaitGroup
	var myResult MyResult
	results := make(Results)
	myResult.Results = results
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsRegular() && info.Size() > 0 {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				DoWork(path, &myResult)
			}(path)
			// results[pair.hash] = append(results[pair.hash], pair.path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	return myResult.Results
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("program use: ./program dir")
	}

	start := time.Now()
	_ = SearchDir(os.Args[1])
	t := time.Since(start).Round(time.Millisecond) //not counting the printing time
	// for hash, fileList := range results {
	// if len(fileList) > 1 {
	// 	fmt.Printf("%s\n", hash)
	// 	for _, path := range fileList {
	// 		fmt.Printf("\t%s\n", path)
	// 	}
	// 	fmt.Println()
	// }
	// }
	fmt.Println(t)
}
