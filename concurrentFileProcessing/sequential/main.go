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

func SearchDir(dir string) Results {
	results := make(Results)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsRegular() && info.Size() > 0 {
			pair := HashFile(path)
			results[pair.hash] = append(results[pair.hash], pair.path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return results
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
