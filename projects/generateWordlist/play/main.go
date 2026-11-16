package main

import (
	"fmt"
	"sync"
)

var charSet []byte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890~!@#$%^&*()_+=-./<>?,;:")

func main() {
	var wg sync.WaitGroup
	wg.Add(20)
	wg.Done()
	wg.Add(-19)
	wg.Wait()
	fmt.Println("done")
}
