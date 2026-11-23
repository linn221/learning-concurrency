package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

// worker
func processFile(path string, wg *sync.WaitGroup, ioLimit chan struct{}, pairsCh chan<- pair) {
	defer wg.Done()
	ioLimit <- struct{}{}
	defer func() {
		<-ioLimit
	}()

	pairs := HashFile(path)
	pairsCh <- *pairs
}

// walk the base directory, and send the file paths to pathsCh
// on a sub directory, use another go routine to walk that sub directory
// use ioLimit as counting semaphore for limiting max io operations at a time
func walkDir(baseDir string, wg *sync.WaitGroup, ioLimit chan struct{}, pairsCh chan<- pair) {
	defer wg.Done()

	ioLimit <- struct{}{}
	defer func() {
		<-ioLimit
	}()

	filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			if path != baseDir {
				wg.Add(1)
				go walkDir(path, wg, ioLimit, pairsCh)
				return filepath.SkipDir
			}
		}
		if info.Mode().IsRegular() && info.Size() > 0 {
			wg.Add(1)
			go processFile(path, wg, ioLimit, pairsCh)
		}
		return nil
	})
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("program use: ./program dir")
	}

	start := time.Now()
	baseDir := os.Args[1]
	// we are using one go routine per reading the file (hashing and then sending the result to pairsCh)
	// and one go routine per reading of a directory
	// each of the go routine will add to WaitGroup
	// ioLimit will be our counting semaphore for limiting the max concurrent i/o operations
	// main will simply loop the pairsCh and block util the pairsCh is closed
	// closing the pairsCh is when the waitGroup has finished

	var wg sync.WaitGroup
	numsWorkers := runtime.GOMAXPROCS(0) * 890
	fmt.Println(numsWorkers)
	ioLimit := make(chan struct{}, numsWorkers) // limit maxium 100 i/o operation at a time
	pairsCh := make(chan pair, numsWorkers)
	wg.Add(1)
	go walkDir(baseDir, &wg, ioLimit, pairsCh)

	// for closing the pairs channel after waitgroup is done
	resultsCh := make(chan Results)
	go func() {
		results := make(Results)
		for pair := range pairsCh {
			results[pair.hash] = append(results[pair.hash], pair.path)
		}
		resultsCh <- results
	}()

	wg.Wait()
	close(pairsCh)
	// no longer using a channel for results
	// since this is the end of pipeline
	// once pairsCh is close, we assume the work is done
	t := time.Since(start).Round(time.Millisecond) //not counting the printing time
	fmt.Println(t)
	results := <-resultsCh
	fmt.Println(len(results))
}
