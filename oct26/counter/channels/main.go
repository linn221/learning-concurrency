package main

import (
	"fmt"
	"net/http"
)

func count(ch chan<- int) {
	for i := 1; ; i++ {
		ch <- i
	}
}

func main() {
	ch := make(chan int)

	go count(ch)
	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		counterValue := <-ch
		fmt.Fprintf(w, "%d", counterValue)
	})
	http.ListenAndServe(":8080", nil)
}
