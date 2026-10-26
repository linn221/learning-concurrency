A `select` blocks until one of its cases can run, then it executes that case. It chooses one at random if multiple are ready.

The `default` case in a `select` is run if no other case is ready.

- listens on more than one channel at a time.
- whatever channel return first on the case statement will be read.
### Example
```go
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:  // write to the channel c, not listening to it
		// writing to the channel is inside the case, therefore, will keep writing to it until quit channel is read to be read, then return the function
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)
	// reads the channel for 10 times, then write to quit channel to signal cancellation
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	
	// write to the channel, listens to the cancellation channel using select
	fibonacci(c, quit)
}

```

functions that return channels
- time.After(time.Second * 5)
- time.NewTicker(time.Second * 5).C
listening to those channels is a way to control the execution.

[context](context.md) has a Done that returns a signaling channel.
