package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func generator(min int, max int, chars string, outputfile string) (func(current []byte, pos int), func()) {
	charset := []byte(chars)
	var f func([]byte, int)
	file, err := os.Create(outputfile)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)
	close := func() {
		w.Flush()
		file.Close()
	}
	f = func(buffer []byte, pos int) {
		if pos >= min {
			// w.WriteString(string(buffer[:pos]))
			w.Write((buffer[:pos]))
			w.WriteString("\n")
			if pos == max {
				return
			}
		}

		for _, b := range charset {
			buffer[pos] = b
			f(buffer, pos+1)
		}
	}
	return f, close
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
	buffer := make([]byte, max)
	generate(buffer, 0)
	close()
	ttk := time.Since(start).Round(time.Millisecond)
	fmt.Println(ttk)
}
