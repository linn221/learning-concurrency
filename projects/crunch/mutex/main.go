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
this employ multiple go routines
and use mutex to lock writing to file
faster than using channels but still not as fast as sequential
sequential takes 14 seconds, while this takes 2 minute on commented input
*/

func generator(min int, max int, chars string, outputfile string) (func(current []byte, pos int), func()) {
	charset := []byte(chars)
	var f func([]byte, int)
	file, err := os.Create(outputfile)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)
	var mu sync.Mutex
	f = func(buffer []byte, pos int) {
		if pos >= min {
			mu.Lock()
			w.Write(buffer[:pos])
			w.WriteString("\n")
			mu.Unlock()
			if pos == max {
				return
			}
		}

		for _, b := range charset {
			buffer[pos] = b
			f(buffer, pos+1)
		}
	}
	return f, func() {
		w.Flush()
		file.Close()
	}
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
	generate, close := generator(min, max, chars, outputfile)
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

	wg.Wait()
	close()
	ttk := time.Since(start).Round(time.Millisecond)
	fmt.Println(ttk)
}
