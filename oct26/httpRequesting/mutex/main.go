package main

import (
	"fmt"
	"sync"

	"github.com/linn221/learning-concurrency/oct26/httpRequesting/lib"
)

func main() {
	var mu sync.Mutex
	var results []lib.Result
	var wg sync.WaitGroup

	for _, url := range lib.Urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := lib.FetchUrl(url)
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
			r.Print()
		}()
	}

	wg.Wait()
	fmt.Println(len(results))
}
