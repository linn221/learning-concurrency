package main

import "github.com/linn221/learning-concurrency/oct26/httpRequesting/lib"

func makeRequest(url string, ch chan<- lib.Result) {
	result := lib.FetchUrl(url)
	ch <- result
}

func main() {
	ch := make(chan lib.Result)

	for _, url := range lib.Urls {
		go makeRequest(url, ch)
	}

	for range lib.Urls {
		result := <-ch
		result.Print()
	}

	// for r := range ch {
	// 	r.Print()
	// }
}
