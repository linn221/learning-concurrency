package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	counter := 1
	var mu sync.Mutex

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", counter)
		mu.Lock()
		counter++
		mu.Unlock()
	})
	http.ListenAndServe(":8080", nil)
}
