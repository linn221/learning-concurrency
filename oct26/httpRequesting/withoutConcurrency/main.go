package main

import "github.com/linn221/learning-concurrency/oct26/httpRequesting/lib"

func main() {
	for _, url := range lib.Urls {
		result := lib.FetchUrl(url)
		result.Print()
	}
}
