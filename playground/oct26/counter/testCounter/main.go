package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	url := "http://localhost:8080/count"
	for range 2000 {
		wg.Go(func() {
			http.Get(url)
		})
	}

	wg.Wait()
	resp, _ := http.Get(url)
	text, _ := io.ReadAll(resp.Body)
	// the counter must be 2001
	fmt.Println(string(text))
}
