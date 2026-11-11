package main

import (
	"fmt"
	"time"
)

var count int

func work() {
	count++
}

func main() {

	for range 2000 {
		go work()
	}
	time.Sleep(time.Second * 2)
	fmt.Println(count)
}
