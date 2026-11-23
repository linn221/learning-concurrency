package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

/*
this employ multiple go routines
but each go routine write to its own file, thereforce removing the need for synchronization entirely
and later merging to one output file

sequential takes 14 seconds, while this takes 8.5 with commented input
this runs slower than sequential when the charset is not large enough, (0-9) due to the cost of creating separate files and merging is larger than generating the words
*/

func merge(output string, files []string) error {
	dst, err := os.Create(output)
	if err != nil {
		return err
	}
	defer dst.Close()

	for _, filename := range files {
		src, err := os.Open(filename)
		if err != nil {
			return err
		}

		_, err = io.Copy(dst, src)
		src.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func generator(min int, max int, chars string) func(current []byte, pos int, fileWrite func([]byte)) {
	charset := []byte(chars)
	var f func([]byte, int, func([]byte))
	f = func(buffer []byte, pos int, fileWrite func([]byte)) {
		if pos >= min {
			fileWrite(buffer[:pos])
			if pos == max {
				return
			}
		}

		for _, b := range charset {
			buffer[pos] = b
			f(buffer, pos+1, fileWrite)
		}
	}
	return f
}

func main() {
	min := 8
	max := 8
	chars := "0123456789"
	start := time.Now()
	// min := 5
	// max := 5
	// chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	outputfile := "words.txt"
	generate := generator(min, max, chars)
	var wg sync.WaitGroup
	charset := []byte(chars)
	wg.Add(len(charset))
	go func() {
		for i, char := range charset {
			buffer := make([]byte, max)
			buffer[0] = char
			go func() {
				defer wg.Done()
				file, err := os.Create(fmt.Sprintf("worker_%d.txt", i))
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				w := bufio.NewWriter(file)
				defer w.Flush()
				generate(buffer, 1, func(b []byte) {
					w.Write(b)
					w.WriteString("\n")
				})
			}()
		}
	}()

	wg.Wait()
	files := make([]string, len(charset))
	for i := range charset {
		files[i] = fmt.Sprintf("worker_%d.txt", i)
	}
	merge(outputfile, files)
	ttk := time.Since(start).Round(time.Millisecond)
	fmt.Println(ttk)
}
