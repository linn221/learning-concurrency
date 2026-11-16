package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func work(i int, ch <-chan int) {
	v := <-ch
	fmt.Printf("%d got the value %d\n", i, v)
}

func main() {
	// Creates a context that cancels when SIGINT (Ctrl+C) or SIGTERM is received
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ch := make(chan int)
	for i := range 20 {
		go work(i, ch)
	}
	time.Sleep(time.Second * 1)
	for i := range 20 {
		ch <- i
		time.Sleep(200)
	}

	// Blocks main until an interrupt signal is received
	<-ctx.Done()
	fmt.Println("\nShutting down gracefully...")
}
