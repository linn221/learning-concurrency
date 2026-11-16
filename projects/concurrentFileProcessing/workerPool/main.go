package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

// worker
func processPaths(pathsCh <-chan string, pairsCh chan<- pair, doneCh chan<- bool) {
	for path := range pathsCh {
		pairs := HashFile(path)
		pairsCh <- *pairs
	}
	doneCh <- true
}

func collectHashes(pairsCh <-chan pair, resultsCh chan<- Results) {
	results := make(Results)
	for pair := range pairsCh {
		results[pair.hash] = append(results[pair.hash], pair.path)
	}
	resultsCh <- results
}

// the main problem was who close the Pairs ch
// we solved the problem by passing doneCh to worker pool, and then,
// closing the Pairs ch only after doneCh has been sent
// and reading from the doneCh the number of workers time
// we could also use a waitGroup instead of a doneCh (adding one on creating a worker, done when processing finish, and wait on main function)
func SearchDir(dir string) Results {
	pathsCh := make(chan string)
	pairsCh := make(chan pair)
	doneCh := make(chan bool)
	resultsCh := make(chan Results)

	numWorkers := 120
	for i := 0; i < numWorkers; i++ {
		go processPaths(pathsCh, pairsCh, doneCh)
	}

	go collectHashes(pairsCh, resultsCh)

	results := make(Results)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsRegular() && info.Size() > 0 {
			pathsCh <- path
			// pair := HashFile(path)
			// results[pair.hash] = append(results[pair.hash], pair.path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	close(pathsCh)
	// after pathsCh is closed, we will wait for closeCh of each worker
	for range numWorkers {
		<-doneCh
	}
	// then close the pairs ch, indicating the processing has been completed
	close(pairsCh)
	// now we will read from the resultsCh

	results = <-resultsCh
	return results
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("program use: ./program dir")
	}

	start := time.Now()
	results := SearchDir(os.Args[1])
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
	fmt.Println(len(results))
}
