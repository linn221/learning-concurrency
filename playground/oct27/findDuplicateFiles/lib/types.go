package lib

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
