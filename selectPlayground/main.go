package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		for i := range 10 {
			log.Println("channel 1 value generating")
			ch1 <- i
		}
	}()
	go func() {
		for i := range 10 {
			log.Println("channel 2 value generating")
			ch2 <- i * 2
		}
	}()
	go func() {
		for i := range 10 {
			log.Println("channel 3 value generating")
			ch3 <- i * 3
		}
	}()

	for {
		select {
		case v := <-ch1:
			log.Println("channel 1 selected, value is " + fmt.Sprint(v) + " and sleeping for 2 second")
			time.Sleep(time.Second * 2)
		case v := <-ch2:
			log.Println("channel 2 selected, value is " + fmt.Sprint(v) + " and sleeping for 2 second")
			time.Sleep(time.Second * 2)
		case v := <-ch3:
			log.Println("channel 3 selected, value is " + fmt.Sprint(v) + " and sleeping for 2 second")
			time.Sleep(time.Second * 2)
		default:
			log.Printf("default case selected")
			log.Println("default case started working, sleeping for 3 seconds")
			time.Sleep(3 * time.Second)
			log.Println("default case working done")
		}
	}
}
