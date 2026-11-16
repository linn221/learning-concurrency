package main

import (
	"fmt"
	"log"
)

func write(v int, ch chan<- int) {

	log.Println("starting to write")
	ch <- v
	log.Println("done writing")
}

func read(ch <-chan int) int {

	// time.Sleep(time.Second * 1)
	log.Println("starting to read")
	v := <-ch
	log.Println("done reading")
	return v
}

func main() {

	// ch := make(chan int)
	ch := make(chan int, 3)

	for i := range 10 {
		go write(i, ch)
		go read(ch)
	}
	fmt.Scanln()
}
