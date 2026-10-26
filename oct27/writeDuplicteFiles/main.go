package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

const (
	totalFiles          = 500_000
	duplicatePercentage = 70
	maxDuplicateFiles   = 15
	maxContentSize      = 5000
)

// run => go run . ../dumps

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 \n\t.,!?-_:;"

func randomInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

func randomName() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b) + ".txt"
}

func randomContent() string {
	length := randomInt(maxContentSize + 1)

	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		sb.WriteByte(charset[randomInt(len(charset))])
	}

	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}

	dir := os.Args[1]

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Create directory:", err)
		os.Exit(1)
	}

	duplicateCount := totalFiles * duplicatePercentage / 100
	files := make([]string, 0, totalFiles)
	seen := make(map[string]bool)

	// Generate duplicate-content groups.
	remaining := duplicateCount

	for remaining > 0 {
		groupSize := randomInt(maxDuplicateFiles-1) + 2

		if groupSize > remaining {
			groupSize = remaining
		}

		content := randomContent()
		for seen[content] {
			content = randomContent()
		}
		seen[content] = true

		for i := 0; i < groupSize; i++ {
			files = append(files, content)
		}

		remaining -= groupSize
	}

	// Generate the remaining files with unique content.
	for len(files) < totalFiles {
		content := randomContent()

		if seen[content] {
			continue
		}

		seen[content] = true
		files = append(files, content)
	}

	// Shuffle the contents so duplicate files aren't grouped together.
	for i := len(files) - 1; i > 0; i-- {
		j := randomInt(i + 1)
		files[i], files[j] = files[j], files[i]
	}

	// Write the files using random filenames.
	for i, content := range files {
		var path string

		for {
			path = filepath.Join(dir, randomName())

			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				break
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "Check filename:", err)
				os.Exit(1)
			}
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Write file %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Created %d files in %s\n", totalFiles, dir)
	fmt.Printf("Files in duplicate groups: %d (%d%%)\n",
		duplicateCount, duplicatePercentage)
	fmt.Printf("Files with unique content: %d\n",
		totalFiles-duplicateCount)
}
