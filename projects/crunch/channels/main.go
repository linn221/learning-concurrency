package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

/*
this actually run slower than sequential
due to channels overhead
*/

func generator(min int, max int, chars string, wordsCh chan<- []byte) func(current []byte, pos int) {
	charset := []byte(chars)
	var f func([]byte, int)
	f = func(buffer []byte, pos int) {
		if pos >= min {
			wordsCh <- buffer[:pos]
			if pos == max {
				return
			}
		}

		for _, b := range charset {
			buffer[pos] = b
			f(buffer, pos+1)
		}
	}
	return f
}
func main() {
	min := 8
	max := 8
	start := time.Now()
	chars := "0123456789"
	outputfile := "words.txt"
	wordsCh := make(chan []byte, 4096)
	generate := generator(min, max, chars, wordsCh)
	var wg sync.WaitGroup

	charset := []byte(chars)
	wg.Add(len(charset))
	go func() {
		for _, char := range charset {
			buffer := make([]byte, max)
			buffer[0] = char
			go func() {
				defer wg.Done()
				generate(buffer, 1)
			}()
		}
	}()

	go func() {
		wg.Wait()
		close(wordsCh)
	}()
	file, err := os.Create(outputfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	defer w.Flush()
	for word := range wordsCh {
		w.Write(word)
		w.WriteString("\n")
	}

	ttk := time.Since(start).Round(time.Millisecond)
	fmt.Println(ttk)
}
